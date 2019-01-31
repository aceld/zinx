package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool

	//该连接的处理方法router
	Router  ziface.IRouter

	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool

	//给缓冲队列发送数据的channel，
	// 如果向缓冲队列发送数据，那么把数据发送到这个channel下
//	SendBuffChan chan []byte

}


//创建连接的方法
func NewConntion(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection{
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router: router,
		ExitBuffChan: make(chan bool, 1),
		//		SendBuffChan: make(chan []byte, 512),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for  {
		//读取我们最大的数据到buf中
		buf := make([]byte, utils.GlobalObject.MaxPacketSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			continue
		}
		//得到当前客户端请求的Request数据
		req := Request{
			conn:c,
			data:buf,
		}
		//从路由Routers 中找到注册绑定Conn的对应Handle
		go func (request ziface.IRequest) {
			//执行注册的路由方法
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

		//if err := c.handleAPI(c.Conn, buf, cnt); err !=nil {
		//	fmt.Println("connID ", c.ConnID, " handle is error")
		//	c.ExitBuffChan <- true
		//	return
		//}
	}
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {

	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	for {
		select {
		case <- c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}

	//1 开启用于写回客户端数据流程的Goroutine
	//2 开启用户从客户端读取数据流程的Goroutine
}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	// 关闭socket链接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- true

	//关闭该链接全部管道
	close(c.ExitBuffChan)
	//close(c.SendBuffChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32{
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将数据发送数据给远程的TCP客户端
func (c *Connection) Send(data []byte) error {
	return nil
}

//将数据发送给缓冲队列，通过专门从缓冲队列读数据的go写给客户端
func (c *Connection) SendBuff(data []byte) error {
	return nil
}
