package ziface

type IHeartbeatChecker interface {
	SetOnRemoteNotAlive(OnRemoteNotAlive)
	//Deprecated: Heartbeat 在业务层添加对应的路由即可
	SetHeartbeatMsgFunc(HeartBeatMsgFunc)
	//Deprecated: Heartbeat 在业务层添加对应的路由即可
	SetHeartbeatFunc(HeartBeatFunc)
	//Deprecated: Heartbeat 在业务层添加对应的路由即可
	BindRouter(uint32, IRouter)
	//Deprecated: Heartbeat 在业务层添加对应的路由即可
	BindRouterSlices(uint32, ...RouterHandler)
	Start()
	Stop()
	//Deprecated: Heartbeat 在业务层处理相关的逻辑即可
	SendHeartBeatMsg() error
	BindConn(IConnection)
	Clone() IHeartbeatChecker
	//Deprecated: Heartbeat 在业务层添加对应的路由即可
	MsgID() uint32
	//Deprecated: Heartbeat 在业务层添加对应的路由即可
	Router() IRouter
	//Deprecated: Heartbeat 在业务层添加对应的路由即可
	RouterSlices() []RouterHandler
}

// User-defined method for handling heartbeat detection messages
// (用户自定义的心跳检测消息处理方法)
type HeartBeatMsgFunc func(IConnection) []byte

// HeartBeatFunc User-defined heartbeat function
// (用户自定义心跳函数)
type HeartBeatFunc func(IConnection) error

// OnRemoteNotAlive User-defined method for handling remote connections that are not alive
// 用户自定义的远程连接不存活时的处理方法
type OnRemoteNotAlive func(IConnection)

// Deprecated: Heartbeat 在业务层处理相关的逻辑
type HeartBeatOption struct {
	MakeMsg          HeartBeatMsgFunc // User-defined method for handling heartbeat detection messages(用户自定义的心跳检测消息处理方法)
	OnRemoteNotAlive OnRemoteNotAlive // User-defined method for handling remote connections that are not alive(用户自定义的远程连接不存活时的处理方法)
	HeadBeatMsgID    uint32           // User-defined ID for heartbeat detection messages(用户自定义的心跳检测消息ID)
	Router           IRouter          // User-defined business processing route for heartbeat detection messages(用户自定义的心跳检测消息业务处理路由)
	RouterSlices     []RouterHandler  //新版本的路由处理函数的集合
}

const (
	HeartBeatDefaultMsgID uint32 = 99999
)
