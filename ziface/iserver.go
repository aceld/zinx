// Package ziface 主要提供zinx全部抽象层接口定义.
// 包括:
//		IServer 服务mod接口
//		IRouter 路由mod接口
//		IConnection 连接mod层接口
//      IMessage 消息mod接口
//		IDataPack 消息拆解接口
//      IMsgHandler 消息处理及协程池接口
//
// 当前文件描述:
// @Title  iserver.go
// @Description  提供Server抽象层全部接口声明
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

//定义服务接口
type IServer interface {
	Start()		//启动服务器方法
	Stop() 		//停止服务器方法
	Serve() 	//开启业务服务方法
	AddRouter(msgID uint32, router IRouter) //路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	GetConnMgr() IConnManager 				//得到链接管理
	SetOnConnStart(func(IConnection)) 		//设置该Server的连接创建时Hook函数
	SetOnConnStop(func(IConnection)) 		//设置该Server的连接断开时的Hook函数
	CallOnConnStart(conn IConnection) 		//调用连接OnConnStart Hook函数
	CallOnConnStop(conn IConnection) 		//调用连接OnConnStop Hook函数
	Packet() Packet
}
