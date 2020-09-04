package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync/atomic"
	"time"
	"wsserver/configs"
	"wsserver/iserverface"
)

type Server struct {
	//服务器名称
	Name string
	//IPversion
	IPversion string
	//服务器IP地址
	IP string
	//端口号
	Port int
	//Server的消息管理模块
	MsgHandler iserverface.IMsgHandle
	//当前Server链接管理器
	ConnMgr iserverface.IConnMgr
	//当前Server连接创建时的hook函数
	OnConnStart func(conn iserverface.IConnection)
	//当前Server连接断开时的hook函数
	OnConnStop func(conn iserverface.IConnection)
}

var (
	GWServer iserverface.IServer
)

func NewServer() iserverface.IServer {
	return &Server{
		Name:        configs.GConf.Name,
		IPversion:   configs.GConf.IpVersion,
		IP:          configs.GConf.Ip,
		Port:        configs.GConf.Port,
		ConnMgr:     NewConnManager(),
		MsgHandler:  NewMsgHandle(),
	}
}

func (s *Server) Start(c *gin.Context) {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	//开启一个go去做服务端Linster业务
	go func() {


		//TODO server.go 应该有一个自动生成ID的方法
		curConnId := uint64(time.Now().Unix())
		connId := atomic.AddUint64(&curConnId, 1)
		//3.1 阻塞等待客户端建立连接请求
		var (
			err      error
			wsSocket *websocket.Conn
			wsUpgrader = websocket.Upgrader{
				// 允许所有CORS跨域请求
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
		)

		if wsSocket, err = wsUpgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
			return
		}
		fmt.Println("Get conn remote addr = ", wsSocket.RemoteAddr().String())
		//3 启动server网络连接业务


		//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
		/*
		if s.ConnMgr.Len() >= configs.GConf.MaxConn {
			wsSocket.Close()
			continue
		}
		**/
		//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
		dealConn := NewConnection(s, wsSocket, connId, s.MsgHandler)

		fmt.Println("Current connId:",connId)
		//3.4 启动当前链接的处理业务
		go dealConn.Start()

	}()
}

//停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Websocket server , name ", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

//运行服务
func (s *Server) Serve(c *gin.Context) {
	s.Start(c)

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgId string, router iserverface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

//得到链接管理
func (s *Server) GetConnMgr() iserverface.IConnMgr {
	return s.ConnMgr
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(iserverface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(iserverface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn iserverface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn iserverface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

