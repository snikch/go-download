package core

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
)

type Downloadable struct {
	Downloaded ByteSize
	Size       ByteSize

	client   *http.Client
	Resource *Resource
}

func (d *Downloadable) Progress() int {
	return int(float64(100) * float64(d.Downloaded) / float64(d.Size))
}

type Downloader struct {
	Downloadable
	io.WriteCloser

	Chunks         []*Chunk
	CompleteChunks int
	MaxChunks      int
	SpeedMonitor   *SpeedMonitor
	TotalChunks    int

	nextChunkIndex int
	progressChan   chan *Downloader
	resume         bool
}

func NewDownloader(url string, updater chan *Downloader) (d Downloader, err error) {
	resource, err := NewResource(url)
	if err != nil {
		return
	}
	d = Downloader{progressChan: updater}
	d.Resource = &resource
	return
}

func (d *Downloader) Start(resume bool) {
	d.MaxChunks = d.Resource.Hoster.MaxChunks()
	d.resume = resume
	go d.startChunks()
}

func (d *Downloader) startChunks() (err error) {
	// Get the expected filesize
	client := &http.Client{}
	req, err := http.NewRequest("HEAD", d.Resource.Url.String(), nil)
	if err != nil {
		return
	}
	//req.Header.Add("")
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	d.Size = ByteSize(resp.ContentLength)

	// Create an appropriate amount of chunks
	//	 downloader.chunks || (resource.size / Hoster.chunkSize)

	// How many chunks are we gonna need?
	d.TotalChunks = int(
		math.Ceil(float64(d.Size) / float64(d.Resource.Hoster.ChunkSize())))

	if d.TotalChunks < d.MaxChunks {
		d.MaxChunks = d.TotalChunks
	}
	d.nextChunkIndex = d.MaxChunks

	d.Chunks = make([]*Chunk, d.TotalChunks)

	var start ByteSize
	chunkProgress := make(chan *Chunk)
	chunkComplete := make(chan *Chunk)
	amountDownloaded := make(chan ByteSize)

	// Create a speed monitor and start monitoring
	d.SpeedMonitor = NewSpeedMonitor(amountDownloaded)
	go d.SpeedMonitor.Start()

	for i := 0; i < d.TotalChunks; i++ {
		end := start + ByteSize(d.Resource.Hoster.ChunkSize())
		if end > d.Size {
			end = d.Size
		}
		chunk, err := NewChunk(d.Resource, int(i), start, end-1, d.resume)
		if err != nil {
			panic(err)
		}
		chunk.progress = chunkProgress
		chunk.complete = chunkComplete
		chunk.client = client
		d.Chunks[i] = &chunk
		if i < d.MaxChunks {
			go chunk.Download()
		}
		start = end
	}

	for {
		select {
		case _ = <-chunkProgress:
			d.updateProgress()
			d.progressChan <- d
			amountDownloaded <- d.Downloaded

		case _ = <-chunkComplete:
			d.CompleteChunks += 1
			if d.CompleteChunks < d.TotalChunks {
				d.startNextChunk()
			} else {
				err = d.createFinalFile()
				if err != nil {
					panic(err)
				}
				d.Close()
			}
		}
	}

	return nil
}

func (d *Downloader) startNextChunk() error {
	if d.nextChunkIndex >= d.TotalChunks {
		return AllChunksDownloadingError{d.nextChunkIndex, d.TotalChunks}
	}
	go d.Chunks[d.nextChunkIndex].Download()
	d.nextChunkIndex += 1
	return nil
}

// Concatenate all the chunk files into the final file
func (d *Downloader) createFinalFile() (err error) {

	if err != nil {
		return
	}

	defer d.Close()

	for i := 0; i < len(d.Chunks); i++ {
		filename := d.Chunks[i].DestinationFile()
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(d, file)
		if err != nil {
			return err
		}
		file.Close()
		os.Remove(filename)
	}

	return
}

func (d *Downloader) Write(p []byte) (n int, err error) {
	if d.WriteCloser == nil {
		filename := d.DestinationFile()
		os.Remove(filename)
		dst, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return 0, err
		}
		d.WriteCloser = dst
	}
	n, err = d.WriteCloser.Write(p)
	return

}

func (d *Downloader) DestinationFile() string {

	return fmt.Sprintf("%s/%s", d.Resource.Name, d.Resource.Name)
}

func (d *Downloader) Close() {
	d.SpeedMonitor.Stop()
	d.WriteCloser.Close()
}

// Move to main package
type AllChunksDownloadingError struct {
	attempted int
	total     int
}

func (e AllChunksDownloadingError) Error() string {
	return fmt.Sprintf(
		"Attempted to download chunk %d, but only %d chunks exist",
		e.attempted,
		e.total,
	)
}
func (d *Downloader) AreAllChunksDownloading(err error) (ok bool) {
	_, ok = err.(AllChunksDownloadingError)
	return
}

func (d *Downloader) updateProgress() {
	var downloaded ByteSize
	for _, chunk := range d.Chunks {
		downloaded += chunk.Downloaded
	}
	d.Downloaded = downloaded
}
