package znet

import (
	"github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/zpack"
	"math"
	"net"
	"time"
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
	//断粘包解码器
	LengthField ziface.LengthField
}

func NewClient(ip string, port int, opts ...ClientOption) ziface.IClient {

	c := &Client{
		Ip:         ip,
		Port:       port,
		msgHandler: NewMsgHandle(),
		packet:     zpack.Factory().NewPack(ziface.ZinxDataPack), //默认使用zinx的TLV封包方式
		// +---------------+---------------+---------------+
		// |      Tag      |    Length     |     Value     |
		// | uint32(4byte) | uint32(4byte) |     n byte    |
		// +---------------+---------------+---------------+
		// Tag：   uint32类型，占4字节
		// Length：uint32类型，占4字节，Length标记Value长度
		// Value： 占n字节
		//
		//说明:
		//    lengthFieldOffset   = 4            (Length的字节位索引下标是4) 长度字段的偏差
		//    lengthFieldLength   = 4            (Length是4个byte) 长度字段占的字节数
		//    lengthAdjustment    = 0            (Length只表示Value长度，程序只会读取Length个字节就结束，后面没有来，故为0，若Value后面还有crc占2字节的话，那么此处就是2。若Length标记的是Tag+Length+Value总长度，那么此处是-8)
		//    initialBytesToStrip = 0            (这个0表示返回完整的协议内容Tag+Length+Value，如果只想返回Value内容，去掉Tag的4字节和Length的4字节，此处就是8) 从解码帧中第一次去除的字节数
		//    maxFrameLength      = 2^32 + 4 + 4 (Length为uint32类型，故2^32次方表示Value最大长度，此外Tag和Length各占4字节)
		//默认使用TLV封包方式
		LengthField: ziface.LengthField{
			MaxFrameLength:      math.MaxUint32 + 4 + 4,
			LengthFieldOffset:   4,
			LengthFieldLength:   4,
			LengthAdjustment:    0,
			InitialBytesToStrip: 0,
		},
	}

	//应用Option设置
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (this *Client) AddInterceptor(interceptor ziface.Interceptor) {
	this.msgHandler.AddInterceptor(interceptor)
}

func (this *Client) GetLengthField() ziface.LengthField {
	return this.LengthField
}

// 启动客户端，发送请求且建立链接
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

// 启动心跳检测
func (c *Client) StartHeartBeat(interval time.Duration) {
	checker := NewHeartbeatCheckerC(interval, c)

	//添加心跳检测的路由
	c.AddRouter(checker.msgID, checker.router)

	go checker.Start()
}

// 启动心跳检测(自定义回调)
func (c *Client) StartHeartBeatWithOption(interval time.Duration, option *ziface.HeartBeatOption) {
	checker := NewHeartbeatCheckerC(interval, c)

	if option != nil {
		checker.SetHeartbeatMsgFunc(option.MakeMsg)
		checker.SetOnRemoteNotAlive(option.OnRemoteNotAlive)
		checker.BindRouter(option.HeadBeatMsgID, option.Router)
	}

	//添加心跳检测的路由
	c.AddRouter(checker.msgID, checker.router)

	go checker.Start()
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

// 设置该Client的连接创建时Hook函数
func (c *Client) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	c.onConnStart = hookFunc
}

// 设置该Client的连接断开时的Hook函数
func (c *Client) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	c.onConnStop = hookFunc
}

// GetOnConnStart 得到该Server的连接创建时Hook函数
func (c *Client) GetOnConnStart() func(ziface.IConnection) {
	return c.onConnStart
}

// 得到该Server的连接断开时的Hook函数
func (c *Client) GetOnConnStop() func(ziface.IConnection) {
	return c.onConnStop
}

// 获取Client绑定的数据协议封包方式
func (c *Client) GetPacket() ziface.IDataPack {
	return c.packet
}

// 设置Client绑定的数据协议封包方式
func (c *Client) SetPacket(packet ziface.IDataPack) {
	c.packet = packet
}

func (c *Client) GetMsgHandler() ziface.IMsgHandle {
	return c.msgHandler
}
