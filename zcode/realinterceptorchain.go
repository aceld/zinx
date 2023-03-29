/**
 * @author uuxia
 * @date 15:56 2023/3/10
 * @description 责任链模式
 **/

package zcode

import "github.com/aceld/zinx/ziface"

type RealInterceptorChain struct {
	req          ziface.IcReq
	position     int
	interceptors []ziface.Interceptor
}

func (ric *RealInterceptorChain) Request() ziface.IcReq {
	return ric.req
}

func (ric *RealInterceptorChain) Proceed(request ziface.IcReq) ziface.IcResp {
	if ric.position < len(ric.interceptors) {
		chain := NewRealInterceptorChain(ric.interceptors, ric.position+1, request)
		interceptor := ric.interceptors[ric.position]
		response := interceptor.Intercept(chain)
		return response
	}
	return request
}

func NewRealInterceptorChain(list []ziface.Interceptor, pos int, req ziface.IcReq) ziface.Chain {
	return &RealInterceptorChain{
		req:          req,
		position:     pos,
		interceptors: list,
	}
}
