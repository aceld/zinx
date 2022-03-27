package main

import (
	"github.com/aceld/zinx/zd"
)

func main() {
	server := zd.NewNode()
	server.Start()

	select {}
}
