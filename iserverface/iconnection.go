package iserverface

import (
	"github.com/gorilla/websocket"
	"net"
)

type IConnection interface {
	//启动链接开始工作
	Start()
	//关闭链接停止工作
	Close()
	//获取websocket链接
	GetConnection() *websocket.Conn
	//获取当前连接ID
	GetConnID() uint64
	//获取远程客户端地址信息
	RemoteAddr() net.Addr
	//发送数据
	SendMessage(msgType int,msgData []byte) error
	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)

}