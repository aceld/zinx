package znet

import (
	"github.com/aceld/zinx/ziface"
	"sync"
)

const (
	PRE_HANDLE  ziface.HandleStep = iota //PreHandle 预处理
	HANDLE                               //Handle 处理
	POST_HANDLE                          //PostHandle 后处理

	HANDLE_OVER
)

type BaseRequest struct {
}

func (br *BaseRequest) GetConnection() ziface.IConnection       { return nil }
func (br *BaseRequest) GetData() []byte                         { return nil }
func (br *BaseRequest) GetMsgID() uint32                        { return 0 }
func (br *BaseRequest) GetMessage() ziface.IMessage             { return nil }
func (br *BaseRequest) GetResponse() ziface.IcResp              { return nil }
func (br *BaseRequest) SetResponse(resp ziface.IcResp)          {}
func (br *BaseRequest) BindRouter(router ziface.IRouter)        {}
func (br *BaseRequest) Call()                                   {}
func (br *BaseRequest) Abort()                                  {}
func (br *BaseRequest) Goto(ziface.HandleStep)                  {}
func (br *BaseRequest) BindRouterSlices([]ziface.RouterHandler) {}
func (br *BaseRequest) RouterSlicesNext()                       {}
func (br *BaseRequest) RouterAbort()                            {}

// RequestFunc 函数调用
type RequestFunc struct {
	BaseRequest
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

func NewFuncRequest(conn ziface.IConnection, callFunc func()) ziface.IRequest {
	req := new(RequestFunc)
	req.conn = conn
	req.callFunc = callFunc
	return req
}

// Request 请求
type Request struct {
	BaseRequest
	conn     ziface.IConnection     //已经和客户端建立好的 链接
	msg      ziface.IMessage        //客户端请求的数据
	router   ziface.IRouter         //请求处理的函数
	steps    ziface.HandleStep      //用来控制路由函数执行
	stepLock *sync.RWMutex          //并发互斥
	needNext bool                   //是否需要执行下一个路由函数
	icResp   ziface.IcResp          //拦截器返回数据
	handlers []ziface.RouterHandler //路由函数切片
	index    int8                   //路由函数切片索引
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

	return req
}

// GetMessage 获取消息实体
func (r *Request) GetMessage() ziface.IMessage {
	return r.msg
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取请求的消息的ID
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
	r.stepLock.Lock()
	r.steps = HANDLE_OVER
	r.stepLock.Unlock()
}

// 新版本路由操作

func (r *Request) BindRouterSlices(handlers []ziface.RouterHandler) {
	r.handlers = handlers
	r.index = -1
}

func (r *Request) RouterSlicesNext() {
	r.index++
	for r.index < int8(len(r.handlers)) {
		r.handlers[r.index](r)
		r.index++
	}
}

func (r *Request) RouterAbort() {
	r.index = int8(len(r.handlers))
}
