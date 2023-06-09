// @Title  iclient.go
// @Description  Provides all interface declarations for the Client abstraction layer.
// @Author  Aceld - 2023-2-28

package ziface

import "time"

type IClient interface {
	Restart()
	Start()
	Stop()
	AddRouter(msgID uint32, router IRouter)
	Conn() IConnection

	// SetOnConnStart Set the Hook function to be called when a connection is created for this Client
	// (设置该Client的连接创建时Hook函数)
	SetOnConnStart(func(IConnection))

	// SetOnConnStop Set the Hook function to be called when a connection is closed for this Client
	// (设置该Client的连接断开时的Hook函数)
	SetOnConnStop(func(IConnection))

	// GetOnConnStart Get the Hook function that is called when a connection is created for this Client
	// (获取该Client的连接创建时Hook函数)
	GetOnConnStart() func(IConnection)

	// GetOnConnStop Get the Hook function that is called when a connection is closed for this Client
	// (设置该Client的连接断开时的Hook函数)
	GetOnConnStop() func(IConnection)

	// GetPacket Get the data protocol packet binding method for this Client
	// (获取Client绑定的数据协议封包方式)
	GetPacket() IDataPack

	// SetPacket Set the data protocol packet binding method for this Client
	// (设置Client绑定的数据协议封包方式)
	SetPacket(IDataPack)

	// GetMsgHandler Get the message handling module bound to this Client
	// (获取Client绑定的消息处理模块)
	GetMsgHandler() IMsgHandle

	// StartHeartBeat Start heartbeat detection(启动心跳检测)
	StartHeartBeat(time.Duration)

	// StartHeartBeatWithOption Start heartbeat detection with custom callbacks 启动心跳检测(自定义回调)
	StartHeartBeatWithOption(time.Duration, *HeartBeatOption)

	// GetLengthField Get the length field of this Client
	GetLengthField() *LengthField

	// SetDecoder Set the decoder for this Client 设置解码器
	SetDecoder(IDecoder)

	// AddInterceptor Add an interceptor for this Client 添加拦截器
	AddInterceptor(IInterceptor)

	// Get the error channel for this Client 获取客户端错误管道
	GetErrChan() chan error

	// Set the name of this Clien
	// 设置客户端Client名称
	SetName(string)

	// Get the name of this Client
	// 获取客户端Client名称
	GetName() string
}
