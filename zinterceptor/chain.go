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

func NewChain(list []ziface.IInterceptor, pos int, req ziface.IcReq) ziface.IChain {
	return &Chain{
		req:          req,
		position:     pos,
		interceptors: list,
	}
}

func (c *Chain) Request() ziface.IcReq {
	return c.req
}

func (c *Chain) Proceed(request ziface.IcReq) ziface.IcResp {
	if c.position < len(c.interceptors) {
		chain := NewChain(c.interceptors, c.position+1, request)
		interceptor := c.interceptors[c.position]
		response := interceptor.Intercept(chain)
		return response
	}
	return request
}

// GetIMessage  从Chain中获取IMessage
func (c *Chain) GetIMessage() ziface.IMessage {

	req := c.Request()
	if req == nil {
		return nil
	}

	iRequest := c.ShouldIRequest(req)
	if iRequest == nil {
		return nil
	}

	return iRequest.GetMessage()
}

// Next 通过IMessage和解码后数据进入下一个责任链任务
// iMessage 为解码后的IMessage
// response 为解码后的数据
func (c *Chain) ProceedWithIMessage(iMessage ziface.IMessage, response ziface.IcReq) ziface.IcResp {
	if iMessage == nil || response == nil {
		return c.Proceed(c.Request())
	}

	req := c.Request()
	if req == nil {
		return c.Proceed(c.Request())
	}

	iRequest := c.ShouldIRequest(req)
	if iRequest == nil {
		return c.Proceed(c.Request())
	}

	//设置chain的request下一次请求
	iRequest.SetResponse(response)

	return c.Proceed(iRequest)
}

// ShouldIRequest 判断是否是IRequest
func (c *Chain) ShouldIRequest(icReq ziface.IcReq) ziface.IRequest {
	if icReq == nil {
		return nil
	}

	switch icReq.(type) {
	case ziface.IRequest:
		return icReq.(ziface.IRequest)
	default:
		return nil
	}
}
