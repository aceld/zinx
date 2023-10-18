package main

import (
	"github.com/gstones/zinx/examples/zinx_decoder/bili/router"
	"github.com/gstones/zinx/zdecoder"
	"github.com/gstones/zinx/ziface"
	"github.com/gstones/zinx/znet"
)

func DoConnectionBegin(conn ziface.IConnection) {
}

func DoConnectionLost(conn ziface.IConnection) {
}

func main() {
	server := znet.NewServer(func(s *znet.Server) {
		s.Port = 9090
		/*
			s.LengthField = ziface.LengthField{
				MaxFrameLength:      math.MaxUint8 + 4,
				LengthFieldOffset:   2,
				LengthFieldLength:   1,
				LengthAdjustment:    2,
				InitialBytesToStrip: 0,
			}
		*/
	})
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLost)
	server.AddInterceptor(zdecoder.NewHTLVCRCDecoder())
	server.AddRouter(0x10, &router.Data0x10Router{})
	server.AddRouter(0x13, &router.Data0x13Router{})
	server.AddRouter(0x14, &router.Data0x14Router{})
	server.AddRouter(0x15, &router.Data0x15Router{})
	server.AddRouter(0x16, &router.Data0x16Router{})
	server.Serve()
}
