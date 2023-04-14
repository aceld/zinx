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

// chainBuilder 责任链构造器
type chainBuilder struct {
	body       []ziface.IInterceptor
	head, tail ziface.IInterceptor
}

func newChainBuilder() *chainBuilder {
	return &chainBuilder{
		body: make([]ziface.IInterceptor, 0),
	}
}

// Head 将拦截器添加到责任链头部
func (ic *chainBuilder) Head(interceptor ziface.IInterceptor) {
	ic.head = interceptor
}

// Tail 将拦截器添加到责任链尾部
func (ic *chainBuilder) Tail(interceptor ziface.IInterceptor) {
	ic.tail = interceptor
}

// AddInterceptor 顺位添加一个拦截器到责任链中
func (ic *chainBuilder) AddInterceptor(interceptor ziface.IInterceptor) {
	ic.body = append(ic.body, interceptor)
}

// Execute 依次执行当前责任链上所有拦截器
func (ic *chainBuilder) Execute(req ziface.IcReq) ziface.IcResp {

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
	chain := zinterceptor.NewChain(interceptors, 0, req)

	//进入责任链执行
	return chain.Proceed(req)
}
