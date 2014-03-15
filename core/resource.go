package core

import (
	"net/url"
	"path"

	"github.com/snikch/go-download/providers"
)

type Resource struct {
	Url      *url.URL
	Name     string
	size     ByteSize
	provider providers.Provider
}

func NewResource(address string) (res Resource, err error) {
	url, err := url.Parse(address)
	if err != nil {
		return
	}
	res = Resource{
		Name:     path.Base(url.Path),
		Url:      url,
		provider: providers.NewProvider(url),
	}
	return
}
