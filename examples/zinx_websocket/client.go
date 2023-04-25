package main

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_client/c_router"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
	"time"
)

type PositionClientRouter struct {
	znet.BaseRouter
}

func (this *PositionClientRouter) Handle(request ziface.IRequest) {

}

// 客户端自定义业务
func business(conn ziface.IConnection) {

	for {
		err := conn.SendMsg(1, []byte("ping ping ping ..."))
		if err != nil {
			fmt.Println(err)

		}
		time.Sleep(1 * time.Second)
	}
}

// 创建连接的时候执行
func DoClientConnectedBegin(conn ziface.IConnection) {
	//设置两个链接属性，在连接创建之后
	conn.SetProperty("Name", "刘丹冰")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	go business(conn)
}

func wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
}

func main() {
	// Create a Client.
	client := znet.NewWsClient("127.0.0.1", 9000)

	// Add business logic for when the connection is first established.(添加首次建立链接时的业务)
	client.SetOnConnStart(DoClientConnectedBegin)
	// Register business routing for receiving messages from the server.(注册收到服务器消息业务路由)
	client.AddRouter(2, &c_router.PingRouter{})
	client.AddRouter(3, &c_router.HelloRouter{})
	// Start the client.
	client.Start()
	select {
	case err := <-client.GetErrChan():
		// Handle the errors returned by the client.(处理客户端返回的错误)
		zlog.Ins().ErrorF("client err:%v", err)
	}

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
	// Clean up the client.(清理客户端)
	client.Stop()
	time.Sleep(time.Second * 2)
}
