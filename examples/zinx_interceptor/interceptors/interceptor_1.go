package interceptors

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
)

// 自定义拦截器1

type MyInterceptor struct{}

func (m *MyInterceptor) Intercept(chain ziface.IChain) ziface.IcResp {
	request := chain.Request()
	// 这一层是自定义拦截器处理逻辑，这里只是简单打印输入
	iRequest := request.(ziface.IRequest)
	zlog.Ins().InfoF("自定义拦截器, 收到消息：%s", iRequest.GetData())
	return chain.Proceed(chain.Request())
}
