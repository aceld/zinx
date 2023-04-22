package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"time"
)

func DefaultTest1(request ziface.IRequest) {
	fmt.Println("test1")
}
func DefaultTest2(request ziface.IRequest) {
	time.Sleep(1)
	panic("test")
}

func main() {
	s := znet.NewDefaultRouterSlicesServer()
	s.AddRouterSlices(1, DefaultTest1, DefaultTest2)
	s.Serve()
}
