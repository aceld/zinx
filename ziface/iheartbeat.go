package ziface

type IHeartbeatChecker interface {
	SetOnRemoteNotAlive(OnRemoteNotAlive)
	SetHeartbeatMsgFunc(HeartBeatMsgFunc)
	BindRouter(uint32, IRouter)
	Start()
	Stop()
	SendHeartBeatMsg() error
	BindConn(IConnection)
	Clone() IHeartbeatChecker
	MsgID() uint32
	Router() IRouter
}

// 用户自定义的心跳检测消息处理方法
type HeartBeatMsgFunc func(IConnection) []byte

// 用户自定义的远程连接不存活时的处理方法
type OnRemoteNotAlive func(IConnection)

type HeartBeatOption struct {
	MakeMsg          HeartBeatMsgFunc //用户自定义的心跳检测消息处理方法
	OnRemoteNotAlive OnRemoteNotAlive //用户自定义的远程连接不存活时的处理方法
	HeadBeatMsgID    uint32           //用户自定义的心跳检测消息ID
	Router           IRouter          //用户自定义的心跳检测消息业务处理路由
}

const (
	HeartBeatDefaultMsgID uint32 = 99999
)
