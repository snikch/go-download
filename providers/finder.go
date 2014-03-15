package providers

import "net/url"

type Finder struct {
	providers []Provider
}

func newFinder() Finder {
	f := Finder{make([]Provider, 0, 1)}
	return f
}

// Returns a finder that matches the given url
func (f *Finder) providerFor(url *url.URL) Provider {
	return &baseProvider{
		"Default",
		1024 * 256,
		4,
	}
}

func (f *Finder) Register(p Provider) {
	f.providers = append(f.providers, p)
}
