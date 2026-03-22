package c_router

import (
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/zlog"
	"github.com/aceld/zinx/v3/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call PingRouter Handle")

	zlog.Debug("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	if err := request.GetConnection().SendBuffMsg(1, []byte("Hello[from client]")); err != nil {
		zlog.Error(err)
	}
}
