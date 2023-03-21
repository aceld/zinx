package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
	"time"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from server : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	/*
		if err := request.GetConnection().SendBuffMsg(1, []byte("Hello[from client]")); err != nil {
			zlog.Error(err)
		}
	*/
}

// 客户端自定义业务
func business(conn ziface.IConnection) {

	for {
		err := conn.SendMsg(1, []byte("Ping...[FromClient]"))
		if err != nil {
			fmt.Println(err)
			break
		}

		time.Sleep(30 * time.Second)
	}
}

// 创建连接的时候执行
func DoClientConnectedBegin(conn ziface.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")

	//设置两个链接属性，在连接创建之后
	conn.SetProperty("Name", "刘丹冰")
	conn.SetProperty("Home", "https://yuque.com/aceld")

	go business(conn)
}

// 连接断开的时候执行
func DoClientConnectedLost(conn ziface.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}

	fmt.Println("DoClientConnectedLost is Called ... ")
}

func main() {
	//创建一个Client句柄，使用Zinx的API
	client := znet.NewClient("127.0.0.1", 8999)

	//添加首次建立链接时的业务
	client.SetOnConnStart(DoClientConnectedBegin)
	client.SetOnConnStop(DoClientConnectedLost)

	//注册收到服务器消息业务路由
	client.AddRouter(0, &PingRouter{})

	//启动心跳检测
	client.StartHeartBeat(5 * time.Second)

	//启动客户端client
	client.Start()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)
}
