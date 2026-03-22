package router

import (
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/zlog"
	"github.com/aceld/zinx/v3/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

func (hr *HelloRouter) Handle(request ziface.IRequest) {
	zlog.Ins().InfoF(string(request.GetData()))
}
