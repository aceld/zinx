package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zinx_app_demo/mmo_game/pb"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
)

type PositionServerRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (this *PositionServerRouter) Handle(request ziface.IRequest) {

	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Position Unmarshal error ", err, " data = ", request.GetData())
		return
	}

	fmt.Printf("recv from client : msgId=%+v, data=%+v\n", request.GetMsgID(), msg)

	msg.X += 1
	msg.Y += 1
	msg.Z += 1
	msg.V += 1

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("proto Marshal error = ", err, " msg = ", msg)
		return
	}

	err = request.GetConnection().SendMsg(0, data)

	if err != nil {
		zlog.Error(err)
	}
}

func main() {
	s := znet.NewServer()

	s.AddRouter(0, &PositionServerRouter{})

	s.Serve()
}
