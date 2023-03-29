// Package ziface 主要提供zinx全部抽象层接口定义.
// 包括:
//
//			IServer 服务mod接口
//			IRouter 路由mod接口
//			IConnection 连接mod层接口
//	     IMessage 消息mod接口
//			IDataPack 消息拆解接口
//	     IMsgHandler 消息处理及协程池接口
//
// 当前文件描述:
// @Title  imsghandler.go
// @Description  提供worker启动、处理消息业务调用等接口
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

/*
消息管理抽象层
*/
type IMsgHandle interface {
	//为消息添加具体的处理逻辑, msgID，支持整型，字符串
	AddRouter(msgID uint32, router IRouter)
	StartWorkerPool()                    //启动worker工作池
	SendMsgToTaskQueue(request IRequest) //将消息交给TaskQueue,由worker进行处理

	Execute(request IRequest)                //
	AddInterceptor(interceptor IInterceptor) //注册责任链任务入口，每个拦截器处理完后，数据都会传递至下一个拦截器，使得消息可以层层处理层层传递，顺序取决于注册顺序
}
