package main

import (
	"github.com/gstones/zinx/znet"
)

func main() {
	s := znet.NewServer()

	// Start heartbeating detection.
	s.Serve()
}
