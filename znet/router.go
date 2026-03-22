package znet

import (
	"strconv"
	"sync"

	"github.com/aceld/zinx/ziface"
)

// BaseRouter is used as the base class when implementing a router.
// Depending on the needs, the methods of this base class can be overridden.
// (实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写)
type BaseRouter struct{}

// Here, all of BaseRouter's methods are empty, because some routers may not want to have PreHandle or PostHandle.
// Therefore, inheriting all routers from BaseRouter has the advantage that PreHandle and PostHandle do not need to be
// implemented to instantiate a router.
// (这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化)

// PreHandle -
func (br *BaseRouter) PreHandle(req ziface.IRequest) {}

// Handle -
func (br *BaseRouter) Handle(req ziface.IRequest) {}

// PostHandle -
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}

// New slice-based router
// The new version of the router has basic logic that allows users to pass in varying numbers of router handlers.
// The router will save all of these router handler functions and find them when a request comes in, then execute them using IRequest.
// The router can set globally shared components using the Use method.
// The router can be grouped using Group, and groups also have their own Use method for setting group-shared components.
// (新切片集合式路由
// 新版本路由基本逻辑,用户可以传入不等数量的路由路由处理器
// 路由本体会讲这些路由处理器函数全部保存,在请求来的时候找到，并交由IRequest去执行
// 路由可以设置全局的共用组件通过Use方法
// 路由可以分组,通过Group,分组也有自己对应Use方法设置组共有组件)

type RouterSlices struct {
	Apis     map[uint32][]ziface.RouterHandler
	Handlers []ziface.RouterHandler
	sync.RWMutex
}

func NewRouterSlices() *RouterSlices {
	return &RouterSlices{
		Apis:     make(map[uint32][]ziface.RouterHandler, 10),
		Handlers: make([]ziface.RouterHandler, 0, 6),
	}
}

func (r *RouterSlices) Use(handles ...ziface.RouterHandler) {
	r.Handlers = append(r.Handlers, handles...)
}

func (r *RouterSlices) AddHandler(msgId uint32, Handlers ...ziface.RouterHandler) {
	// 1. Check if the API handler method bound to the current msg already exists
	if _, ok := r.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}

	finalSize := len(r.Handlers) + len(Handlers)
	mergedHandlers := make([]ziface.RouterHandler, finalSize)
	copy(mergedHandlers, r.Handlers)
	copy(mergedHandlers[len(r.Handlers):], Handlers)
	r.Apis[msgId] = append(r.Apis[msgId], mergedHandlers...)
}

func (r *RouterSlices) GetHandlers(MsgId uint32) ([]ziface.RouterHandler, bool) {
	r.RLock()
	defer r.RUnlock()
	handlers, ok := r.Apis[MsgId]
	return handlers, ok
}

func (r *RouterSlices) Group(start, end uint32, Handlers ...ziface.RouterHandler) ziface.IGroupRouterSlices {
	return NewGroup(start, end, r, Handlers...)
}

type GroupRouter struct {
	start    uint32
	end      uint32
	Handlers []ziface.RouterHandler
	router   ziface.IRouterSlices
}

func NewGroup(start, end uint32, router *RouterSlices, Handlers ...ziface.RouterHandler) *GroupRouter {
	g := &GroupRouter{
		start:    start,
		end:      end,
		Handlers: make([]ziface.RouterHandler, 0, len(Handlers)),
		router:   router,
	}
	g.Handlers = append(g.Handlers, Handlers...)
	return g
}

func (g *GroupRouter) Use(Handlers ...ziface.RouterHandler) {
	g.Handlers = append(g.Handlers, Handlers...)
}

func (g *GroupRouter) AddHandler(MsgId uint32, Handlers ...ziface.RouterHandler) {
	if MsgId < g.start || MsgId > g.end {
		panic("add router to group err in msgId:" + strconv.Itoa(int(MsgId)))
	}

	finalSize := len(g.Handlers) + len(Handlers)
	mergedHandlers := make([]ziface.RouterHandler, finalSize)
	copy(mergedHandlers, g.Handlers)
	copy(mergedHandlers[len(g.Handlers):], Handlers)

	g.router.AddHandler(MsgId, mergedHandlers...)
}

