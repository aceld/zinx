package znet

import (
	"fmt"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/zdecoder"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/zpack"
	"github.com/gorilla/websocket"
	"time"
)

type WsClient struct {
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
	decoder ziface.IDecoder
	//心跳检测器
	hc ziface.IHeartbeatChecker

	dialer websocket.Dialer
}

func NewWsClient(ip string, port int, opts ...ClientOption) ziface.IClient {

	c := &WsClient{
		Ip:         ip,
		Port:       port,
		msgHandler: NewMsgHandle(),
		packet:     zpack.Factory().NewPack(ziface.ZinxDataPack), //默认使用zinx的TLV封包方式
		decoder:    zdecoder.NewTLVDecoder(),                     //默认使用zinx的TLV解码器
	}

	//应用Option设置
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// 启动客户端，发送请求且建立链接
func (c *WsClient) Start() {

	c.exitChan = make(chan struct{})

	// 将解码器添加到拦截器
	if c.decoder != nil {
		c.msgHandler.AddInterceptor(c.decoder)
	}

	//客户端将协程池关闭
	zconf.GlobalObject.WorkerPoolSize = 0

	go func() {
		addr := fmt.Sprintf("ws://%s:%d", c.Ip, c.Port)

		//创建原始Socket，得到net.Conn
		conn, _, err := c.dialer.Dial(addr, nil)
		if err != nil {
			//创建链接失败
			zlog.Ins().ErrorF("WsClient connect to server failed, err:%v", err)
			panic(err)
		}

		//创建Connection对象
		c.conn = newWsClientConn(c, conn)
		zlog.Ins().InfoF("[START] Zinx WsClient LocalAddr: %s, RemoteAddr: %s\n", conn.LocalAddr(), conn.RemoteAddr())

		//HeartBeat心跳检测
		if c.hc != nil {
			//创建链接成功，绑定链接与心跳检测器
			c.hc.BindConn(c.conn)
		}

		//启动链接
		go c.conn.Start()

		select {
		case <-c.exitChan:
			zlog.Ins().InfoF("WsClient exit.")
		}
	}()
}

// StartHeartBeat 启动心跳检测
// interval 每次发送心跳的时间间隔
func (c *WsClient) StartHeartBeat(interval time.Duration) {
	checker := NewHeartbeatChecker(interval)

	//添加心跳检测的路由
	c.AddRouter(checker.MsgID(), checker.Router())

	//client绑定心跳检测器
	c.hc = checker
}

// 启动心跳检测(自定义回调)
func (c *WsClient) StartHeartBeatWithOption(interval time.Duration, option *ziface.HeartBeatOption) {
	checker := NewHeartbeatChecker(interval)

	if option != nil {
		checker.SetHeartbeatMsgFunc(option.MakeMsg)
		checker.SetOnRemoteNotAlive(option.OnRemoteNotAlive)
		checker.BindRouter(option.HeadBeatMsgID, option.Router)
	}

	//添加心跳检测的路由
	c.AddRouter(checker.MsgID(), checker.Router())

	//client绑定心跳检测器
	c.hc = checker
}

func (c *WsClient) Stop() {
	zlog.Ins().InfoF("[STOP] Zinx WsClient LocalAddr: %s, RemoteAddr: %s\n", c.conn.LocalAddr(), c.conn.RemoteAddr())
	c.conn.Stop()
	c.exitChan <- struct{}{}
	close(c.exitChan)
}

func (c *WsClient) AddRouter(msgID uint32, router ziface.IRouter) {
	c.msgHandler.AddRouter(msgID, router)
}

func (c *WsClient) Conn() ziface.IConnection {
	return c.conn
}

// 设置该Client的连接创建时Hook函数
func (c *WsClient) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	c.onConnStart = hookFunc
}

// 设置该Client的连接断开时的Hook函数
func (c *WsClient) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	c.onConnStop = hookFunc
}

// GetOnConnStart 得到该Server的连接创建时Hook函数
func (c *WsClient) GetOnConnStart() func(ziface.IConnection) {
	return c.onConnStart
}

// 得到该Server的连接断开时的Hook函数
func (c *WsClient) GetOnConnStop() func(ziface.IConnection) {
	return c.onConnStop
}

// 获取Client绑定的数据协议封包方式
func (c *WsClient) GetPacket() ziface.IDataPack {
	return c.packet
}

// 设置Client绑定的数据协议封包方式
func (c *WsClient) SetPacket(packet ziface.IDataPack) {
	c.packet = packet
}

func (c *WsClient) GetMsgHandler() ziface.IMsgHandle {
	return c.msgHandler
}

func (c *WsClient) AddInterceptor(interceptor ziface.IInterceptor) {
	c.msgHandler.AddInterceptor(interceptor)
}

func (c *WsClient) SetDecoder(decoder ziface.IDecoder) {
	c.decoder = decoder
}
func (c *WsClient) GetLengthField() *ziface.LengthField {
	if c.decoder != nil {
		return c.decoder.GetLengthField()
	}
	return nil
}
