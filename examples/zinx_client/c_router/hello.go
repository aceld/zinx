package c_router

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

//HelloZinxRouter Handle
func (this *HelloRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	zlog.Debug("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

}
