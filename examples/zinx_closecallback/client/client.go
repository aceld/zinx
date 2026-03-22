package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/aceld/zinx/v3/examples/zinx_closecallback/router"
	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

// business handles the main business logic for sending ping messages
// business 处理发送ping消息的主要业务逻辑
func business(conn ziface.IConnection) {
	for i := 0; i < 3; i++ {
		err := conn.SendMsg(1, []byte(fmt.Sprintf("Ping %d", i+1)))
		if err != nil {
			fmt.Println("SendMsg error:", err)
			break
		}

		time.Sleep(1 * time.Second)
	}

	// Actively disconnect after sending is complete / 发送完成后主动断开连接
	fmt.Println("Client actively disconnects")
	conn.Stop()
}

// DoClientConnectedBegin is the callback function when connection starts
// DoClientConnectedBegin 是连接开始时的回调函数
func DoClientConnectedBegin(conn ziface.IConnection) {
	fmt.Println("Client connection started")

	// Set connection properties / 设置连接属性
	conn.SetProperty("StartTime", time.Now())

	// Add close callback function - record connection statistics / 添加关闭回调函数 - 记录连接统计
	conn.AddCloseCallback("stats", "connection-stats", func() {
		if startTime, err := conn.GetProperty("StartTime"); err == nil {
			duration := time.Since(startTime.(time.Time))
			fmt.Printf("Client connection duration: %v\n", duration)
		}
	})

	// Start business processing / 启动业务处理
	go business(conn)
}

// DoClientConnectedLost is the callback function when connection is lost
// DoClientConnectedLost 是连接断开时的回调函数
func DoClientConnectedLost(conn ziface.IConnection) {
	fmt.Println("Client connection lost")
}

func main() {
	// Create client / 创建客户端
	client := znet.NewClient("127.0.0.1", 8999)

	// Set connection start and stop callbacks / 设置连接开始和断开的回调
	client.SetOnConnStart(DoClientConnectedBegin)
	client.SetOnConnStop(DoClientConnectedLost)

	// Add router / 添加路由
	client.AddRouter(0, &router.PingRouter{})

	// Start client / 启动客户端
	fmt.Println("Client starting")
	client.Start()

	// Wait for interrupt signal / 等待中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
