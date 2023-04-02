package main

import (
	"github.com/aceld/zinx/examples/zinx_server/s_router"
	"github.com/aceld/zinx/znet"
)

func main() {
	//znet.NewUserConfServer(&zconf.Config{
	//	TCPPort: 9999,
	//})
	//创建一个server句柄
	s := znet.NewServer()

	//配置路由
	s.AddRouter(100, &s_router.PingRouter{})
	s.AddRouter(1, &s_router.HelloZinxRouter{})

	//开启服务
	s.Serve()
}