// RouterSlicesContext is the context-based router slices implementation
// (基于Context的路由切片实现)
type RouterSlicesContext struct {
	Apis     map[uint32][]ziface.HandlerFunc
	Handlers []ziface.HandlerFunc
	sync.RWMutex
}

// NewRouterSlicesContext creates a new RouterSlicesContext
// (创建一个新的RouterSlicesContext)
func NewRouterSlicesContext() *RouterSlicesContext {
	return &RouterSlicesContext{
		Apis:     make(map[uint32][]ziface.HandlerFunc, 10),
		Handlers: make([]ziface.HandlerFunc, 0, 6),
	}
}

// Use adds global middleware handlers
// (添加全局中间件处理程序)
func (r *RouterSlicesContext) Use(handles ...ziface.HandlerFunc) {
	r.Handlers = append(r.Handlers, handles...)
}

// AddHandler adds route handlers for a specific message ID
// (为特定消息ID添加路由处理程序)
func (r *RouterSlicesContext) AddHandler(msgId uint32, Handlers ...ziface.HandlerFunc) {
	// 1. Check if the API handler method bound to the current msg already exists
	if _, ok := r.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}

	finalSize := len(r.Handlers) + len(Handlers)
	mergedHandlers := make([]ziface.HandlerFunc, finalSize)
	copy(mergedHandlers, r.Handlers)
	copy(mergedHandlers[len(r.Handlers):], Handlers)
	r.Apis[msgId] = append(r.Apis[msgId], mergedHandlers...)
}

// GetHandlers returns the handlers for a specific message ID
// (返回特定消息ID的处理程序)
func (r *RouterSlicesContext) GetHandlers(MsgId uint32) ([]ziface.HandlerFunc, bool) {
	r.RLock()
	defer r.RUnlock()
	handlers, ok := r.Apis[MsgId]
	return handlers, ok
}

// Group creates a route group
// (创建路由分组)
func (r *RouterSlicesContext) Group(start, end uint32, Handlers ...ziface.HandlerFunc) ziface.IGroupRouterSlicesContext {
	return NewGroupContext(start, end, r, Handlers...)
}

// GroupRouterContext is the context-based group router implementation
// (基于Context的分组路由实现)
type GroupRouterContext struct {
	start    uint32
	end      uint32
	Handlers []ziface.HandlerFunc
	router   ziface.IRouterSlicesContext
}

// NewGroupContext creates a new GroupRouterContext
// (创建一个新的GroupRouterContext)
func NewGroupContext(start, end uint32, router *RouterSlicesContext, Handlers ...ziface.HandlerFunc) *GroupRouterContext {
	g := &GroupRouterContext{
		start:    start,
		end:      end,
		Handlers: make([]ziface.HandlerFunc, 0, len(Handlers)),
		router:   router,
	}
	g.Handlers = append(g.Handlers, Handlers...)
	return g
}

// Use adds global middleware handlers to the group
// (向组添加全局中间件处理程序)
func (g *GroupRouterContext) Use(Handlers ...ziface.HandlerFunc) {
	g.Handlers = append(g.Handlers, Handlers...)
}

// AddHandler adds route handlers for a specific message ID in the group
// (在组中为特定消息ID添加路由处理程序)
func (g *GroupRouterContext) AddHandler(MsgId uint32, Handlers ...ziface.HandlerFunc) {
	if MsgId < g.start || MsgId > g.end {
		panic("add router to group err in msgId:" + strconv.Itoa(int(MsgId)))
	}

	finalSize := len(g.Handlers) + len(Handlers)
	mergedHandlers := make([]ziface.HandlerFunc, finalSize)
	copy(mergedHandlers, g.Handlers)
	copy(mergedHandlers[len(g.Handlers):], Handlers)

	g.router.AddHandler(MsgId, mergedHandlers...)
}
