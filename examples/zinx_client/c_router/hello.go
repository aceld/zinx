package c_router

import (
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/zlog"
	"github.com/aceld/zinx/v3/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

// HelloZinxRouter Handle
func (this *HelloRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call HelloZinxRouter Handle")

	zlog.Debug("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
}
