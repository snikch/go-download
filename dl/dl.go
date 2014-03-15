package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/snikch/go-download/core"

	"code.google.com/p/go.crypto/ssh/terminal"
)

var (
	resume = flag.Bool("resume", true, "Attempt to resume unfinished downloads")
	file   = flag.String("file", "", "URL of file to download")
	chunks = flag.Int("chunks", 6, "Number of chunks to download at once")
)

func main() {

	flag.Parse()
	if *file == "" {
		panic(fmt.Errorf("No file provided"))
	}
	downloader, err := core.NewDownloader(*file)
	if err != nil {
		panic(err)
	}

	progress, complete, err := downloader.Start(*resume)
	if err != nil {
		panic(err)
	}
	disp := Display{downloader: &downloader}
	terminalW, _, _ := terminal.GetSize(0)
	disp.width = terminalW
Loop:
	for {
		select {
		case percent := <-progress:
			disp.downloading = true
			disp.percent = percent
		case complete := <-complete:
			if complete {
				disp.draw()
				fmt.Println("\nFinished downloading")
				break Loop
			}

		default:
			disp.draw()
			time.Sleep(100 * time.Millisecond)
		}
	}
}

type Display struct {
	downloader  *core.Downloader
	percent     int
	downloading bool
	width       int
}

func (d *Display) draw() {
	if d.downloading == false {
		return
	}

	lineS := fmt.Sprintf(
		"%s %d%% (%s/%s) %s/s [",
		d.downloader.Resource.Name,
		d.percent,
		d.downloader.Downloaded,
		d.downloader.Size,
		d.downloader.SpeedMonitor.Speed,
	)
	lineE := fmt.Sprintf(
		"] %d/%d Chunks",
		d.downloader.CompleteChunks,
		d.downloader.TotalChunks,
	)

	barW := d.width - len(lineS) - len(lineE)

	chunkW := int(math.Floor(float64(barW) / float64(d.downloader.TotalChunks)))

	var bar string
	for i := 0; i < d.downloader.TotalChunks; i++ {
		c := d.downloader.Chunks[i]
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
	os.Stdout.Write([]byte(line + "\r"))
	os.Stdout.Sync()
}
