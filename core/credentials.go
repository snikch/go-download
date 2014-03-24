package core

import "net/http"

type Credentials interface {
	Authenticate(c *http.Client, req *http.Request) error
}

type BasicAuthCredentials struct {
	user     string
	password string
}

func (c *BasicAuthCredentials) Authenticate(client *http.Client, req *http.Request) (err error) {
	return
}
