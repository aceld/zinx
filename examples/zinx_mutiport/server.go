package main

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_server/s_router"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
)

// Execute when creating a connection (创建连接的时候执行)
func DoConnectionBegin(conn ziface.IConnection) {
	zlog.Ins().InfoF("DoConnecionBegin is Called ...")

	//设置两个链接属性，在连接创建之后
	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://yuque.com/@aceld")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		zlog.Error(err)
	}
}

// Execute when connection lost (连接断开的时候执行)
func DoConnectionLost(conn ziface.IConnection) {
	// Query the Name and Home properties of the conn before destroying the connectio
	// 在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		zlog.Ins().InfoF("Conn Property Name = %v", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		zlog.Ins().InfoF("Conn Property Home = %v", home)
	}

	zlog.Ins().InfoF("Conn is Lost")
}

func main() {

	var i = 0

	for i < 2 {

		port := 8999 + i
		s := znet.NewUserConfServer(&zconf.Config{
			TCPPort: port,
			Name:    fmt.Sprintf("MyZinxServer-port:%d", port),
		})

		s.SetOnConnStart(DoConnectionBegin)
		s.SetOnConnStop(DoConnectionLost)

		s.AddRouter(100, &s_router.PingRouter{})
		s.AddRouter(1, &s_router.HelloZinxRouter{})

		go s.Serve()

		i++
	}

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
}
