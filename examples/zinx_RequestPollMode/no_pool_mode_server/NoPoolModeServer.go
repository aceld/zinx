package main

import (
	"fmt"
	"time"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

// 如果不使用对象池模式则可以直接传递但是产生大量的 Request 对象

func NoPoll1(request ziface.IRequest) {
	request.Set("num", 1)
	go NoPoll2(request)
}

func NoPoll2(request ziface.IRequest) {
	time.Sleep(time.Second * 3)
	get, _ := request.Get("num")
	fmt.Printf("num:%v \n", get)

}

func NoPoll4(request ziface.IRequest) {
	// 非对象池模式永远不会影响原本的 Request
	request.Set("num", 3)
}

func main() {

	// 关闭 Request 对象池模式
	server := znet.NewUserConfServer(&zconf.Config{RouterSlicesMode: true, TCPPort: 8999, Host: "127.0.0.1", RequestPoolMode: false})
	server.AddRouterSlices(1, NoPoll1)
	server.AddRouterSlices(2, NoPoll4)
	server.Serve()
}
