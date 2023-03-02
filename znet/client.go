package znet

import (
	"github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/zpack"
	"net"
)

type Client struct {
	//目标链接服务器的IP
	Ip string
	//目标链接服务器的端口
	Port int
	//客户端链接
	conn ziface.IConnection
	//该client的连接创建时Hook函数
	onConnStart func(conn ziface.IConnection)
	//该client的连接断开时的Hook函数
	onConnStop func(conn ziface.IConnection)
	//数据报文封包方式
	packet ziface.IDataPack
	//异步捕获链接关闭状态
	exitChan chan struct{}
	//消息管理模块
	msgHandler ziface.IMsgHandle
}

func NewClient(ip string, port int, opts ...ClientOption) ziface.IClient {

	c := &Client{
		Ip:         ip,
		Port:       port,
		msgHandler: NewMsgHandle(),
		packet:     zpack.Factory().NewPack(ziface.ZinxDataPack), //默认使用zinx的TLV封包方式
	}

	//应用Option设置
	for _, opt := range opts {
		opt(c)
	}

	return c
}

//启动客户端，发送请求且建立链接
func (c *Client) Start() {
	c.exitChan = make(chan struct{})

	//客户端将协程池关闭
	utils.GlobalObject.WorkerPoolSize = 0

	go func() {
		addr := &net.TCPAddr{
			IP:   net.ParseIP(c.Ip),
			Port: c.Port,
			Zone: "", //for ipv6, ignore
		}

		//创建原始Socket，得到net.Conn
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			//创建链接失败
			zlog.Ins().ErrorF("client connect to server failed, err:%v", err)
			panic(err)
			return
		}

		//创建Connection对象
		c.conn = newClientConn(c, conn)
		zlog.Ins().InfoF("[START] Zinx Client LocalAddr: %s, RemoteAddr: %s\n", conn.LocalAddr(), conn.RemoteAddr())

		//启动链接
		go c.conn.Start()

		select {
		case <-c.exitChan:
			zlog.Ins().InfoF("client exit.")
		}
	}()
}

func (c *Client) Stop() {
	zlog.Ins().InfoF("[STOP] Zinx Client LocalAddr: %s, RemoteAddr: %s\n", c.conn.LocalAddr(), c.conn.RemoteAddr())
	c.conn.Stop()
}

func (c *Client) AddRouter(msgID uint32, router ziface.IRouter) {
	c.msgHandler.AddRouter(msgID, router)
}

func (c *Client) Conn() ziface.IConnection {
	return c.conn
}

//设置该Client的连接创建时Hook函数
func (c *Client) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	c.onConnStart = hookFunc
}

//设置该Client的连接断开时的Hook函数
func (c *Client) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	c.onConnStop = hookFunc
}

//GetOnConnStart 得到该Server的连接创建时Hook函数
func (c *Client) GetOnConnStart() func(ziface.IConnection) {
	return c.onConnStart
}

//得到该Server的连接断开时的Hook函数
func (c *Client) GetOnConnStop() func(ziface.IConnection) {
	return c.onConnStop
}

//获取Client绑定的数据协议封包方式
func (c *Client) GetPacket() ziface.IDataPack {
	return c.packet
}

//设置Client绑定的数据协议封包方式
func (c *Client) SetPacket(packet ziface.IDataPack) {
	c.packet = packet
}

func (c *Client) GetMsgHandler() ziface.IMsgHandle {
	return c.msgHandler
}
