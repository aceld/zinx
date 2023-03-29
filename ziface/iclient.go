// Package ziface 主要提供zinx全部抽象层接口定义.
// 包括:
//
//			IServer 服务mod接口
//			IRouter 路由mod接口
//			IConnection 连接mod层接口
//	     IMessage 消息mod接口
//			IDataPack 消息拆解接口
//	     IMsgHandler 消息处理及协程池接口
//	     IClient 客户端接口
//
// 当前文件描述:
// @Title  iclient.go
// @Description  提供Client抽象层全部接口声明
// @Author  Aceld - 2023-2-28
package ziface

import "time"

type IClient interface {
	Start()
	Stop()
	AddRouter(msgID uint32, router IRouter)
	Conn() IConnection
	SetOnConnStart(func(IConnection))                         //设置该Client的连接创建时Hook函数
	SetOnConnStop(func(IConnection))                          //设置该Client的连接断开时的Hook函数
	GetOnConnStart() func(IConnection)                        //获取该Client的连接创建时Hook函数
	GetOnConnStop() func(IConnection)                         //设置该Client的连接断开时的Hook函数
	GetPacket() IDataPack                                     //获取Client绑定的数据协议封包方式
	SetPacket(IDataPack)                                      //设置Client绑定的数据协议封包方式
	GetMsgHandler() IMsgHandle                                //获取Client绑定的消息处理模块
	StartHeartBeat(time.Duration)                             //启动心跳检测
	StartHeartBeatWithOption(time.Duration, *HeartBeatOption) //启动心跳检测(自定义回调)
	GetLengthField() *LengthField
	SetDecoder(IDecoder)
	AddInterceptor(IInterceptor)
}
