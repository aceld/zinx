package main

import (
	"fmt"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

//Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func main() {
	//创建一个server句柄
	// s := znet.NewServer("[zinx V0.3]")
	s := znet.NewServer()

	// s.AddRouter(&PingRouter{})
	s.AddRouter(3, &PingRouter{})
	//2 开启服务
	s.Serve()
}
