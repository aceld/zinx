/**
 * @author uuxia
 * @date 15:54 2023/3/10
 * @description
 **/

package ziface

// Input data for interceptor
// (拦截器输入数据)
type IcReq interface{}

// Output data for interceptor
// (拦截器输出数据)
type IcResp interface{}

// Interceptor
// (拦截器)
type IInterceptor interface {
	Intercept(IChain) IcResp
	// The interception method of the interceptor (defined by the developer)
	// (拦截器的拦截处理方法,由开发者定义)
}

// Responsibility chain
// (责任链)
type IChain interface {
	Request() IcReq        // Get the request data in the current chain (current interceptor)-获取当前责任链中的请求数据(当前拦截器)
	GetIMessage() IMessage // Get IMessage from Chain (从Chain中获取IMessage)
	Proceed(IcReq) IcResp  // Enter and execute the next interceptor, and pass the request data to the next interceptor (进入并执行下一个拦截器，且将请求数据传递给下一个拦截器)
	ProceedWithIMessage(IMessage, IcReq) IcResp
}
