package router

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_decoder/decode"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type TLVBusinessRouter struct {
	znet.BaseRouter
}

func (this *TLVBusinessRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call TLVRouter Handle", request.GetMessage().GetMsgID(), request.GetMessage().GetData())
	msgID := request.GetMessage().GetMsgID()
	if msgID == 0x00000001 {
		_response := request.GetResponse()
		if _response != nil {
			switch _response.(type) {
			case decode.TlvData:
				tlvData := _response.(decode.TlvData)
				fmt.Println("do msgid=0x00000001 data business", tlvData)
			}
		}
	}

}
