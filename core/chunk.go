package core

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
)

type Chunk struct {
	io.ReadCloser
	Downloadable

	Percent int

	end      ByteSize
	complete chan *Chunk
	progress chan *Chunk
	request  *http.Request
	response *http.Response
	resume   bool
	start    ByteSize
	state    string
	store    *ChunkStore
	index    int
}

func NewChunk(res *Resource, index int, start, end ByteSize, resume bool) (c Chunk, err error) {
	c = Chunk{
		end:    end,
		index:  index,
		resume: resume,
		start:  start,
		state:  "new",
	}

	c.Resource = res

	err = c.Setup()

	return
}
func (c *Chunk) Setup() (err error) {
	if _, err := os.Stat(c.DestinationFolder()); os.IsNotExist(err) {
		err := os.MkdirAll(c.DestinationFolder(), 0700)
		if err != nil {
			return err
		}
	}

	c.Size = c.end - c.start + 1

	c.store = NewChunkStore(c.DestinationFile())
	return
}

func (c *Chunk) Download() {
	err := c.Dl()
	if err != nil {
		panic(err)
	}
}

func (c *Chunk) Dl() (err error) {
	if c.state == "downloading" {
		panic(fmt.Errorf("Attempting to download an already downloading chunk"))
	}
	c.state = "downloading"

	startByte, err := c.initializeRequest()
	if c.IsChunkCompleteError(err) {
		c.Close()
		return nil
	}
	if err != nil {
		return
	}

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c.response, err = client.Do(c.request)
	if err != nil {
		return
	}

	if c.response.StatusCode != 206 {
		return fmt.Errorf(
			"Expected 206 Partial Content, received %d instead",
			c.response.StatusCode,
		)
	}
	// Check for size? Ensure bytes header accepted?
	responseLength := c.response.ContentLength
	expectedLength := int64(c.end) - int64(startByte) + 1
	if responseLength != expectedLength {
		return fmt.Errorf(
			"Expected response content length to be %d, but was %d",
			expectedLength,
			responseLength,
		)
	}
	c.ReadCloser = c.response.Body

	defer c.Close()

	io.Copy(c, c)

	return nil
}

func (c *Chunk) initializeRequest() (startByte int, err error) {
	// If we're resuming, we don't need the entire chunk
	if c.resume == true {
		existingBytes, err := c.store.Size()
		if err != nil {
			return startByte, err
		}
		startByte = int(c.start + existingBytes)
		c.Downloaded = existingBytes
		if c.Downloaded == c.Size {
			return startByte, ChunkCompleteError{}
		}
		if c.Downloaded >= c.Size {
			// Corrupt?
			c.resume = false
		}
	}
	if c.resume == false {
		c.store.Remove()
		startByte = int(c.start)
	}

	err = c.store.Open()
	if err != nil {
		return
	}

	c.request, err = http.NewRequest("GET", c.Resource.Url.String(), nil)
	if err != nil {
		return
	}

	rangeS := fmt.Sprintf(
		"bytes=%d-%d",
		startByte,
		int(c.end),
	)

	c.request.Header.Add(
		"Range",
		rangeS,
	)
	return
}

func (c *Chunk) DestinationFile() string {
	return fmt.Sprintf(
		"%s/%d-%d-%d",
		c.DestinationFolder(),
		c.index,
		int(c.start),
		int(c.end),
	)
}

func (c *Chunk) DestinationFolder() string {
	return fmt.Sprintf("%s", c.Resource.Name)
}

func (c *Chunk) Read(p []byte) (n int, err error) {
	n, err = c.ReadCloser.Read(p)
	if err != nil {
		return
	}
	c.Downloaded += ByteSize(n)
	if c.Downloaded > c.Size {
		return n, fmt.Errorf(
			"Downloaded %s, which is bigger than expected size of %s",
			c.Downloaded,
			c.Size,
		)
	}
	c.updateProgress()
	return
}

func (c *Chunk) updateProgress() {
	if c.Size > 0 {
		p := int(math.Ceil(float64(c.Downloaded) / float64(c.Size) * float64(100)))
		c.Percent = p
		c.progress <- c
	}
}

func (c *Chunk) Write(p []byte) (n int, err error) {
	n, err = c.store.Write(p)
	return
}

func (c *Chunk) Close() (err error) {
	c.state = "downloaded"
	if c.ReadCloser != nil {
		err = c.ReadCloser.Close()
		if err != nil {
			return
		}
	}
	c.updateProgress()
	c.complete <- c
	return
}

type ChunkCompleteError struct {
	c *Chunk
}

func (e ChunkCompleteError) Error() string {
	return fmt.Sprintf("Chunk %d is already complete", e.c.index)
}

func (c *Chunk) IsChunkCompleteError(err error) (ok bool) {
	_, ok = err.(ChunkCompleteError)
	return
}
