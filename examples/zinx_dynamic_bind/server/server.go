package main

import (
	"sync/atomic"
	"time"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

func OnConnectionAdd(conn ziface.IConnection) {
	zlog.Debug("OnConnectionAdd:", conn.GetConnection().RemoteAddr())
}

func OnConnectionLost(conn ziface.IConnection) {
	zlog.Debug("OnConnectionLost:", conn.GetConnection().RemoteAddr())
}

type blockRouter struct {
	znet.BaseRouter
}

var Block = int32(1)

// 模拟阻塞操作
func (r *blockRouter) Handle(request ziface.IRequest) {
	//read client data
	zlog.Infof("recv from client:%s, msgId=%d, data=%s\n", request.GetConnection().RemoteAddr(), request.GetMsgID(), string(request.GetData()))

	// 第一次处理时，模拟任务阻塞操作, Hash 模式下，后面的连接的任务得不到处理
	// DynamicBind 模式下，看后面的连接的任务会得到即使处理，不会因为前面连接的任务阻塞而得不到处理
	// 这里只模拟一次阻塞操作。
	if atomic.CompareAndSwapInt32(&Block, 1, 0) {
		zlog.Infof("blockRouter handle start, msgId=%d, remote:%v\n", request.GetMsgID(), request.GetConnection().RemoteAddr())
		time.Sleep(time.Second * 10)
		//阻塞操作结束
		zlog.Infof("blockRouter handle end, msgId=%d, remote:%v\n", request.GetMsgID(), request.GetConnection().RemoteAddr())
	}

	err := request.GetConnection().SendMsg(2, []byte("pong from server"))
	if err != nil {
		zlog.Error(err)
		return
	}
	zlog.Infof("send pong over, client:%s\n", request.GetConnection().RemoteAddr())
}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	s.AddRouter(1, &blockRouter{})

	s.Serve()
}
