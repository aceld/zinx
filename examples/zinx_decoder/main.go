package main

import (
	"github.com/aceld/zinx/examples/zinx_decoder/bili"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

func DoConnectionBegin(conn ziface.IConnection) {
}

func DoConnectionLost(conn ziface.IConnection) {
}

func main() {
	server := znet.NewServer(func(s *znet.Server) {
		s.Port = 9090
	})
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLost)
	coder := bili.HtlvcrcInterceptor{}
	server.AddInterceptor(coder.GetDecoder())
	server.AddInterceptor(&coder)
	server.AddRouter(0x10, &bili.Data0x10Router{})
	server.AddRouter(0x13, &bili.Data0x13Router{})
	server.AddRouter(0x14, &bili.Data0x14Router{})
	server.AddRouter(0x15, &bili.Data0x15Router{})
	server.AddRouter(0x16, &bili.Data0x16Router{})
	server.Serve()

	//arr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	//fmt.Println(arr[0])
	//fmt.Println(arr[1])
	//fmt.Println(arr[2])
	//fmt.Println(arr[3 : len(arr)-2])
	//fmt.Println(arr[len(arr)-2 : len(arr)])
	//fmt.Println(arr)
}
