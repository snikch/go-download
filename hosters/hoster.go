package hosters

import "net/url"

type Hoster interface {
	Name() string
	ChunkSize() int64
	MaxChunks() int
	Match(url *url.URL) bool
	URLPreflight(url *url.URL) error
}

var (
	finder = newFinder()
)
