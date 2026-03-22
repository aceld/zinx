package main

import (
	"fmt"
	"time"

	"github.com/aceld/zinx/v3/examples/zinx_closecallback/router"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

// DoConnectionBegin is the callback function when connection starts
// DoConnectionBegin 是连接开始时的回调函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("Server connection started")

	// Set connection properties / 设置连接属性
	conn.SetProperty("StartTime", time.Now())

	// Add close callback function - cleanup resources / 添加关闭回调函数 - 清理资源
	conn.AddCloseCallback("cleanup", "resources", func() {
		fmt.Printf("Cleanup resources for connection %d\n", conn.GetConnID())
	})

	// Add close callback function - record statistics / 添加关闭回调函数 - 记录统计信息
	conn.AddCloseCallback("stats", "record", func() {
		if startTime, err := conn.GetProperty("StartTime"); err == nil {
			duration := time.Since(startTime.(time.Time))
			fmt.Printf("Connection %d duration: %v\n", conn.GetConnID(), duration)
		}
	})

	// Add close callback function - notify other components / 添加关闭回调函数 - 通知其他组件
	conn.AddCloseCallback("notification", "broadcast", func() {
		fmt.Printf("Notify other components: connection %d disconnected\n", conn.GetConnID())
	})
}

// DoConnectionLost is the callback function when connection is lost
// DoConnectionLost 是连接断开时的回调函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("Server connection lost")
}

func main() {
	// Create server / 创建服务器
	s := znet.NewServer()

	// Set connection start and stop callbacks / 设置连接开始和断开的回调
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// Add router / 添加路由
	s.AddRouter(1, &router.PingRouter{})

	// Start server / 启动服务器
	fmt.Println("Server starting")
	s.Serve()
}
