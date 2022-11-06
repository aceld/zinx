package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

// test handleRouter
func Handle1(router ziface.IRouter, req ziface.IRequest) {
	fmt.Println("1")
	if err := req.GetConnection().SendMsg(0, []byte("test1")); err != nil {
		fmt.Println(err)
	}
}

// test handleRouter
func Handle2(router ziface.IRouter, req ziface.IRequest) {
	fmt.Println("2")
	if err := req.GetConnection().SendMsg(0, []byte("test2")); err != nil {
		fmt.Println(err)
	}
}

// test handleRouter
func Handle3(router ziface.IRouter, req ziface.IRequest) {
	fmt.Println("3")
	if err := req.GetConnection().SendMsg(0, []byte("test3")); err != nil {
		fmt.Println(err)
	}
}

// test handleRouter
func Handle4(router ziface.IRouter, req ziface.IRequest) {
	fmt.Println("4")
	if err := req.GetConnection().SendMsg(0, []byte("test4")); err != nil {
		fmt.Println(err)
	}

}

func main() {
	s := znet.NewServer()
	s.AddRouter(1, Handle1, Handle2, Handle3, Handle4)

	s.Serve()
}
