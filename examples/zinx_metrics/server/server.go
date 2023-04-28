package main

import (
	"github.com/aceld/zinx/examples/zinx_server/s_router"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

func DoConnectionBegin(conn ziface.IConnection) {
	zlog.Ins().InfoF("DoConnecionBegin is Called ...")

	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		zlog.Error(err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	if name, err := conn.GetProperty("Name"); err == nil {
		zlog.Ins().InfoF("Conn Property Name = %v", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		zlog.Ins().InfoF("Conn Property Home = %v", home)
	}

	zlog.Ins().InfoF("Conn is Lost")
}

// usage:$  curl 0.0.0.0:20004/metrics
// to get Metrics
func main() {
	s := znet.NewServer()

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(100, &s_router.PingRouter{})
	s.AddRouter(1, &s_router.HelloZinxRouter{})

	s.Serve()
}
