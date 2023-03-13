package main

import (
	"github.com/aceld/zinx/examples/zinx_decoder/bili"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"math"
)

func DoConnectionBegin(conn ziface.IConnection) {
}

func DoConnectionLost(conn ziface.IConnection) {
}

func main() {
	server := znet.NewServer(func(s *znet.Server) {
		s.Port = 9090
		s.LengthField = ziface.LengthField{
			MaxFrameLength:      math.MaxUint8 + 4,
			LengthFieldOffset:   2,
			LengthFieldLength:   1,
			LengthAdjustment:    2,
			InitialBytesToStrip: 0,
		}
	})
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLost)
	coder := bili.HtlvcrcInterceptor{}
	//server.AddInterceptor(coder.GetDecoder())
	server.AddInterceptor(&coder)
	server.AddRouter(0x10, &bili.Data0x10Router{})
	server.AddRouter(0x13, &bili.Data0x13Router{})
	server.AddRouter(0x14, &bili.Data0x14Router{})
	server.AddRouter(0x15, &bili.Data0x15Router{})
	server.AddRouter(0x16, &bili.Data0x16Router{})
	server.Serve()

}
