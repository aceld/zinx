package router

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_decoder/decode"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type HtlvCrcBusinessRouter struct {
	znet.BaseRouter
}

func (this *HtlvCrcBusinessRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HtlvCrcBusinessRouter Handle", request.GetMessage().GetMsgID(), request.GetMessage().GetData())
	msgID := request.GetMessage().GetMsgID()
	if msgID == 0x10 {
		_response := request.GetResponse()
		if _response != nil {
			switch _response.(type) {
			case decode.HtlvCrcData:
				tlvData := _response.(decode.HtlvCrcData)
				fmt.Println("do msgid=0x10 data business", tlvData)
			}
		}
	}
}
