package main

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

// DoConnectionBegin is the callback function when connection starts
// DoConnectionBegin 是连接开始时的回调函数
func DoConnectionBegin(conn ziface.IConnection) {
	zlog.Debug("Server connection started")

	// Check client connection properties / 检查客户端连接属性
	name, err := conn.GetProperty("name")
	if err != nil {
		zlog.Error("Failed to get client name property: ", err)
		// Disconnect if property not found / 如果属性不存在则断开连接
		conn.Stop()
		return
	}

	// Verify if name is "test" / 验证name是否为"test"
	if name != "test" {
		zlog.Error("Invalid client name: ", name)
		// Disconnect if name is not valid / 如果name无效则断开连接
		conn.Stop()
		return
	}

	zlog.Debug("Client connected with valid name: ", name)

	// Get additional connection properties / 获取额外的连接属性
	version, err := conn.GetProperty("version")
	if err != nil {
		zlog.Debug("Client version property not set")
	} else {
		zlog.Debug("Client version: ", version)
	}

	userID, err := conn.GetProperty("user_id")
	if err != nil {
		zlog.Debug("Client user_id property not set")
	} else {
		zlog.Debug("Client user_id: ", userID)
	}
}

// DoConnectionLost is the callback function when connection is lost
// DoConnectionLost 是连接断开时的回调函数
func DoConnectionLost(conn ziface.IConnection) {
	zlog.Debug("Server connection lost")
}

// PingRouter ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call PingRouter Handle")
	// Read the data from the client first, then send back "ping...ping...ping"
	zlog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(2, []byte("pong-server"))
	if err != nil {
		zlog.Error(err)
	}
}

func main() {
	// Create server / 创建服务器
	s := znet.NewServer()

	// Set connection start and stop callbacks / 设置连接开始和断开的回调
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// Add router / 添加路由
	s.AddRouter(1, &PingRouter{})

	// Start server / 启动服务器
	zlog.Debug("Server starting")
	s.Serve()
}