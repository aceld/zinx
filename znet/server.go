package znet

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/aceld/zinx/logo"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/zdecoder"
	"github.com/aceld/zinx/zlog"
	"github.com/gorilla/websocket"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zpack"
	"github.com/xtaci/kcp-go"
)

// Server interface implementation, defines a Server service class
// (接口实现，定义一个Server服务类)
type Server struct {
	// Name of the server (服务器的名称)
	Name string
	//tcp4 or other
	IPVersion string
	// IP version (e.g. "tcp4") - 服务绑定的IP地址
	IP string
	// IP address the server is bound to (服务绑定的端口)
	Port int
	// 服务绑定的websocket 端口 (Websocket port the server is bound to)
	WsPort int
	// 服务绑定的kcp 端口 (kcp port the server is bound to)
	KcpPort int

	// Current server's message handler module, used to bind MsgID to corresponding processing methods
	// (当前Server的消息管理模块，用来绑定MsgID和对应的处理方法)
	msgHandler ziface.IMsgHandle

	// Routing mode (路由模式)
	RouterSlicesMode bool

	// Current server's connection manager (当前Server的链接管理器)
	ConnMgr ziface.IConnManager

	// Hook function called when a new connection is established
	// (该Server的连接创建时Hook函数)
	onConnStart func(conn ziface.IConnection)

	// Hook function called when a connection is terminated
	// (该Server的连接断开时的Hook函数)
	onConnStop func(conn ziface.IConnection)

	// Data packet encapsulation method
	// (数据报文封包方式)
	packet ziface.IDataPack

	// Asynchronous capture of connection closing status
	// (异步捕获链接关闭状态)
	exitChan chan struct{}

	// Decoder for dealing with message fragmentation and reassembly
	// (断粘包解码器)
	decoder ziface.IDecoder

	// Heartbeat checker
	// (心跳检测器)
	hc ziface.IHeartbeatChecker

	// websocket
	upgrader *websocket.Upgrader

	// websocket connection authentication
	websocketAuth func(r *http.Request) error

	kcpConfig *KcpConfig

	// connection id
	cID uint64
}

type KcpConfig struct {
	// changes ack flush option, set true to flush ack immediately,
	// (改变ack刷新选项，设置为true立即刷新ack)
	KcpACKNoDelay bool
	// toggles the stream mode on/off
	// (切换流模式开/关)
	KcpStreamMode bool
	// Whether nodelay mode is enabled, 0 is not enabled; 1 enabled.
	// (是否启用nodelay模式，0不启用；1启用)
	KcpNoDelay int
	// Protocol internal work interval, in milliseconds, such as 10 ms or 20 ms.
	// (协议内部工作的间隔，单位毫秒，比如10ms或者20ms)
	KcpInterval int
	// Fast retransmission mode, 0 represents off by default, 2 can be set (2 ACK spans will result in direct retransmission)
	// (快速重传模式，默认为0关闭，可以设置2（2次ACK跨越将会直接重传）
	KcpResend int
	// Whether to turn off flow control, 0 represents “Do not turn off” by default, 1 represents “Turn off”.
	// (是否关闭流控，默认是0代表不关闭，1代表关闭)
	KcpNc int
	// SND_BUF, this unit is the packet, default 32.
	// (SND_BUF发送缓冲区大小，单位是包，默认是32)
	KcpSendWindow int
	// RCV_BUF, this unit is the packet, default 32.
	// (RCV_BUF接收缓冲区大小，单位是包，默认是32)
	KcpRecvWindow int
}

