// @Title irequest.go
// @Description Provides all interface declarations for connection requests
// @Author Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

type HandleStep int

// IFuncRequest function message interface (函数消息接口)
type IFuncRequest interface {
	CallFunc()
}

// IRequest interface:
// It actually packages the connection information and request data of the client request into Request
// (实际上是把客户端请求的链接信息 和 请求的数据 包装到了 Request里)
type IRequest interface {
	GetConnection() IConnection // Get the connection information of the request(获取请求连接信息)

	GetData() []byte  // Get the data of the request message(获取请求消息的数据)
	GetMsgID() uint32 // Get the message ID of the request(获取请求的消息ID)

	GetMessage() IMessage // Get the raw data of the request message (获取请求消息的原始数据 add by uuxia 2023-03-10)

	GetResponse() IcResp // Get the serialized data after parsing(获取解析完后序列化数据)
	SetResponse(IcResp)  // Set the serialized data after parsing(设置解析完后序列化数据)

	BindRouter(router IRouter) // Bind which router handles this request(绑定这次请求由哪个路由处理)
	// Move on to the next handler to start execution, but the function that calls this method will execute in reverse order of their order
	// (转进到下一个处理器开始执行 但是调用此方法的函数会根据先后顺序逆序执行)
	Call()

	//erminate the execution of the processing function, but the function that calls this method will be executed until completion
	// 终止处理函数的运行 但调用此方法的函数会执行完毕
	Abort()

	//Specify which Handler function to execute next in the Handle
	// (指定接下来的Handle去执行哪个Handler函数)
	// Be careful, it will cause loop calling
	// (慎用，会导致循环调用)
	Goto(HandleStep)

	// New router operation
	// (新路由操作)
	BindRouterSlices([]RouterHandler)

	// Execute the next function
	// (执行下一个函数)
	RouterSlicesNext()

	// 重置一个 Request
	//Reset(conn IConnection, msg IMessage)
	// 复制一份 Request 对象
	Copy() IRequest
	//Set 在 Request 中存放一个上下文
	Set(key string, value interface{})
	//Get 从 Request 中获取一个上下文信息
	Get(key string) (value interface{}, exists bool)
}

type BaseRequest struct{}

func (br *BaseRequest) GetConnection() IConnection       { return nil }
func (br *BaseRequest) GetData() []byte                  { return nil }
func (br *BaseRequest) GetMsgID() uint32                 { return 0 }
func (br *BaseRequest) GetMessage() IMessage             { return nil }
func (br *BaseRequest) GetResponse() IcResp              { return nil }
func (br *BaseRequest) SetResponse(resp IcResp)          {}
func (br *BaseRequest) BindRouter(router IRouter)        {}
func (br *BaseRequest) Call()                            {}
func (br *BaseRequest) Abort()                           {}
func (br *BaseRequest) Goto(HandleStep)                  {}
func (br *BaseRequest) BindRouterSlices([]RouterHandler) {}
func (br *BaseRequest) RouterSlicesNext()                {}
func (br *BaseRequest) Copy() IRequest                   { return nil }

func (br *BaseRequest) Set(key string, value interface{}) {}

func (br *BaseRequest) Get(key string) (value interface{}, exists bool) { return nil, false }
