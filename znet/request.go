package znet

import (
	"sync"

	"github.com/aceld/zinx/zconf"
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
	router   ziface.IRouter         // the router that handles this request(请求处理的函数)
	steps    ziface.HandleStep      // used to control the execution of router functions(用来控制路由函数执行)
	stepLock *sync.RWMutex          // concurrency lock(并发互斥)
	needNext bool                   // whether to execute the next router function(是否需要执行下一个路由函数)
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
	req.steps = PRE_HANDLE
	req.conn = conn
	req.msg = msg
	req.stepLock = new(sync.RWMutex)
	req.needNext = true
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
	r.router = router
}

func (r *Request) next() {
	if r.needNext == false {
		r.needNext = true
		return
	}

	r.stepLock.Lock()
	r.steps++
	r.stepLock.Unlock()
}

func (r *Request) Goto(step ziface.HandleStep) {
	r.stepLock.Lock()
	r.steps = step
	r.needNext = false
	r.stepLock.Unlock()
}

func (r *Request) Call() {

	if r.router == nil {
		return
	}

	for r.steps < HANDLE_OVER {
		switch r.steps {
		case PRE_HANDLE:
			r.router.PreHandle(r)
		case HANDLE:
			r.router.Handle(r)
		case POST_HANDLE:
			r.router.PostHandle(r)
		}

		r.next()
	}

	r.steps = PRE_HANDLE
}

func (r *Request) Abort() {
	if zconf.GlobalObject.RouterSlicesMode {
		r.index = int8(len(r.handlers))
	} else {
		r.stepLock.Lock()
		r.steps = HANDLE_OVER
		r.stepLock.Unlock()
	}
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
