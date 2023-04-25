package router

import (
	"github.com/aceld/zinx/zdecoder"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

type TLVBusinessRouter struct {
	znet.BaseRouter
}

func (this *TLVBusinessRouter) Handle(request ziface.IRequest) {

	msgID := request.GetMessage().GetMsgID()
	zlog.Ins().DebugF("Call TLVRouter Handle %d %+v\n", msgID, request.GetMessage().GetData())

	resp := request.GetResponse()
	if resp == nil {
		return
	}

	tlvData := resp.(zdecoder.TLVDecoder)
	zlog.Ins().DebugF("do msgid=0x00000001 data business %+v\n", tlvData)
}
