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
	interceptors []ziface.Interceptor
	request      ziface.Request
}

func NewInterceptorBuilder() ziface.InterceptorBuilder {
	return &InterceptorChain{
		interceptors: make([]ziface.Interceptor, 0),
	}
}

func (this *InterceptorChain) AddInterceptor(interceptor ziface.Interceptor) {
	this.interceptors = append(this.interceptors, interceptor)
}

func (this *InterceptorChain) Execute(request ziface.Request) ziface.Response {
	this.request = request
	chain := NewRealInterceptorChain(this.interceptors, 0, request)
	return chain.Proceed(this.request)
}
