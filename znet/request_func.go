package znet

import "github.com/aceld/zinx/ziface"

type RequestFunc struct {
	ziface.BaseRequest
	conn     ziface.IConnection
	callFunc func()
}

func (rf *RequestFunc) GetConnection() ziface.IConnection {
	return rf.conn
}

func (rf *RequestFunc) CallFunc() {
	if rf.callFunc != nil {
		rf.callFunc()
	}
}

func NewFuncRequest(conn ziface.IConnection, callFunc func()) ziface.IRequest {
	req := new(RequestFunc)
	req.conn = conn
	req.callFunc = callFunc
	return req
}
