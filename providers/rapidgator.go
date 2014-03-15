package providers

import (
	"net/url"
	"regexp"
)

type rapidgator struct {
	baseProvider
}

func (r *rapidgator) Match(url *url.URL) (isMatch bool) {
	isMatch, err := regexp.Match(
		"^https?://rapidgator.com",
		[]byte(url.String()),
	)

	if err != nil {
		isMatch = false
	}

	return
}

func init() {
	finder.Register(&rapidgator{})
}
