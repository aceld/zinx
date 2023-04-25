package main

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"time"
)

func main() {
	client := znet.NewClient("127.0.0.1", 8999)

	client.SetOnConnStart(func(connection ziface.IConnection) {
		_ = connection.SendMsg(1, []byte("hello zinx"))
	})

	client.Start()

	time.Sleep(time.Second)
}
