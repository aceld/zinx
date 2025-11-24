package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

// DoClientConnectedBegin is the callback function when client connects to server
// DoClientConnectedBegin 是客户端连接到服务器时的回调函数
func DoClientConnectedBegin(conn ziface.IConnection) {
	zlog.Debug("Client connected to server")

	// Set connection properties with invalid name / 设置带有无效名称的连接属性
	conn.SetProperty("name", "invalid_name")

	// Try to send a test message to server / 尝试向服务器发送测试消息
	// This message may not be sent because the server may have already disconnected us / 此消息可能无法发送，因为服务器可能已经断开了我们的连接
	if err := conn.SendMsg(1, []byte("Hello server, I'm client with invalid name")); err != nil {
		zlog.Error(err)
	}
}

// DoClientConnectedLost is the callback function when client disconnects from server
// DoClientConnectedLost 是客户端与服务器断开连接时的回调函数
func DoClientConnectedLost(conn ziface.IConnection) {
	zlog.Debug("Client disconnected from server (expected, since we used invalid name)")
	os.Exit(0)
}

func main() {
	// Create client / 创建客户端
	client := znet.NewClient("127.0.0.1", 8999)

	// Set connection start and stop callbacks / 设置连接开始和断开的回调
	client.SetOnConnStart(DoClientConnectedBegin)
	client.SetOnConnStop(DoClientConnectedLost)

	// Start client / 启动客户端
	zlog.Debug("Client (with invalid name) starting")
	client.Start()

	// Wait for interrupt signal to exit / 等待中断信号以退出
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Wait for a short time to see if we get disconnected by server / 等待一段时间，看看是否会被服务器断开连接
	select {
	case <-ch:
		// Stop client / 停止客户端
		client.Stop()
	case <-time.After(5 * time.Second):
		zlog.Warn("Client is still connected, which is unexpected")
		client.Stop()
		os.Exit(1)
	}
}
