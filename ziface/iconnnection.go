package ziface

import "net"

//定义连接接口
type IConnection interface {
	//启动连接，让当前连接开始工作
	Start()
	//停止连接，结束当前连接状态M
	Stop()
	//从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn
	//获取当前连接ID
	GetConnID() uint32
	//获取远程客户端地址信息
	RemoteAddr() net.Addr
	//直接将数据发送数据给远程的TCP客户端
	Send(data []byte) error
	//将数据发送给缓冲队列，通过专门从缓冲队列读数据的go写给客户端
	SendBuff(data []byte) error
}


