package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/aceld/zinx/examples/zinx_test_stopChanelPanic/router"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

func onClientStart(conn ziface.IConnection) {
	fmt.Println("[Client] 连接建立成功，开始发送测试消息...")

	// 发送测试消息
	go func() {
		time.Sleep(500 * time.Millisecond) // 等待一下确保连接稳定
		err := conn.SendMsg(1, []byte("Test Panic"))
		if err != nil {
			fmt.Println("[Client] 发送失败:", err)
		} else {
			fmt.Println("[Client] 测试消息已发送")
		}
	}()
}

func main() {
	client := znet.NewClient("127.0.0.1", 8999)
	client.AddRouter(1, &router.PanicTestRouter{znet.BaseRouter{}})
	client.SetOnConnStart(onClientStart)
	// Start the client
	client.Start()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
	client.Stop()
	time.Sleep(time.Second * 2)
}
