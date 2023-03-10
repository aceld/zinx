/**
 * @author uuxia
 * @date 15:56 2023/3/10
 * @description 责任链模式
 **/

package zcode

import "github.com/aceld/zinx/ziface"

type RealInterceptorChain struct {
	request      ziface.Request
	position     int
	interceptors []ziface.Interceptor
}

func (this *RealInterceptorChain) Request() ziface.Request {
	return this.request
}

func (this *RealInterceptorChain) Proceed(request ziface.Request) ziface.Response {
	if this.position < len(this.interceptors) {
		chain := NewRealInterceptorChain(this.interceptors, this.position+1, request)
		interceptor := this.interceptors[this.position]
		response := interceptor.Intercept(chain)
		return response
	}
	return request
}

func NewRealInterceptorChain(list []ziface.Interceptor, pos int, request ziface.Request) ziface.Chain {
	return &RealInterceptorChain{
		request:      request,
		position:     pos,
		interceptors: list,
	}
}
