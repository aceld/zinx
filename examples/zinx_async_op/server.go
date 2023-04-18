package main

import (
	"github.com/aceld/zinx/examples/zinx_async_op/router"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

// OnConnectionAdd 当客户端建立连接的时候的hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	zlog.Debug("zinx_async_op OnConnectionAdd ===>")
}

// OnConnectionLost 当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	zlog.Debug("zinx_async_op OnConnectionLost ===>")
}

func main() {
	//创建服务器句柄
	s := znet.NewServer()

	//注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	//注册路由
	s.AddRouter(1, &router.LoginRouter{})

	//启动服务
	s.Serve()
}
