package znet

import "github.com/aceld/zinx/ziface"

//BaseRouter 实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
//type BaseRouter struct{}

//这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化

////PreHandle -
//func (br *BaseRouter) PreHandle(req ziface.IRequest) {}
//
////Handle -
//func (br *BaseRouter) Handle(req ziface.IRequest) {}
//
////PostHandle -
//func (br *BaseRouter) PostHandle(req ziface.IRequest) {}

type Router struct {
	index    int8 //函数索引
	handlers []ziface.RouterHandler
}

func (r *Router) Next(request ziface.IRequest) {
	r.index++
	for r.index < int8(len(r.handlers)) {
		r.handlers[r.index](r, request)
		r.index++
	}
}

func (r *Router) Abort() {
	r.index = int8(len(r.handlers))
}

func (r *Router) IsAbort() bool {
	return r.index >= int8(len(r.handlers))
}

func (r *Router) Reset() {
	r.index = -1
	r.handlers = make([]ziface.RouterHandler, 0, 1)
}

func (r *Router) Reindx() {
	r.index = -1
}
func (r *Router) AddHandler(handler ...ziface.RouterHandler) {
	r.handlers = append(r.handlers, handler...)
}
