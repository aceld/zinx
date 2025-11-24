package main

import (
	"os"
	"os/signal"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

// DoClientConnectedBegin is the callback function when client connects to server
// DoClientConnectedBegin 是客户端连接到服务器时的回调函数
func DoClientConnectedBegin(conn ziface.IConnection) {
	zlog.Debug("Client connected to server")

	// Set connection properties / 设置连接属性
	// Set name to "test" which is required by server / 设置name为服务器要求的"test"
	conn.SetProperty("name", "test")

	// Set additional connection properties / 设置额外的连接属性
	conn.SetProperty("version", "1.0")
	conn.SetProperty("user_id", 12345)

	// Send a test message to server / 向服务器发送测试消息
	if err := conn.SendMsg(1, []byte("Hello server, I'm client with valid name and properties")); err != nil {
		zlog.Error(err)
	}
}

// DoClientConnectedLost is the callback function when client disconnects from server
// DoClientConnectedLost 是客户端与服务器断开连接时的回调函数
func DoClientConnectedLost(conn ziface.IConnection) {
	zlog.Debug("Client disconnected from server")
	os.Exit(0)
}

// PingRouter ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call PingRouter Handle")
	// Read the data from the server
	zlog.Debug("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
}

func main() {
	// Create client / 创建客户端
	client := znet.NewClient("127.0.0.1", 8999)

	// Set connection start and stop callbacks / 设置连接开始和断开的回调
	client.SetOnConnStart(DoClientConnectedBegin)
	client.SetOnConnStop(DoClientConnectedLost)

	// Add router / 添加路由
	client.AddRouter(2, &PingRouter{})

	// Start client / 启动客户端
	zlog.Debug("Client starting")
	client.Start()

	// Wait for interrupt signal to exit / 等待中断信号以退出
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	// Stop client / 停止客户端
	client.Stop()
}