package main

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_client/c_router"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
	"time"
)

func business(conn ziface.IConnection) {

	for {
		err := conn.SendMsg(100, []byte("Ping...[FromClient]"))
		if err != nil {
			fmt.Println(err)
			zlog.Error(err)
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func DoClientConnectedBegin(conn ziface.IConnection) {
	zlog.Debug("DoConnecionBegin is Called ... ")

	//设置两个链接属性，在连接创建之后
	conn.SetProperty("Name", "刘丹冰Aceld")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	go business(conn)
}

func DoClientConnectedLost(conn ziface.IConnection) {
	if name, err := conn.GetProperty("Name"); err == nil {
		zlog.Debug("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		zlog.Debug("Conn Property Home = ", home)
	}

	zlog.Debug("DoClientConnectedLost is Called ... ")
}

func main() {
	client := znet.NewClient("127.0.0.1", 8999)

	client.SetOnConnStart(DoClientConnectedBegin)
	client.SetOnConnStop(DoClientConnectedLost)

	client.AddRouter(2, &c_router.PingRouter{})
	client.AddRouter(3, &c_router.HelloRouter{})

	client.Start()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)

	client.Stop()
	time.Sleep(time.Second * 2)
}
