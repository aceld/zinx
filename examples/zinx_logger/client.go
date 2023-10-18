package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gstones/zinx/examples/zinx_client/c_router"
	"github.com/gstones/zinx/ziface"
	"github.com/gstones/zinx/zlog"
	"github.com/gstones/zinx/znet"
)

func business(conn ziface.IConnection) {

	for {
		err := conn.SendMsg(1, []byte("Ping...[FromClient]"))
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

	conn.SetProperty("Name", "刘丹冰")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	go business(conn)
}

func DoClientConnectedLost(conn ziface.IConnection) {
	if name, err := conn.GetProperty("Name"); err == nil {
		zlog.Error("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		zlog.Error("Conn Property Home = ", home)
	}

	zlog.Debug("DoClientConnectedLost is Called ... ")
}

func main() {
	client := znet.NewClient("127.0.0.1", 8999)

	client.SetOnConnStart(DoClientConnectedBegin)
	client.SetOnConnStop(DoClientConnectedLost)

	client.AddRouter(0, &c_router.PingRouter{})

	client.Start()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
}
