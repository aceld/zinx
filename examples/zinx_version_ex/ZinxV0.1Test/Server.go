package main

import (
	"github.com/aceld/zinx/znet"
)

//Server 模块的测试函数
func main() {

	/*
		服务端测试
	*/
	//1 创建一个server 句柄 s
	// s := znet.NewServer("[zinx V0.1]")
	s := znet.NewServer()

	//2 开启服务
	s.Serve()
}
