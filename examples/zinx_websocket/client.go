package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/aceld/zinx/zpack"
	"github.com/gorilla/websocket"
	"log"
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
func business(conn *websocket.Conn) {
	pack := zpack.Factory().NewPack(ziface.ZinxDataPack)
	for {
		msgPackage := zpack.NewMsgPackage(0, []byte("ping ping ping ..."))
		msgData, err := pack.Pack(msgPackage)
		err = conn.WriteMessage(websocket.BinaryMessage, msgData)
		if err != nil {
			fmt.Println(err)
			break
		}

		time.Sleep(1 * time.Second)
	}
}

// 创建连接的时候执行
func DoClientConnectedBegin(conn ziface.IConnection) {
	//设置两个链接属性，在连接创建之后
	conn.SetProperty("Name", "刘丹冰")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	//go business(conn)
}

func wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
}

func main() {
	//创建一个Client句柄，使用Zinx的API
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://localhost:8999/ws", nil)

	if err != nil {
		log.Println(err)
		return
	}
	//离开作用域关闭连接，go 的常规操作
	defer conn.Close()

	//定时向客户端发送数据
	go business(conn)

	wait()
}
