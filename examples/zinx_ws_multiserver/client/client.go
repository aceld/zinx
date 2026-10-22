package main

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_client/c_router"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func business(connId string, conn ziface.IConnection) {
	for i := 0; i < 5; i++ {
		err := conn.SendMsg(1, []byte(fmt.Sprintf("ping from %s", connId)))
		if err != nil {
			fmt.Printf("[%s] send error: %v\n", connId, err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Connect to both WebSocket servers to verify they work independently
	client1 := znet.NewWsClient("127.0.0.1", 9001)
	client2 := znet.NewWsClient("127.0.0.1", 9002)

	client1.AddRouter(2, &c_router.PingRouter{})
	client2.AddRouter(2, &c_router.PingRouter{})

	client1.SetOnConnStart(func(conn ziface.IConnection) {
		go business("client1", conn)
	})
	client2.SetOnConnStart(func(conn ziface.IConnection) {
		go business("client2", conn)
	})

	client1.Start()
	client2.Start()

	fmt.Println("Both clients connected. Press Ctrl+C to exit.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client1.Stop()
	client2.Stop()
	fmt.Println("===exit===")
}