// newServerWithConfig creates a server handle based on config
// (根据config创建一个服务器句柄)
func newServerWithConfig(config *zconf.Config, ipVersion string, opts ...Option) ziface.IServer {
	logo.PrintLogo()

	s := &Server{
		Name:             config.Name,
		IPVersion:        ipVersion,
		IP:               config.Host,
		Port:             config.TCPPort,
		WsPort:           config.WsPort,
		KcpPort:          config.KcpPort,
		msgHandler:       newMsgHandle(),
		RouterSlicesMode: config.RouterSlicesMode,
		ConnMgr:          newConnManager(),
		exitChan:         nil,
		// Default to using Zinx's TLV data pack format
		// (默认使用zinx的TLV封包方式)
		packet:  zpack.Factory().NewPack(ziface.ZinxDataPack),
		decoder: zdecoder.NewTLVDecoder(), // Default to using TLV decode (默认使用TLV的解码方式)
		upgrader: &websocket.Upgrader{
			ReadBufferSize: int(config.IOReadBuffSize),
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		kcpConfig: &KcpConfig{
			KcpACKNoDelay: config.KcpACKNoDelay,
			KcpStreamMode: config.KcpStreamMode,
			KcpNoDelay:    config.KcpNoDelay,
			KcpInterval:   config.KcpInterval,
			KcpResend:     config.KcpResend,
			KcpNc:         config.KcpNc,
			KcpSendWindow: config.KcpSendWindow,
			KcpRecvWindow: config.KcpRecvWindow,
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	// Display current configuration information
	// (提示当前配置信息)
	config.Show()

	return s
}

// NewServer creates a server handle
// (创建一个服务器句柄)
func NewServer(opts ...Option) ziface.IServer {
	return newServerWithConfig(zconf.GlobalObject, "tcp", opts...)
}

// NewUserConfServer creates a server handle using user-defined configuration
// (创建一个服务器句柄)
func NewUserConfServer(config *zconf.Config, opts ...Option) ziface.IServer {

	// Refresh user configuration to global configuration variable
	// (刷新用户配置到全局配置变量)
	zconf.UserConfToGlobal(config)

	s := newServerWithConfig(zconf.GlobalObject, "tcp4", opts...)
	return s
}

// NewDefaultRouterSlicesServer creates a server handle with a default RouterRecovery processor.
// (创建一个默认自带一个Recover处理器的服务器句柄)
func NewDefaultRouterSlicesServer(opts ...Option) ziface.IServer {
	zconf.GlobalObject.RouterSlicesMode = true
	s := newServerWithConfig(zconf.GlobalObject, "tcp", opts...)
	s.Use(RouterRecovery)
	return s
}

// NewUserConfDefaultRouterSlicesServer creates a server handle with user-configured options and a default Recover handler.
// If the user does not wish to use the Use method, they should use NewUserConfServer instead.
// (创建一个用户配置的自带一个Recover处理器的服务器句柄，如果用户不希望Use这个方法，那么应该使用NewUserConfServer)
func NewUserConfDefaultRouterSlicesServer(config *zconf.Config, opts ...Option) ziface.IServer {

	if !config.RouterSlicesMode {
		panic("RouterSlicesMode is false")
	}

	// Refresh user configuration to global configuration variable (刷新用户配置到全局配置变量)
	zconf.UserConfToGlobal(config)

	s := newServerWithConfig(zconf.GlobalObject, "tcp4", opts...)
	s.Use(RouterRecovery)
	return s
}

func (s *Server) StartConn(conn ziface.IConnection) {
	// HeartBeat check
	if s.hc != nil {
		// Clone a heart-beat checker from the server side
		heartBeatChecker := s.hc.Clone()

		// Bind current connection
		heartBeatChecker.BindConn(conn)
	}

	// Start processing business for the current connection
	conn.Start()
}

func (s *Server) ListenTcpConn() {
	// 1. Get a TCP address
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		zlog.Ins().ErrorF("[START] resolve tcp addr err: %v\n", err)
		return
	}

	// 2. Listen to the server address
	var listener net.Listener
	if zconf.GlobalObject.CertFile != "" && zconf.GlobalObject.PrivateKeyFile != "" {
		// Read certificate and private key
		crt, err := tls.LoadX509KeyPair(zconf.GlobalObject.CertFile, zconf.GlobalObject.PrivateKeyFile)
		if err != nil {
			panic(err)
		}

		// TLS connection
		tlsConfig := &tls.Config{}
		tlsConfig.Certificates = []tls.Certificate{crt}
		tlsConfig.Time = time.Now
		tlsConfig.Rand = rand.Reader
		listener, err = tls.Listen(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port), tlsConfig)
		if err != nil {
			panic(err)
		}
	} else {
		listener, err = net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}
	}

	// 3. Start server network connection business
	go func() {
		for {
			// 3.1 Set the maximum connection control for the server. If it exceeds the maximum connection, wait.
			// (设置服务器最大连接控制,如果超过最大连接，则等待)
			if s.ConnMgr.Len() >= zconf.GlobalObject.MaxConn {
				zlog.Ins().InfoF("Exceeded the maxConnNum:%d, Wait:%d", zconf.GlobalObject.MaxConn, AcceptDelay.duration)
				AcceptDelay.Delay()
				continue
			}
			// 3.2 Block and wait for a client to establish a connection request.
			// (阻塞等待客户端建立连接请求)
			conn, err := listener.Accept()
			if err != nil {
				//Go 1.17+
				if errors.Is(err, net.ErrClosed) {
					zlog.Ins().ErrorF("Listener closed")
					return
				}
				zlog.Ins().ErrorF("Accept err: %v", err)
				AcceptDelay.Delay()
				continue
			}

			AcceptDelay.Reset()

			// 3.4 Handle the business method for this new connection request. At this time, the handler and conn should be bound.
			// (处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的)
			newCid := atomic.AddUint64(&s.cID, 1)
			dealConn := newServerConn(s, conn, newCid)

			go s.StartConn(dealConn)

		}
	}()
	select {
	case <-s.exitChan:
		err := listener.Close()
		if err != nil {
			zlog.Ins().ErrorF("listener close err: %v", err)
		}
	}
}

