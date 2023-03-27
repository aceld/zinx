package main

import (
	"github.com/aceld/zinx/znet"
	"time"
)

func main() {
	//创建一个Client句柄，使用Zinx的API
	client := znet.NewClient("127.0.0.1", 8999)

	//启动心跳检测
	client.StartHeartBeat(3 * time.Second)

	//启动客户端client
	client.Start()

	// wait
	select {}
}
