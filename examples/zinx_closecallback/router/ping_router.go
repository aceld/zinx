package router

import (
	"fmt"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

// PingRouter handles ping messages
// PingRouter 处理ping消息
type PingRouter struct {
	znet.BaseRouter
}

// PreHandle processes the request before the main handler
// PreHandle 在主处理器之前处理请求
func (r *PingRouter) PreHandle(req ziface.IRequest) {
}

// Handle processes the main ping message and sends pong response
// Handle 处理主要的ping消息并发送pong响应
func (r *PingRouter) Handle(req ziface.IRequest) {
	// Reply to client / 回复客户端
	if err := req.GetConnection().SendMsg(0, []byte("Pong")); err != nil {
		fmt.Println("SendMsg error:", err)
	}
}

// PostHandle processes the request after the main handler
// PostHandle 在主处理器之后处理请求
func (r *PingRouter) PostHandle(req ziface.IRequest) {
}
