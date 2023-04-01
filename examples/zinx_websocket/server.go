package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type PositionServerRouter struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PositionServerRouter) Handle(request ziface.IRequest) {
	fmt.Printf("data: %s\n", string(request.GetData()))
}

func main() {
	//创建一个server句柄
	s := znet.NewServer()

	//配置路由
	s.AddRouter(0, &PositionServerRouter{})

	//开启服务
	s.Serve()
}
