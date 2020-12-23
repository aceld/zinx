package zrouter

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

type HelloZinxRouter struct {
	znet.BaseRouter
}

//HelloZinxRouter Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	zlog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(1, []byte("Hello Zinx Router V0.10"))
	if err != nil {
		zlog.Error(err)
	}
}

