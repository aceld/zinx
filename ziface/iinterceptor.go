/**
 * @author uuxia
 * @date 15:54 2023/3/10
 * @description //TODO
 **/

package ziface

// 请求父类，定义空接口，用于扩展支持任意类型
type IcReq interface{}

// 回复父类，定义空接口，用于扩展支持任意类型
type IcResp interface{}

// 拦截器
type IInterceptor interface {
	Intercept(IChain) IcResp
}

// 责任链
type IChain interface {
	Request() IcReq
	Proceed(IcReq) IcResp
}

type IBuilder interface {
	Head(interceptor IInterceptor)
	Tail(interceptor IInterceptor)
	AddInterceptor(interceptor IInterceptor)
	Execute(request IcReq) IcResp
}
