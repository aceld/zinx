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
	request    ziface.Request
}

func NewInterceptorBuilder() ziface.InterceptorBuilder {
	return &InterceptorChain{
		body: make([]ziface.Interceptor, 0),
	}
}

func (this *InterceptorChain) Head(interceptor ziface.Interceptor) {
	this.head = interceptor
}

func (this *InterceptorChain) Tail(interceptor ziface.Interceptor) {
	this.tail = interceptor
}

func (this *InterceptorChain) AddInterceptor(interceptor ziface.Interceptor) {
	this.body = append(this.body, interceptor)
}

func (this *InterceptorChain) Execute(request ziface.Request) ziface.Response {
	this.request = request
	var interceptors []ziface.Interceptor
	if this.head != nil {
		interceptors = append(interceptors, this.head)
	}
	if len(this.body) > 0 {
		interceptors = append(interceptors, this.body...)
	}
	if this.tail != nil {
		interceptors = append(interceptors, this.tail)
	}
	chain := NewRealInterceptorChain(interceptors, 0, request)
	return chain.Proceed(this.request)
}
