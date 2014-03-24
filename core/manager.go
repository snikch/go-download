// Provides a state management interface for all current
// downloads and the settings object
package core

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/nsf/termbox-go"

	"code.google.com/p/go.crypto/ssh/terminal"
)

type Manager struct {
	settings     *Settings
	Downloads    map[url.URL]*Downloader
	Publisher    *Publisher
	progressChan chan *Downloader
}

func NewManager(s *Settings) *Manager {
	m := &Manager{
		Downloads:    make(map[url.URL]*Downloader),
		progressChan: make(chan *Downloader),
		settings:     s,
	}
	return m
}

func (m *Manager) RestoreState(s *Settings) error {
	return nil
}

func (m *Manager) AddUrl(url string) (err error) {
	downloader, err := NewDownloader(url, m.progressChan)
	if err != nil {
		return
	}

	if m.Downloads[*downloader.Resource.Url] != nil {
		return fmt.Errorf("Already downloading %s", url)
	}
	m.Downloads[*downloader.Resource.Url] = &downloader
	downloader.Start(true)
	return
}

func (m *Manager) Start() {
	disp := Display{manager: m}
	terminalW, _, _ := terminal.GetSize(0)
	disp.width = terminalW
	disp.drawLine("Starting", 0)
	for {
		select {
		case _, _ = <-m.progressChan:
			disp.downloading = true
			disp.draw()
			time.Sleep(100 * time.Millisecond)

		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

type Display struct {
	percent     int
	downloading bool
	width       int
	manager     *Manager
}

func (d *Display) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	y := 0
	d.drawLine(fmt.Sprintf("%d downloads", len(d.manager.Downloads)), y)
	for _, dwn := range d.manager.Downloads {
		y += 1
		d.drawDownload(dwn, y)
	}
	termbox.Flush()
}
func (d *Display) drawLine(s string, y int) {
	c := termbox.ColorDefault
	termbox.SetCursor(0, 0)
	x := 0
	for _, r := range s {
		termbox.SetCell(x, y, r, c, c)
		x += 1
	}
}
func (d *Display) drawDownload(dwn *Downloader, y int) {

	if dwn.Downloaded == ByteSize(0) {
		d.drawLine(
			fmt.Sprintf("%s download starting", dwn.Resource.Hoster.Name()),
			y,
		)
		return
	}
	percent := dwn.Progress()
	lineS := fmt.Sprintf(
		"%s %d%% (%s/%s) %s/s [",
		dwn.Resource.Name,
		percent,
		dwn.Downloaded,
		dwn.Size,
		dwn.SpeedMonitor.Speed,
	)
	lineE := fmt.Sprintf(
		"] %d/%d Chunks",
		dwn.CompleteChunks,
		dwn.TotalChunks,
	)

	barW := d.width - len(lineS) - len(lineE)

	chunkW := int(math.Floor(float64(barW) / float64(dwn.TotalChunks)))

	var bar string
	for i := 0; i < dwn.TotalChunks; i++ {
		c := dwn.Chunks[i]
		if c.Percent == 0 {
			bar += strings.Repeat(".", chunkW)
			continue
		}
		if c.Percent == 100 {
			bar += strings.Repeat("+", chunkW)
			continue
		}
		completeW := int(math.Floor(float64(c.Percent) / float64(100) * float64(chunkW)))
		if completeW > 0 {
			completeW -= 1
		}
		incompleteW := chunkW - completeW - 1
		var incompleteS string
		if incompleteW > 0 {
			incompleteS = strings.Repeat(".", incompleteW)
		}
		bar += strings.Repeat("+", completeW) + ">" + incompleteS
	}

	line := lineS + bar + lineE
	spaceW := d.width - len(line)
	if spaceW > 0 {
		line += strings.Repeat(" ", spaceW)
	}
	d.drawLine(line, y)
}
