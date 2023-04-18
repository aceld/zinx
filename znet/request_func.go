package znet

import (
	"github.com/aceld/zinx/ziface"
)

func NewFuncRequest(conn ziface.IConnection, callFunc func()) ziface.IFuncRequest {
	req := new(RequestFunc)
	req.conn = conn
	req.callFunc = callFunc
	return req
}

type RequestFunc struct {
	conn     ziface.IConnection //已经和客户端建立好的 链接
	callFunc func()
}

// GetConnection 获取请求连接信息
func (rf *RequestFunc) GetConnection() ziface.IConnection {
	return rf.conn
}

func (rf *RequestFunc) CallFunc() {
	if rf.callFunc != nil {
		rf.callFunc()
	}
}
