package providers

import (
	"net/url"
	"reflect"
)

type Provider interface {
	Name() string
	ChunkSize() int64
	MaxChunks() int
	Match(url *url.URL) bool
}

type baseProvider struct {
	name      string
	chunkSize int64
	maxChunks int
}

// If no name is explicitly set, use reflection to set/get one
func (p *baseProvider) Name() string {
	if p.name != "" {
		return p.name
	}
	t := reflect.TypeOf(p)
	p.name = t.Name()
	return p.name
}

func (p *baseProvider) ChunkSize() int64 {
	return p.chunkSize
}

func (p *baseProvider) MaxChunks() int {
	return p.maxChunks
}

func (p *baseProvider) Match(url *url.URL) (isMatch bool) {
	return true
}

var (
	finder = newFinder()
)

func NewProvider(url *url.URL) Provider {
	return finder.providerFor(url)
}
