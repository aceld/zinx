package main

import (
	"github.com/aceld/zinx/v3/examples/zinx_server/s_router"
	"github.com/aceld/zinx/v3/zconf"
	"github.com/aceld/zinx/v3/znet"
)

func main() {
	// Set up as WebSocket before starting. (在启动之前设置为 websocket)
	zconf.GlobalObject.Mode = ""
	zconf.GlobalObject.LogFile = ""

	s := znet.NewServer()

	s.AddRouter(100, &s_router.PingRouter{})
	s.AddRouter(1, &s_router.HelloZinxRouter{})

	s.Serve()
}
