package main

import (
	"github.com/aceld/zinx/examples/zinx_server/s_router"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/znet"
)

func main() {
	// 在启动之前设置为 websocket
	zconf.GlobalObject.Mode = ""
	zconf.GlobalObject.LogFile = ""

	//创建一个server句柄
	s := znet.NewServer()
	//配置路由
	s.AddRouter(100, &s_router.PingRouter{})
	s.AddRouter(1, &s_router.HelloZinxRouter{})

	//开启服务
	s.Serve()
}
