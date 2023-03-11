package router

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type TLVRouter struct {
	znet.BaseRouter
}

func (this *TLVRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call TLVRouter Handle", request.GetMessage().GetData())
}
