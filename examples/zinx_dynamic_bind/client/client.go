package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

const (
	PingType = 1
	PongType = 2
)

// ping response router
type PongRouter struct {
	znet.BaseRouter
	client string
}

// Hash 工作模式下，需要等待接受到client1的pong后，才会收到client2和client3的pong
// DynamicBind工作模式下，client2, client3 都会立马收到pong, 但client1的pong会被阻塞十秒后才收到
func (p *PongRouter) Handle(request ziface.IRequest) {
	//read server pong data
	zlog.Infof("---------client:%s, recv from server:%s, msgId=%d, data=%s ----------\n",
		p.client, request.GetConnection().RemoteAddr(), request.GetMsgID(), string(request.GetData()))
}

func onClient1Start(conn ziface.IConnection) {
	zlog.Infof("client1 connection start, %s->%s\n", conn.LocalAddrString(), conn.RemoteAddrString())
	//send ping
	err := conn.SendMsg(PingType, []byte("Ping From Client1"))
	if err != nil {
		zlog.Error(err)
	}
}

func onClient2Start(conn ziface.IConnection) {
	zlog.Infof("client2 connection start, %s->%s\n", conn.LocalAddrString(), conn.RemoteAddrString())
	//send ping
	err := conn.SendMsg(PingType, []byte("Ping From Client2"))
	if err != nil {
		zlog.Error(err)
	}
}

func onClient3Start(conn ziface.IConnection) {
	zlog.Infof("client3 connection start, %s->%s\n", conn.LocalAddrString(), conn.RemoteAddrString())
	//send ping
	err := conn.SendMsg(PingType, []byte("Ping From Client3"))
	if err != nil {
		zlog.Error(err)
	}
}

func main() {
	//Create a client client
	client1 := znet.NewClient("127.0.0.1", 8999)
	client1.SetOnConnStart(onClient1Start)
	client1.AddRouter(PongType, &PongRouter{client: "client1"})
	client1.Start()

	time.Sleep(time.Second)

	client2 := znet.NewClient("127.0.0.1", 8999)
	client2.SetOnConnStart(onClient2Start)
	client2.AddRouter(PongType, &PongRouter{client: "client2"})
	client2.Start()

	time.Sleep(time.Second)

	client3 := znet.NewClient("127.0.0.1", 8999)
	client3.SetOnConnStart(onClient3Start)
	client3.AddRouter(PongType, &PongRouter{client: "client3"})
	client3.Start()

	//Prevent the process from exiting, waiting for an interrupt signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	client1.Stop()
	client2.Stop()
	client3.Stop()

	time.Sleep(time.Second)
}