func (s *Server) ListenWebsocketConn() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. Check if the server has reached the maximum allowed number of connections
		// (设置服务器最大连接控制,如果超过最大连接，则等待)
		if s.ConnMgr.Len() >= zconf.GlobalObject.MaxConn {
			zlog.Ins().InfoF("Exceeded the maxConnNum:%d, Wait:%d", zconf.GlobalObject.MaxConn, AcceptDelay.duration)
			AcceptDelay.Delay()
			return
		}
		// 2. If websocket authentication is required, set the authentication information
		// (如果需要 websocket 认证请设置认证信息)
		if s.websocketAuth != nil {
			err := s.websocketAuth(r)
			if err != nil {
				zlog.Ins().ErrorF(" websocket auth err:%v", err)
				w.WriteHeader(401)
				AcceptDelay.Delay()
				return
			}
		}
		// 3. Check if there is a subprotocol specified in the header
		// (判断 header 里面是有子协议)
		if len(r.Header.Get("Sec-Websocket-Protocol")) > 0 {
			s.upgrader.Subprotocols = websocket.Subprotocols(r)
		}
		// 4. Upgrade the connection to a websocket connection
		// (升级成 websocket 连接)
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			zlog.Ins().ErrorF("new websocket err:%v", err)
			w.WriteHeader(500)
			AcceptDelay.Delay()
			return
		}
		AcceptDelay.Reset()
		// 5. Handle the business logic of the new connection, which should already be bound to a handler and conn
		// 5. 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
		newCid := atomic.AddUint64(&s.cID, 1)
		wsConn := newWebsocketConn(s, conn, newCid)
		go s.StartConn(wsConn)

	})

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.IP, s.WsPort), nil)
	if err != nil {
		panic(err)
	}
}

