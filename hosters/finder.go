package hosters

import (
	"fmt"
	"net/url"
)

type Finder struct {
	hosters []Hoster
}

func newFinder() Finder {
	f := Finder{
		make([]Hoster, 0, 1),
	}
	return f
}

// Returns a finder that matches the given url
func (f *Finder) hosterForUrl(url *url.URL) (p Hoster, err error) {
	for _, p = range f.hosters {
		if p.Match(url) {
			return
		}
	}
	return nil, fmt.Errorf("No hoster matched %s", url.String())
}

func (f *Finder) hosterForName(name string) (p Hoster, e error) {
	for _, p = range f.hosters {
		if name == p.Name() {
			return
		}
	}
	return nil, fmt.Errorf("No hoster named %s", name)
}

func (f *Finder) Register(p Hoster) {
	f.hosters = append(f.hosters, p)
}

func FindHosterForUrl(url *url.URL) (Hoster, error) {
	return finder.hosterForUrl(url)
}
func FindHosterByName(name string) (Hoster, error) {
	return finder.hosterForName(name)
}
