package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/snikch/go-download/core"
)

var (
	resume = flag.Bool("resume", true, "Attempt to resume unfinished downloads")
	file   = flag.String("file", "", "URL of file to download")
	chunks = flag.Int("chunks", 6, "Number of chunks to download at once")
)

func main() {
	flag.Parse()
	if *file == "" {
		panic(fmt.Errorf("No file provided"))
	}

	client, err := core.NewRpcClient()
	if err != nil {
		panic(err)
	}

	var reply bool
	err = client.Call("RpcController.AddDownload", *file, &reply)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	fmt.Printf("Added %s: %s", file, reply)
}
