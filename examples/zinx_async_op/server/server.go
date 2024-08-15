package main

import (
	"github.com/aceld/zinx/examples/zinx_async_op/router"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

func OnConnectionAdd(conn ziface.IConnection) {
	zlog.Debug("zinx_async_op OnConnectionAdd ===>")
}

func OnConnectionLost(conn ziface.IConnection) {
	zlog.Debug("zinx_async_op OnConnectionLost ===>")
}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	s.AddRouter(1, &router.LoginRouter{})

	s.Serve()
}
