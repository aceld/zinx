package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"time"
)

type TestRouter struct {
	znet.BaseRouter
}

// Handle -
func (t *TestRouter) Handle(req ziface.IRequest) {
	fmt.Println("--> Call Handle, reveived msg: ", string(req.GetData()), " msgID: ", req.GetMsgID(), " connID: ", req.GetConnection().GetConnID())

	if err := req.GetConnection().SendMsg(0, []byte("hello i am server")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer()

	s.AddRouter(1, &TestRouter{})

	//启动心跳检测
	s.StartHeartBeat(5 * time.Second)

	s.Serve()
}
