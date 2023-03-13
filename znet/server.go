package znet

import (
	"errors"
	"fmt"
	"github.com/aceld/zinx/zlog"
	"math"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zpack"
)

var zinxLogo = `                                        
              ██                        
              ▀▀                        
 ████████   ████     ██▄████▄  ▀██  ██▀ 
     ▄█▀      ██     ██▀   ██    ████   
   ▄█▀        ██     ██    ██    ▄██▄   
 ▄██▄▄▄▄▄  ▄▄▄██▄▄▄  ██    ██   ▄█▀▀█▄  
 ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀  ▀▀    ▀▀  ▀▀▀  ▀▀▀ 
                                        `
var topLine = `┌──────────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└──────────────────────────────────────────────────────┘`

// Server 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler ziface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr ziface.IConnManager
	//该Server的连接创建时Hook函数
	onConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	onConnStop func(conn ziface.IConnection)
	//数据报文封包方式
	packet ziface.IDataPack
	//异步捕获链接关闭状态
	exitChan chan struct{}
	//断粘包解码器
	LengthField ziface.LengthField
}

func (this *Server) GetLengthField() ziface.LengthField {
	return this.LengthField
}

// NewServer 创建一个服务器句柄
func NewServer(opts ...Option) ziface.IServer {
	printLogo()

	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		exitChan:   nil,
		//默认使用zinx的TLV封包方式
		packet: zpack.Factory().NewPack(ziface.ZinxDataPack),
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

	for _, opt := range opts {
		opt(s)
	}

	//提示当前配置信息
	utils.GlobalObject.Show()

	return s
}

// NewServer 创建一个服务器句柄
func NewUserConfServer(config *utils.Config, opts ...Option) ziface.IServer {
	//打印logo
	printLogo()

	s := &Server{
		Name:       config.Name,
		IPVersion:  config.TcpVersion,
		IP:         config.Host,
		Port:       config.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		exitChan:   nil,
		packet:     zpack.Factory().NewPack(ziface.ZinxDataPack),
		//默认使用TLV封包方式
		LengthField: ziface.LengthField{
			MaxFrameLength:      math.MaxUint32 + 4 + 4,
			LengthFieldOffset:   4,
			LengthFieldLength:   4,
			LengthAdjustment:    0,
			InitialBytesToStrip: 0,
		},
	}
	//更替打包方式
	for _, opt := range opts {
		opt(s)
	}
	//刷新用户配置到全局配置变量
	utils.UserConfToGlobal(config)

	//提示当前配置信息
	utils.GlobalObject.Show()

	return s
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

func (this *Server) AddInterceptor(interceptor ziface.Interceptor) {
	this.msgHandler.AddInterceptor(interceptor)
}

// Start 开启网络服务
func (s *Server) Start() {
	zlog.Ins().InfoF("[START] Server name: %s,listenner at IP: %s, Port %d is starting", s.Name, s.IP, s.Port)
	s.exitChan = make(chan struct{})

	//开启一个go去做服务端Linster业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			zlog.Ins().ErrorF("[START] resolve tcp addr err: %v\n", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}

		//已经监听成功
		zlog.Ins().InfoF("[START] start Zinx server  %s succ, now listenning...", s.Name)

		//TODO server.go 应该有一个自动生成ID的方法
		var cID uint32
		cID = 0

		go func() {
			//3 启动server网络连接业务
			for {
				//3.1 设置服务器最大连接控制,如果超过最大连接，则等待
				if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
					zlog.Ins().InfoF("Exceeded the maxConnNum:%d, Wait:%d", utils.GlobalObject.MaxConn, AcceptDelay.duration)
					AcceptDelay.Delay()
					continue
				}

				//3.2 阻塞等待客户端建立连接请求
				conn, err := listener.AcceptTCP()
				if err != nil {
					//Go 1.16+
					if errors.Is(err, net.ErrClosed) {
						zlog.Ins().ErrorF("Listener closed")
						return
					}
					zlog.Ins().ErrorF("Accept err: %v", err)
					AcceptDelay.Delay()
					continue
				}

				AcceptDelay.Reset()

				//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
				dealConn := newServerConn(s, conn, cID)
				cID++

				//3.4 启动当前链接的处理业务
				go dealConn.Start()
			}
		}()

		select {
		case <-s.exitChan:
			err := listener.Close()
			if err != nil {
				zlog.Ins().ErrorF("listener close err: %v", err)
			}
		}
	}()
}

// Stop 停止服务
func (s *Server) Stop() {
	zlog.Ins().InfoF("[STOP] Zinx server , name %s", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
	s.exitChan <- struct{}{}
	close(s.exitChan)
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	//select {}
	c := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill信号
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	zlog.Ins().InfoF("[SERVE] Zinx server , name %s, Serve Interrupt, signal = %v", s.Name, sig)
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.onConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.onConnStop = hookFunc
}

// GetOnConnStart 得到该Server的连接创建时Hook函数
func (s *Server) GetOnConnStart() func(ziface.IConnection) {
	return s.onConnStart
}

// 得到该Server的连接断开时的Hook函数
func (s *Server) GetOnConnStop() func(ziface.IConnection) {
	return s.onConnStop
}

func (s *Server) GetPacket() ziface.IDataPack {
	return s.packet
}

func (s *Server) SetPacket(packet ziface.IDataPack) {
	s.packet = packet
}

func (s *Server) GetMsgHandler() ziface.IMsgHandle {
	return s.msgHandler
}

// StartHeartBeat 启动心跳检测
// interval 每次发送心跳的时间间隔
// msgID    心跳检测消息的msgID
func (s *Server) StartHeartBeat(interval time.Duration) {
	checker := NewHeartbeatCheckerS(interval, s)

	//添加心跳检测的路由
	s.AddRouter(checker.msgID, checker.router)

	go checker.Start()
}

// StartHeartBeatWithFunc 启动心跳检测
// msgFunc  心跳检测消息的自定义生成函数,不需要自定义可以传nil
// notAlive 检测到地方停止心跳的自定义处理函数,不需要自定义可以传nil
func (s *Server) StartHeartBeatWithOption(interval time.Duration, option *ziface.HeartBeatOption) {
	checker := NewHeartbeatCheckerS(interval, s)

	if option != nil {
		checker.SetHeartbeatMsgFunc(option.MakeMsg)
		checker.SetOnRemoteNotAlive(option.OnRemoteNotAlive)
		checker.BindRouter(option.HeadBeatMsgID, option.Router)
	}

	//添加心跳检测的路由
	s.AddRouter(checker.msgID, checker.router)

	go checker.Start()
}

func printLogo() {
	fmt.Println(zinxLogo)
	fmt.Println(topLine)
	fmt.Println(fmt.Sprintf("%s [Github] https://github.com/aceld                    %s", borderLine, borderLine))
	fmt.Println(fmt.Sprintf("%s [tutorial] https://www.yuque.com/aceld/npyr8s/bgftov %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
}

func init() {}
