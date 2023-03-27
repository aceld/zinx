package main

import (
	"fmt"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
	"time"
)

func wait() {
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
}

func main() {
	//创建一个Client句柄，使用Zinx的API
	client := znet.NewClient("127.0.0.1", 8999)

	//启动心跳检测
	client.StartHeartBeat(3 * time.Second)

	//启动客户端client
	client.Start()

	wait()
}
