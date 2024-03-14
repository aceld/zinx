package znet

import (
	"math"
	"sync"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zpack"
)

const (
	PRE_HANDLE  ziface.HandleStep = iota // PreHandle for pre-processing
	HANDLE                               // Handle for processing
	POST_HANDLE                          // PostHandle for post-processing

	HANDLE_OVER
)

var RequestPool = new(sync.Pool)

func init() {
	RequestPool.New = func() interface{} {
		return allocateRequest()
	}
}

// Request 请求
type Request struct {
	ziface.BaseRequest
	conn     ziface.IConnection     // the connection which has been established with the client(已经和客户端建立好的链接)
	msg      ziface.IMessage        // the request data sent by the client(客户端请求的数据)
	router   ziface.IRouter         // the router that handles this request(请求处理的函数)
	steps    ziface.HandleStep      // used to control the execution of router functions(用来控制路由函数执行)
	stepLock sync.RWMutex           // concurrency lock(并发互斥)
	needNext bool                   // whether to execute the next router function(是否需要执行下一个路由函数)
	icResp   ziface.IcResp          // response data returned by the interceptors (拦截器返回数据)
	handlers []ziface.RouterHandler // router function slice(路由函数切片)
	index    int8                   // router function slice index(路由函数切片索引)
	keys     map[string]interface{} // keys 路由处理时可能会存取的上下文信息
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
	req.stepLock = sync.RWMutex{}
	req.needNext = true
	req.index = -1
	return req
}

func GetRequest(conn ziface.IConnection, msg ziface.IMessage) ziface.IRequest {
	// 从对象池中取得一个 Request 对象,如果池子中没有可用的 Request 对象则会调用 allocateRequest 函数构造一个新的对象分配
	r := RequestPool.Get().(*Request)
	// 因为取出的 Request 对象可能是已存在也可能是新构造的,无论是哪种情况都应该初始化再返回使用
	r.Reset(conn, msg)
	return r
}

func PutRequest(request ziface.IRequest) {
	RequestPool.Put(request)
}

func allocateRequest() ziface.IRequest {
	req := new(Request)
	req.steps = PRE_HANDLE
	req.needNext = true
	req.index = -1
	return req
}

func (r *Request) Reset(conn ziface.IConnection, msg ziface.IMessage) {
	r.steps = PRE_HANDLE
	r.conn = conn
	r.msg = msg
	r.needNext = true
	r.index = -1
	r.keys = nil

}

// Copy 在执行路由函数的时候可能会出现需要再起一个协程的需求,但是 Request 对象由对象池管理后无法保证新协程中的 Request 参数一致
// 通过 Copy 方法复制一份 Request 对象保持创建协程时候的参数一致。但新开的协程不应该在对原始的执行过程有影响，所以不包含链接和路由对象。
func (r *Request) Copy() ziface.IRequest {
	// 构造一个新的 Request 对象，复制部分原始对象的参数,但是复制的 Request 不应该再对原始链接操作,所以不能含有链接参数
	// 同理也不应该再执行路由方法,路由函数也不包含
	newRequest := &Request{
		conn:     nil,
		router:   nil,
		steps:    r.steps,
		needNext: false,
		icResp:   nil,
		handlers: nil,
		index:    math.MaxInt8,
	}

	// 复制原本的上下文信息
	newRequest.keys = make(map[string]interface{})
	for k, v := range r.keys {
		newRequest.keys[k] = v
	}

	// 复制一份原本的 msg 信息
	newRequest.msg = zpack.NewMessageByMsgId(r.msg.GetMsgID(), r.msg.GetDataLen(), r.msg.GetRawData())

	return newRequest
}

// Set 在 Request 中存放一个上下文，如果 keys 为空会实例化一个
func (r *Request) Set(key string, value interface{}) {
	r.stepLock.Lock()
	if r.keys == nil {
		r.keys = make(map[string]interface{})
	}

	r.keys[key] = value
	r.stepLock.Unlock()
}

// Get 在 Request 中取出一个上下文信息
func (r *Request) Get(key string) (value interface{}, exists bool) {
	r.stepLock.RLock()
	value, exists = r.keys[key]
	r.stepLock.RUnlock()
	return
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

// BindRouterSlices New version
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
