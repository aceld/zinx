package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"time"
)

// User-defined heartbeat message processing method
// 用户自定义的心跳检测消息处理方法
func myHeartBeatMsg(conn ziface.IConnection) []byte {
	return []byte("heartbeat, I am server, I am alive")
}

// User-defined handling method for remote connection not alive.
// 用户自定义的远程连接不存活时的处理方法
func myOnRemoteNotAlive(conn ziface.IConnection) {
	fmt.Println("myOnRemoteNotAlive is Called, connID=", conn.GetConnID(), "remoteAddr = ", conn.RemoteAddr())
	//关闭链接
	conn.Stop()
}

// User-defined method for handling heartbeat messages (用户自定义的心跳检测消息处理方法)
type myHeartBeatRouter struct {
	znet.BaseRouter
}

func (r *myHeartBeatRouter) Handle(request ziface.IRequest) {
	fmt.Println("in MyHeartBeatRouter Handle, recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
}

func main() {
	s := znet.NewServer()

	myHeartBeatMsgID := 88888

	// Start heartbeating detection. (启动心跳检测)
	s.StartHeartBeatWithOption(1*time.Second, &ziface.HeartBeatOption{
		MakeMsg:          myHeartBeatMsg,
		OnRemoteNotAlive: myOnRemoteNotAlive,
		Router:           &myHeartBeatRouter{},
		HeartBeatMsgID:   uint32(myHeartBeatMsgID),
	})

	s.Serve()
}
