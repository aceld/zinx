package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zinx_app_demo/mmo_game/pb"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
	"os"
	"os/signal"
	"time"
)

type PositionClientRouter struct {
	znet.BaseRouter
}

func (this *PositionClientRouter) Handle(request ziface.IRequest) {
	fmt.Println("Handle....")

	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Position Unmarshal error ", err, " data = ", request.GetData())
		return
	}

	fmt.Printf("recv from server : msgId=%+v, data=%+v\n", request.GetMsgID(), msg)
}

// 客户端自定义业务
func business(conn ziface.IConnection) {

	for {

		msg := &pb.Position{}
		msg.X = 1
		msg.Y = 2
		msg.Z = 3
		msg.V = 4

		data, err := proto.Marshal(msg)
		if err != nil {
			fmt.Println("proto Marshal error = ", err, " msg = ", msg)
			break
		}

		err = conn.SendMsg(0, data)
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

	go business(conn)
}

func wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
}

func main() {
	//创建一个Client句柄，使用Zinx的API
	client := znet.NewClient("127.0.0.1", 8999)

	//添加首次建立链接时的业务
	client.SetOnConnStart(DoClientConnectedBegin)

	//注册收到服务器消息业务路由
	client.AddRouter(0, &PositionClientRouter{})

	//启动客户端client
	client.Start()

	wait()
}
