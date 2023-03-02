package c_router

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	zlog.Debug("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	zlog.Debug("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	fmt.Println("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	if err := request.GetConnection().SendBuffMsg(1, []byte("Hello[from client]")); err != nil {
		zlog.Error(err)
	}
}