func (s *Server) ListenKcpConn() {

	// 1. Listen to the server address
	listener, err := kcp.Listen(fmt.Sprintf("%s:%d", s.IP, s.KcpPort))
	if err != nil {
		zlog.Ins().ErrorF("[START] resolve KCP addr err: %v\n", err)
		return
	}

	zlog.Ins().InfoF("[START] KCP server listening at IP: %s, Port %d, Addr %s", s.IP, s.KcpPort, listener.Addr().String())
	// 2. Start server network connection business
	go func() {
		for {
			// 2.1 Set the maximum connection control for the server. If it exceeds the maximum connection, wait.
			// (设置服务器最大连接控制,如果超过最大连接，则等待)
			if s.ConnMgr.Len() >= zconf.GlobalObject.MaxConn {
				zlog.Ins().InfoF("Exceeded the maxConnNum:%d, Wait:%d", zconf.GlobalObject.MaxConn, AcceptDelay.duration)
				AcceptDelay.Delay()
				continue
			}
			// 2.2 Block and wait for a client to establish a connection request.
			// (阻塞等待客户端建立连接请求)
			conn, err := listener.Accept()
			if err != nil {
				zlog.Ins().ErrorF("Accept KCP err: %v", err)
				AcceptDelay.Delay()
				continue
			}

			AcceptDelay.Reset()

			// 3.4 Handle the business method for this new connection request. At this time, the handler and conn should be bound.
			// (处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn 是绑定的)
			newCid := atomic.AddUint64(&s.cID, 1)

			kcpConn := conn.(*kcp.UDPSession)
			kcpConn.SetACKNoDelay(s.kcpConfig.KcpACKNoDelay)
			kcpConn.SetStreamMode(s.kcpConfig.KcpStreamMode)
			kcpConn.SetNoDelay(s.kcpConfig.KcpNoDelay, s.kcpConfig.KcpInterval, s.kcpConfig.KcpResend, s.kcpConfig.KcpNc)
			kcpConn.SetWindowSize(s.kcpConfig.KcpSendWindow, s.kcpConfig.KcpRecvWindow)

			dealConn := newKcpServerConn(s, kcpConn, newCid)

			go s.StartConn(dealConn)
		}
	}()
	select {
	case <-s.exitChan:
		err := listener.Close()
		if err != nil {
			zlog.Ins().ErrorF("KCP listener close err: %v", err)
		}
	}
}

// Start the network service
// (开启网络服务)
func (s *Server) Start() {
	zlog.Ins().InfoF("[START] Server name: %s,listener at IP: %s, Port %d is starting", s.Name, s.IP, s.Port)
	s.exitChan = make(chan struct{})

	// Add decoder to interceptors
	// (将解码器添加到拦截器)
	if s.decoder != nil {
		s.msgHandler.AddInterceptor(s.decoder)
	}
	// Start worker pool mechanism
	// (启动worker工作池机制)
	s.msgHandler.StartWorkerPool()

	// Start a goroutine to handle server listener business
	// (开启一个go去做服务端Listener业务)
	switch zconf.GlobalObject.Mode {
	case zconf.ServerModeTcp:
		go s.ListenTcpConn()
	case zconf.ServerModeWebsocket:
		go s.ListenWebsocketConn()
	case zconf.ServerModeKcp:
		go s.ListenKcpConn()
	default:
		go s.ListenTcpConn()
		go s.ListenWebsocketConn()
	}

}

// Stop stops the server (停止服务)
func (s *Server) Stop() {
	zlog.Ins().InfoF("[STOP] Zinx server , name %s", s.Name)

	// Clear other connection information or other information that needs to be cleaned up
	// (将其他需要清理的连接信息或者其他信息 也要一并停止或者清理)
	s.ConnMgr.ClearConn()
	s.exitChan <- struct{}{}
	close(s.exitChan)
}

