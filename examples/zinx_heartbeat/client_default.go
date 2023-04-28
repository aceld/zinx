package main

import (
	"github.com/aceld/zinx/znet"
	"time"
)

func main() {
	client := znet.NewClient("127.0.0.1", 8999)

	// Start heartbeating detection.
	client.StartHeartBeat(3 * time.Second)

	client.Start()

	// wait
	select {}
}
