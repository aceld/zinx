package main

import (
	"dart/pb"
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
	"time"
)

func main() {
	server := znet.NewServer()
	server.AddRouter(uint32(pb.MessageID_First), &HelloRouter{})
	server.Serve()
}

type HelloRouter struct {
	znet.BaseRouter
}

func (h *HelloRouter) Handle(request ziface.IRequest) {
	req := &pb.HelloRequest{}
	if err := proto.Unmarshal(request.GetData(), req); err != nil {
		zlog.Debugf("proto unmarshal error, %+v", err)
		return
	}
	exits := make(chan struct{}, 1)
	resp := &pb.HelloResponse{}
	go func() {
		for {
			resp.Greeter = fmt.Sprintf("Hello %s, now：%s.", req.Name, time.Now().Format(time.RFC3339))
			bytes, err := proto.Marshal(resp)
			if err != nil {
				exits <- struct{}{}
				return
			}
			if err := request.GetConnection().SendBuffMsg(request.GetMsgID(), bytes); err != nil {
				exits <- struct{}{}
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()
	select {
	case <-exits:
		return
	}
}