// Serve runs the server (运行服务)
func (s *Server) Serve() {
	s.Start()
	// Block, otherwise the listener's goroutine will exit when the main Go exits (阻塞,否则主Go退出， listenner的go将会退出)
	c := make(chan os.Signal, 1)
	// Listen for specified signals: ctrl+c or kill signal (监听指定信号 ctrl+c kill信号)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	zlog.Ins().InfoF("[SERVE] Zinx server , name %s, Serve Interrupt, signal = %v", s.Name, sig)
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	if s.RouterSlicesMode {
		panic("Server RouterSlicesMode is true ")
	}
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) AddRouterSlices(msgID uint32, router ...ziface.RouterHandler) ziface.IRouterSlices {
	if !s.RouterSlicesMode {
		panic("Server RouterSlicesMode is false ")
	}
	return s.msgHandler.AddRouterSlices(msgID, router...)
}

func (s *Server) Group(start, end uint32, Handlers ...ziface.RouterHandler) ziface.IGroupRouterSlices {
	if !s.RouterSlicesMode {
		panic("Server RouterSlicesMode is false")
	}
	return s.msgHandler.Group(start, end, Handlers...)
}

func (s *Server) Use(Handlers ...ziface.RouterHandler) ziface.IRouterSlices {
	if !s.RouterSlicesMode {
		panic("Server RouterSlicesMode is false")
	}
	return s.msgHandler.Use(Handlers...)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.onConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.onConnStop = hookFunc
}

func (s *Server) GetOnConnStart() func(ziface.IConnection) {
	return s.onConnStart
}

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

// StartHeartBeat starts the heartbeat check.
// interval is the time interval between each heartbeat.
// (启动心跳检测
// interval 每次发送心跳的时间间隔)
func (s *Server) StartHeartBeat(interval time.Duration) {
	checker := NewHeartbeatChecker(interval)

	// Add the heartbeat check router. (添加心跳检测的路由)
	//检测当前路由模式
	if s.RouterSlicesMode {
		s.AddRouterSlices(checker.MsgID(), checker.RouterSlices()...)
	} else {
		s.AddRouter(checker.MsgID(), checker.Router())
	}

	// Bind the heartbeat checker to the server. (server绑定心跳检测器)
	s.hc = checker
}

// StartHeartBeatWithOption starts the heartbeat detection with the given configuration.
// interval is the time interval for sending heartbeat messages.
// option is the configuration for heartbeat detection.
// 启动心跳检测
// (option 心跳检测的配置)
func (s *Server) StartHeartBeatWithOption(interval time.Duration, option *ziface.HeartBeatOption) {
	checker := NewHeartbeatChecker(interval)

	// Configure the heartbeat checker with the provided options
	if option != nil {
		checker.SetHeartbeatMsgFunc(option.MakeMsg)
		checker.SetOnRemoteNotAlive(option.OnRemoteNotAlive)
		//检测当前路由模式
		if s.RouterSlicesMode {
			checker.BindRouterSlices(option.HeartBeatMsgID, option.RouterSlices...)
		} else {
			checker.BindRouter(option.HeartBeatMsgID, option.Router)
		}
	}

	// Add the heartbeat checker's router to the server's router (添加心跳检测的路由)
	//检测当前路由模式
	if s.RouterSlicesMode {
		s.AddRouterSlices(checker.MsgID(), checker.RouterSlices()...)
	} else {
		s.AddRouter(checker.MsgID(), checker.Router())
	}

	// Bind the server with the heartbeat checker (server绑定心跳检测器)
	s.hc = checker
}

func (s *Server) GetHeartBeat() ziface.IHeartbeatChecker {
	return s.hc
}

func (s *Server) SetDecoder(decoder ziface.IDecoder) {
	s.decoder = decoder
}

func (s *Server) GetLengthField() *ziface.LengthField {
	if s.decoder != nil {
		return s.decoder.GetLengthField()
	}
	return nil
}

func (s *Server) AddInterceptor(interceptor ziface.IInterceptor) {
	s.msgHandler.AddInterceptor(interceptor)
}

func (s *Server) SetWebsocketAuth(f func(r *http.Request) error) {
	s.websocketAuth = f
}

func (s *Server) ServerName() string {
	return s.Name
}

func init() {}
