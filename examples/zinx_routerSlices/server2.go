package main

import (
	"fmt"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"time"
)

// WorkerPoolGoroutineNums 实例
// 模拟阻塞场景，f1函数会阻塞500ms 如果设置WorkerPoolGoroutineNums数量大于1就可以让同一个消息队列中的chan能够执行而非阻塞
// 如果默认场景则一定会等f1阻塞执行完成 整个chan才会接着读取执行任务
// 警告！WorkerPoolGoroutineNums>1会导致同一chan中的消息执行顺序不可控！注意应该使用在不需要注意消息处理顺序的场景
func f1(request ziface.IRequest) {
	time.Sleep(500 * time.Millisecond)
	fmt.Println("我是阻塞函数")

}
func f2(request ziface.IRequest) {
	fmt.Println("Test2")
}
func f3(request ziface.IRequest) {
	fmt.Println("Test3")
}

func main() {

	server := znet.NewUserConfServer(&zconf.Config{RouterSlicesMode: true, TCPPort: 8999, Host: "127.0.0.1", WorkerPoolGoroutineNums: 1})
	//配置大于1
	//server := znet.NewUserConfServer(&zconf.Config{RouterSlicesMode: true, TCPPort: 8999, Host: "127.0.0.1", WorkerPoolGoroutineNums: 2})
	server.Use(znet.RouterTime)
	server.AddRouterSlices(1, f1)
	server.AddRouterSlices(2, f2)
	server.AddRouterSlices(3, f3)
	server.Serve()
}
