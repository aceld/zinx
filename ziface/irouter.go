// @Title irouter.go
// @Description Provides all interface declarations for message routing
// @Author Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

/*
IRouter is the interface for message routing. The route is the processing
business method set by the framework user for this connection. The IRequest
in the route includes the connection information and the request data
information for this connection.
(路由接口， 这里面路由是 使用框架者给该链接自定的 处理业务方法
路由里的IRequest 则包含用该链接的链接信息和该链接的请求数据信息)
*/
type IRouter interface {
	PreHandle(request IRequest)  //Hook method before processing conn business(在处理conn业务之前的钩子方法)
	Handle(request IRequest)     //Method for processing conn business(处理conn业务的方法)
	PostHandle(request IRequest) //Hook method after processing conn business(处理conn业务之后的钩子方法)
}

/*
RouterHandler is a method slice collection style router. Unlike the old version,
the new version only saves the router method collection, and the specific execution
is handed over to the IRequest of each request.
(方法切片集合式路路由
不同于旧版 新版本仅保存路由方法集合，具体执行交给每个请求的 IRequest)
*/
type RouterHandler func(request IRequest)
type IRouterSlices interface {
	// Add global components (添加全局组件)
	Use(Handlers ...RouterHandler)

	// Add a route (添加业务处理器集合)
	AddHandler(msgId uint32, handlers ...RouterHandler)

	// Router group management （路由分组管理，并且会返回一个组管理器）
	Group(start, end uint32, Handlers ...RouterHandler) IGroupRouterSlices

	// Get the method set collection for processing （获得当前的所有注册在MsgId的处理器集合）
	GetHandlers(MsgId uint32) ([]RouterHandler, bool)
}

type IGroupRouterSlices interface {
	// Add global components (添加全局组件)
	Use(Handlers ...RouterHandler)

	// Add group routing components (添加业务处理器集合)
	AddHandler(MsgId uint32, Handlers ...RouterHandler)
}
