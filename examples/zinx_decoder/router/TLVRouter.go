package router

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_decoder/server/interceptor"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type TLVRouter struct {
	znet.BaseRouter
}

func (this *TLVRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call TLVRouter Handle", request.GetMessage().GetData())
	_response := request.GetResponse()
	if _response != nil {
		switch _response.(type) {
		case interceptor.TLVData:
			tlvData := _response.(interceptor.TLVData)
			fmt.Println(" Response TLVData", tlvData)
		}
	}
}
