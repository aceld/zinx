package iserverface

import "github.com/gin-gonic/gin"

//定义服务器接口
type IServer interface {
	//启动服务器方法
	Start(c *gin.Context)
	//停止服务器方法
	Stop()
	//开启业务服务方法
	Serve(c *gin.Context)
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(msgId string, router IRouter)
	//得到链接管理
	GetConnMgr() IConnMgr
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}