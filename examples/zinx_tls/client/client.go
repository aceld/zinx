package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
	"time"
)

// PongRouter pong test 自定义路由
type PongRouter struct {
	znet.BaseRouter
}

// Handle Pong Handle
func (this *PongRouter) Handle(request ziface.IRequest) {

	zlog.Debug("Call PongRouter Handle")
	//先读取服务器返回的数据
	zlog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

}

func wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
}

func main() {
	// Create a TLS client.
	c := znet.NewTLSClient("127.0.0.1", 8899)

	c.SetOnConnStart(func(connection ziface.IConnection) {
		go func() {
			for {
				err := connection.SendMsg(1, []byte("Ping with TLS"))

				if err != nil {
					fmt.Println(err)
					break
				}

				time.Sleep(1 * time.Second)
			}
		}()

	})

	c.AddRouter(2, &PongRouter{})

	c.Start()

	wait()
}
