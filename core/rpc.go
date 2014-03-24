package core

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

type RpcController struct {
	manager *Manager
}

const (
	RPC_PORT = 7696
)

func (c *RpcController) AddDownload(url string, ok *bool) (err error) {
	err = c.manager.AddUrl(url)
	if err != nil {
		return
	}
	*ok = true
	return nil
}

func StartRpcServer(m *Manager) (err error) {
	c := &RpcController{m}
	rpc.Register(c)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", RPC_PORT))
	if err != nil {
		return
	}
	go http.Serve(l, nil)
	return nil
}

func NewRpcClient() (client *rpc.Client, err error) {
	return rpc.DialHTTP("tcp", fmt.Sprintf(":%d", RPC_PORT))
}
