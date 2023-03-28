/**
 * @author uuxia
 * @date 15:57 2023/3/10
 * @description 拦截器管理
 **/

package zinterceptor

import "github.com/aceld/zinx/ziface"

// Builder 责任链构造器
type Builder struct {
	body       []ziface.IInterceptor
	head, tail ziface.IInterceptor
	req        ziface.IcReq
}

func NewBuilder() ziface.IBuilder {
	return &Builder{
		body: make([]ziface.IInterceptor, 0),
	}
}

func (ic *Builder) Head(interceptor ziface.IInterceptor) {
	ic.head = interceptor
}

func (ic *Builder) Tail(interceptor ziface.IInterceptor) {
	ic.tail = interceptor
}

func (ic *Builder) AddInterceptor(interceptor ziface.IInterceptor) {
	ic.body = append(ic.body, interceptor)
}

func (ic *Builder) Execute(req ziface.IcReq) ziface.IcResp {
	ic.req = req

	//将全部拦截器放入Builder中
	var interceptors []ziface.IInterceptor
	if ic.head != nil {
		interceptors = append(interceptors, ic.head)
	}
	if len(ic.body) > 0 {
		interceptors = append(interceptors, ic.body...)
	}
	if ic.tail != nil {
		interceptors = append(interceptors, ic.tail)
	}

	//创建一个拦截器责任链，执行每一个拦截器
	chain := NewChain(interceptors, 0, req)

	//进入责任链执行
	return chain.Proceed(ic.req)
}
