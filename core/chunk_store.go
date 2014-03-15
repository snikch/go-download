package core

import (
	"io"
	"os"
)

type ChunkStore struct {
	io.ReadWriteCloser
	location string
}

func NewChunkStore(location string) *ChunkStore {
	return &ChunkStore{location: location}
}

func (cs *ChunkStore) Open() (err error) {
	cs.ReadWriteCloser, err = os.OpenFile(
		cs.location,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0644,
	)
	return
}

func (cs *ChunkStore) Remove() (err error) {
	err = os.Remove(cs.location)
	return
}

func (cs *ChunkStore) Size() (b ByteSize, err error) {
	fi, err := os.Stat(cs.location)
	if os.IsNotExist(err) {
		return b, nil
	}
	if err != nil {
		return
	}
	b = ByteSize(fi.Size())
	return
}
