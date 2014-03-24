package hosters

import (
	"net/url"
	"reflect"
	"regexp"
)

type baseHoster struct {
	name      string
	chunkSize int64
	maxChunks int
}

const (
	DEFAULT_CHUNK_SIZE = 1024 * 1024
	DEFAULT_MAX_CHUNKS = 6
)

// If no name is explicitly set, use reflection to set/get one
func (p *baseHoster) Name() string {
	if p.name != "" {
		return p.name
	}
	t := reflect.TypeOf(p)
	p.name = t.Name()
	return p.name
}

func (p *baseHoster) ChunkSize() int64 {
	if p.chunkSize == 0 {
		p.chunkSize = DEFAULT_CHUNK_SIZE
	}
	return p.chunkSize
}

func (p *baseHoster) MaxChunks() int {
	if p.maxChunks == 0 {
		p.maxChunks = SOMETHING
	}
	return p.maxChunks
}

func (p *baseHoster) URLPreflight(url *url.URL) (err error) {
	return
}

// Only handle http(s) requests in the base
func (p *baseHoster) Match(url *url.URL) (isMatch bool) {
	isMatch, err := regexp.Match(
		"^https?",
		[]byte(url.String()),
	)

	if err != nil {
		isMatch = false
	}

	return
}

func init() {
	finder.Register(&baseHoster{
		name: "HTTP",
	})
}
