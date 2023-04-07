/**
 * @author uuxia
 * @date 15:56 2023/3/10
 * @description 责任链模式
 **/

package zinterceptor

import "github.com/aceld/zinx/ziface"

type Chain struct {
	req          ziface.IcReq
	position     int
	interceptors []ziface.IInterceptor
}

func (c *Chain) Request() ziface.IcReq {
	return c.req
}

func (c *Chain) Proceed(request ziface.IcReq) ziface.IcResp {
	if c.position < len(c.interceptors) {
		//chain := NewChain(c.interceptors, c.position+1, request)
		// 不需要每次拷贝实例化，
		interceptor := c.interceptors[c.position]
		c.position += 1
		response := interceptor.Intercept(c)
		return response
	}
	return request
}

func NewChain(list []ziface.IInterceptor, pos int, req ziface.IcReq) ziface.IChain {
	return &Chain{
		req:          req,
		position:     pos,
		interceptors: list,
	}
}
