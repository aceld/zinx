package router

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_decoder/server/interceptor"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

// HTLVCRCRouter test 自定义路由
type HTLVCRCRouter struct {
	znet.BaseRouter
}

// HTLVCRCRouter Handle
func (this *HTLVCRCRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HTLVCRCRouter Handle", request.GetMessage().GetData())
	_response := request.GetResponse()
	if _response != nil {
		switch _response.(type) {
		case interceptor.Data:
			htlvData := _response.(interceptor.Data)
			fmt.Println(" Response HTLVCRCData", htlvData)
		}
	}
}
