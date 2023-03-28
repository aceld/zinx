/**
 * @author uuxia
 * @date 15:54 2023/3/10
 * @description //TODO
 **/

package ziface

// 请求父类，定义空接口，用于扩展支持任意类型

type Request interface {
}

// 回复父类，定义空接口，用于扩展支持任意类型

type Response interface {
}
type Interceptor interface {
	Intercept(Chain) Response
}
type Chain interface {
	Request() Request
	Proceed(Request) Response
}
type InterceptorBuilder interface {
	Head(interceptor Interceptor)
	Tail(interceptor Interceptor)
	AddInterceptor(interceptor Interceptor)
	Execute(request Request) Response
}
