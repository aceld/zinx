/**
 * @author uuxia
 * @date 15:57 2023/3/10
 * @description 拦截器管理
 **/

package zcode

import "github.com/aceld/zinx/ziface"

// InterceptorChain
// HTLV+CRC，H头码，T功能码，L数据长度，V数据内容
// +------+-------+---------+--------+--------+
// | 头码  | 功能码 | 数据长度 | 数据内容 | CRC校验 |
// | 1字节 | 1字节  | 1字节   | N字节   |  2字节  |
// +------+-------+---------+--------+--------+
type InterceptorChain struct {
	body       []ziface.Interceptor
	head, tail ziface.Interceptor
	req        ziface.IcReq
}

func NewInterceptorBuilder() ziface.InterceptorBuilder {
	return &InterceptorChain{
		body: make([]ziface.Interceptor, 0),
	}
}

func (ic *InterceptorChain) Head(interceptor ziface.Interceptor) {
	ic.head = interceptor
}

func (ic *InterceptorChain) Tail(interceptor ziface.Interceptor) {
	ic.tail = interceptor
}

func (ic *InterceptorChain) AddInterceptor(interceptor ziface.Interceptor) {
	ic.body = append(ic.body, interceptor)
}

func (ic *InterceptorChain) Execute(req ziface.IcReq) ziface.IcResp {
	ic.req = req
	var interceptors []ziface.Interceptor
	if ic.head != nil {
		interceptors = append(interceptors, ic.head)
	}
	if len(ic.body) > 0 {
		interceptors = append(interceptors, ic.body...)
	}
	if ic.tail != nil {
		interceptors = append(interceptors, ic.tail)
	}
	chain := NewRealInterceptorChain(interceptors, 0, req)
	return chain.Proceed(ic.req)
}
