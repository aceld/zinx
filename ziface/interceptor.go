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

type Interceptor interface {
	Intercept(Chain) IcResp
}
type Chain interface {
	Request() IcReq
	Proceed(IcReq) IcResp
}
type InterceptorBuilder interface {
	Head(interceptor Interceptor)
	Tail(interceptor Interceptor)
	AddInterceptor(interceptor Interceptor)
	Execute(request IcReq) IcResp
}
