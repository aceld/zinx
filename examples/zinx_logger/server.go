package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"time"
)

type TestRouter struct {
	znet.BaseRouter
}

// PreHandle -
func (t *TestRouter) PreHandle(req ziface.IRequest) {
	start := time.Now()

	fmt.Println("--> Call PreHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test1")); err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(start)
	fmt.Println("elapsed：", elapsed)
}

// Handle -
func (t *TestRouter) Handle(req ziface.IRequest) {
	fmt.Println("--> Call Handle")

	if err := req.GetConnection().SendMsg(0, []byte("test2")); err != nil {
		fmt.Println(err)
	}
}

// PostHandle -
func (t *TestRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("--> Call PostHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test3")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(1, &TestRouter{})
	zlog.SetLogger(new(MyLogger))
	s.Serve()
}
