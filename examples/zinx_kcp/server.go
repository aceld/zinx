package main

import (
	"errors"
	"fmt"
	"github.com/aceld/zinx/zconf"
	"time"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type TestRouter struct {
	znet.BaseRouter
}

var dealTimes = 0

// PreHandle -
func (t *TestRouter) PreHandle(req ziface.IRequest) {
	start := time.Now()

	fmt.Println("--> Call PreHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test1")); err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(start)
	fmt.Println("cost time：", elapsed)
}

// Handle -
func (t *TestRouter) Handle(req ziface.IRequest) {
	fmt.Println("--> Call Handle")

	if err := Err(); err != nil {
		req.Abort()
		fmt.Println("Insufficient permission")
	}

	dealTimes++
	req.GetConnection().AddCloseCallback(nil, nil, func() {
		fmt.Println("run close callback")
	})

	if err := req.GetConnection().SendMsg(0, []byte("test2")); err != nil {
		fmt.Println(err)
	}

	if dealTimes == 5 {
		req.GetConnection().Stop()
	}

	time.Sleep(1 * time.Millisecond)
}

// PostHandle -
func (t *TestRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("--> Call PostHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test3")); err != nil {
		fmt.Println(err)
	}
}

func Err() error {
	//Specific Business Operation (具体业务操作)
	return errors.New("Test")
}

func main() {
	s := znet.NewUserConfServer(&zconf.Config{
		Mode:          "kcp",
		KcpPort:       7777,
		KcpRecvWindow: 128,
		KcpSendWindow: 128,
		KcpStreamMode: true,
		KcpACKNoDelay: false,
		LogDir:        "./",
		LogFile:       "test.log",
	})
	s.AddRouter(1, &TestRouter{})
	s.SetOnConnStart(func(conn ziface.IConnection) {
		fmt.Println("--> OnConnStart")
	})
	s.SetOnConnStop(func(conn ziface.IConnection) {
		fmt.Println("--> OnConnStop")
	})
	s.Serve()
}
