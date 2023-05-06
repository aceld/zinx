package znet

import (
	"github.com/aceld/zinx/ziface"
)

const (
	PRE_HANDLE  ziface.HandleStep = iota // PreHandle for pre-processing
	HANDLE                               // Handle for processing
	POST_HANDLE                          // PostHandle for post-processing

	HANDLE_OVER
)

// Request 请求
type Request struct {
	ziface.BaseRequest
	conn     ziface.IConnection     // the connection which has been established with the client(已经和客户端建立好的链接)
	msg      ziface.IMessage        // the request data sent by the client(客户端请求的数据)
	icResp   ziface.IcResp          // response data returned by the interceptors (拦截器返回数据)
	handlers []ziface.RouterHandler // router function slice(路由函数切片)
	index    int8                   // router function slice index(路由函数切片索引)
}

func (r *Request) GetResponse() ziface.IcResp {
	return r.icResp
}

func (r *Request) SetResponse(response ziface.IcResp) {
	r.icResp = response
}

func NewRequest(conn ziface.IConnection, msg ziface.IMessage) ziface.IRequest {
	req := new(Request)
	req.conn = conn
	req.msg = msg
	req.index = -1
	return req
}

func (r *Request) GetMessage() ziface.IMessage {
	return r.msg
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}

func (r *Request) BindRouter(router ziface.IRouter) {
	r.BindRouterSlices([]ziface.RouterHandler{
		router.PreHandle,
		router.Handle,
		router.PostHandle,
	})
}
func (r *Request) Goto(step ziface.HandleStep) {
	r.index = int8(step) - 1
}

func (r *Request) Call() {
	r.RouterSlicesNext()
	r.index = -1
}

func (r *Request) Abort() {
	r.index = int8(len(r.handlers))
}

// New version
func (r *Request) BindRouterSlices(handlers []ziface.RouterHandler) {
	r.handlers = handlers
}

func (r *Request) RouterSlicesNext() {
	r.index++
	for r.index < int8(len(r.handlers)) {
		r.handlers[r.index](r)
		r.index++
	}
}
