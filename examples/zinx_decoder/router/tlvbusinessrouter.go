package router

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"github.com/aceld/zinx/zpack"
)

type TLVBusinessRouter struct {
	znet.BaseRouter
}

func (this *TLVBusinessRouter) Handle(request ziface.IRequest) {
	zlog.Ins().DebugF("Call TLVRouter Handle %d %+v\n", request.GetMessage().GetMsgID(), request.GetMessage().GetData())
	msgID := request.GetMessage().GetMsgID()
	if msgID == 0x00000001 {
		_response := request.GetResponse()
		if _response != nil {
			switch _response.(type) {
			case zpack.TLVDecoder:
				tlvData := _response.(zpack.TLVDecoder)
				zlog.Ins().DebugF("do msgid=0x00000001 data business %+v\n", tlvData)
			}
		}
	}

}
