/**
* @Author: Aceld
* @Date: 2020/12/24 00:24
* @Mail: danbing.at@gmail.com
*    zinx server demo
 */
package main

import (
	"github.com/aceld/zinx/examples/zinx_server/s_router"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

// DoConnectionBegin Executed when creating a connection.
// 创建连接的时候执行
func DoConnectionBegin(conn ziface.IConnection) {
	zlog.Ins().InfoF("DoConnecionBegin is Called ...")

	//设置两个链接属性，在连接创建之后
	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://www.kancloud.cn/@aceld")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		zlog.Error(err)
	}
}

// 连接断开的时候执行
// DoConnectionLost Executed when the connection is closed.
func DoConnectionLost(conn ziface.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性
	// Query the Name and Home properties of conn before destroying the connection.
	if name, err := conn.GetProperty("Name"); err == nil {
		zlog.Ins().InfoF("Conn Property Name = %v", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		zlog.Ins().InfoF("Conn Property Home = %v", home)
	}

	zlog.Ins().InfoF("Conn is Lost")
}

func main() {
	// Create a server
	s := znet.NewServer()

	// Register a hook callback function for the connection
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// Configure routing.
	s.AddRouter(100, &s_router.PingRouter{})
	s.AddRouter(1, &s_router.HelloZinxRouter{})

	// Start Service
	s.Serve()
}
