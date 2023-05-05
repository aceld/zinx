/**
 * @author uuxia
 * @date 15:57 2023/3/10
 * @description 拦截器管理
 **/

package znet

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zinterceptor"
)

// chainBuilder is a builder for creating a chain of interceptors.
// (责任链构造器)
type chainBuilder struct {
	body       []ziface.IInterceptor
	head, tail ziface.IInterceptor
}

// newChainBuilder creates a new instance of chainBuilder.
func newChainBuilder() *chainBuilder {
	return &chainBuilder{
		body: make([]ziface.IInterceptor, 0),
	}
}

// Head adds an interceptor to the head of the chain.
func (ic *chainBuilder) Head(interceptor ziface.IInterceptor) {
	ic.head = interceptor
}

// Tail adds an interceptor to the tail of the chain.
func (ic *chainBuilder) Tail(interceptor ziface.IInterceptor) {
	ic.tail = interceptor
}

// AddInterceptor adds an interceptor to the body of the chain.
func (ic *chainBuilder) AddInterceptor(interceptor ziface.IInterceptor) {
	ic.body = append(ic.body, interceptor)
}

// Execute executes all the interceptors in the current chain in order.
func (ic *chainBuilder) Execute(req ziface.IcReq) ziface.IcResp {

	// Put all the interceptors into the builder
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

	// Create a new interceptor chain and execute each interceptor
	chain := zinterceptor.NewChain(interceptors, 0, req)

	// Execute the chain
	return chain.Proceed(req)
}
