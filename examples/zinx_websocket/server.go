package main

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_server/s_router"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

func main() {
	//znet.NewUserConfServer(&zconf.Config{
	//	TCPPort: 9999,
	//})
	//创建一个server句柄
	zconf.GlobalObject.Mode = "all"
	s := znet.NewServer()
	s.SetOnConnStart(func(connection ziface.IConnection) {
		fmt.Println("SetOnConnStart")
	})
	//配置路由
	s.AddRouter(100, &s_router.PingRouter{})
	s.AddRouter(1, &s_router.HelloZinxRouter{})

	//开启服务
	s.Serve()
}
