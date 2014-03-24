package core

import (
	"net/url"
	"path"

	"github.com/snikch/go-download/hosters"
)

type Resource struct {
	Name   string
	Hoster hosters.Hoster
	Url    *url.URL
	size   ByteSize
}

func NewResource(address string) (res Resource, err error) {
	url, err := url.Parse(address)
	if err != nil {
		return
	}

	h, err := hosters.FindHosterForUrl(url)
	if err != nil {
		return
	}

	res = Resource{
		Name:   path.Base(url.Path),
		Url:    url,
		Hoster: h,
	}
	return
}
