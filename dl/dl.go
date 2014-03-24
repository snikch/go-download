package main

import (
	"github.com/nsf/termbox-go"
	"github.com/snikch/go-download/core"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.SetCursor(0, 0)
	defer termbox.Close()

	settings, err := core.LoadSettings()
	if err != nil {
		panic(err)
	}

	// Load state
	// Start state manager
	manager := core.NewManager(settings)

	// Start rpc server
	err = core.StartRpcServer(manager)
	if err != nil {
		panic(err)
	}

	manager.AddUrl("http://mirror.cessen.com/blender.org/peach/trailer/trailer_iphone.m4v")
	manager.Start()
}
