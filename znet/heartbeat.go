package znet

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"time"
)

type HeartbeatChecker struct {
	interval time.Duration // 心跳检测时间间隔
	quitChan chan bool     // 退出信号

	makeMsg ziface.HeartBeatMsgFunc //用户自定义的心跳检测消息处理方法

	onRemoteNotAlive ziface.OnRemoteNotAlive //用户自定义的远程连接不存活时的处理方法

	msgID  uint32         // 心跳的消息ID
	router ziface.IRouter //用户自定义的心跳检测消息业务处理路由

	conn ziface.IConnection // 绑定的链接
}

/*
	收到remote心跳消息的默认回调路由业务
*/
type HeatBeatDefaultRouter struct {
	BaseRouter
}

//Handle -
func (r *HeatBeatDefaultRouter) Handle(req ziface.IRequest) {
	zlog.Ins().InfoF("Recv Heartbeat from %s, MsgID = %+v, Data = %s", req.GetConnection().RemoteAddr(), req.GetMsgID(), string(req.GetData()))
}

//默认的心跳消息生成函数
func makeDefaultMsg(conn ziface.IConnection) []byte {
	msg := fmt.Sprintf("heartbeat [%s->%s]", conn.LocalAddr(), conn.RemoteAddr())
	return []byte(msg)
}

//默认的心跳检测函数
func notAliveDefaultFunc(conn ziface.IConnection) {
	zlog.Ins().InfoF("Remote connection %s is not alive, stop it", conn.RemoteAddr())
	conn.Stop()
}

func NewHeartbeatChecker(interval time.Duration, conn ziface.IConnection) ziface.IHeartbeatChecker {
	heatbeat := &HeartbeatChecker{
		interval: interval,
		conn:     conn,
		quitChan: make(chan bool),

		//均使用默认的心跳消息生成函数和远程连接不存活时的处理方法
		makeMsg:          makeDefaultMsg,
		onRemoteNotAlive: notAliveDefaultFunc,
		msgID:            ziface.HeartBeatDefaultMsgID,
		router:           &HeatBeatDefaultRouter{},
	}

	return heatbeat
}

// NewHeartbeatCheckerS Server创建心跳检测器
func NewHeartbeatCheckerS(interval time.Duration, server ziface.IServer) ziface.IHeartbeatChecker {
	heatbeat := &HeartbeatChecker{
		interval: interval,
		s:        server,
		c:        nil,
		quitChan: make(chan bool),

		//均使用默认的心跳消息生成函数和远程连接不存活时的处理方法
		makeMsg:          makeDefaultMsg,
		onRemoteNotAlive: notAliveDefaultFunc,
		msgID:            ziface.HeartBeatDefaultMsgID,
		router:           &HeatBeatDefaultRouter{},
	}

	return heatbeat
}

// NewHeartbeatCheckerC Client创建心跳检测器
func NewHeartbeatCheckerC(interval time.Duration, client ziface.IClient) ziface.IHeartbeatChecker {
	heatbeat := &HeartbeatChecker{
		interval: interval,
		c:        client,
		s:        nil,
		quitChan: make(chan bool),

		//均使用默认的心跳消息生成函数和远程连接不存活时的处理方法
		makeMsg:          makeDefaultMsg,
		onRemoteNotAlive: notAliveDefaultFunc,
		msgID:            ziface.HeartBeatDefaultMsgID,
		router:           &HeatBeatDefaultRouter{},
	}

	return heatbeat
}

func (h *HeartbeatChecker) SetOnRemoteNotAlive(f ziface.OnRemoteNotAlive) {
	if f != nil {
		h.onRemoteNotAlive = f
	}
}

func (h *HeartbeatChecker) SetHeartbeatMsgFunc(f ziface.HeartBeatMsgFunc) {
	if f != nil {
		h.makeMsg = f
	}
}

func (h *HeartbeatChecker) BindRouter(msgID uint32, router ziface.IRouter) {
	if router != nil && msgID != ziface.HeartBeatDefaultMsgID {
		h.msgID = msgID
		h.router = router
	}
}

// 启动心跳检测
func (h *HeartbeatChecker) Start() {
	ticker := time.NewTicker(h.interval)
	for {
		select {
		case <-ticker.C:
			h.check()
		case <-h.quitChan:
			ticker.Stop()
			return
		}
	}
}

// 停止心跳检测
func (h *HeartbeatChecker) Stop() {
	h.quitChan <- true
}

func (h *HeartbeatChecker) checkServer() {
	if h.s.GetConnMgr() != nil {
		// server
		_ = h.s.GetConnMgr().Range(doCheck, h)
	}
}

func (h *HeartbeatChecker) checkClient() {
	if h.c.Conn() != nil {
		//client
		_ = doCheck(0, h.c.Conn(), h)
	}
}

// 检测单个连接 和 发送心跳包
func doCheck(connID uint64, conn ziface.IConnection, args interface{}) (err error) {
	hc := args.(*HeartbeatChecker)

	if !conn.IsAlive() {
		hc.onRemoteNotAlive(conn)
	} else {
		err = hc.SendHeartBeatMsg()
	}

	return err
}

func (h *HeartbeatChecker) SendHeartBeatMsg() error {

	msg := h.makeMsg(h.c.Conn())

	err := h.c.Conn().SendMsg(h.msgID, msg)
	if err != nil {
		zlog.Ins().ErrorF("send heartbeat msg error: %v, msgId=%+v msg=%+v", err, h.msgID, msg)
		return err
	}

	return nil
}

// 执行心跳检测
func (h *HeartbeatChecker) check() {
	if h.s != nil {
		h.checkServer()
	} else if h.c != nil {
		h.checkClient()
	}
}

//深拷贝
func (h *HeartbeatChecker) Clone() ziface.IHeartbeatChecker {

	heatbeat := &HeartbeatChecker{
		interval:         h.interval,
		c:                h.c,
		s:                h.s,
		quitChan:         make(chan bool),
		makeMsg:          h.makeMsg,
		onRemoteNotAlive: h.onRemoteNotAlive,
		msgID:            h.msgID,
		router:           h.router,
	}

	return heatbeat
}

func (h *HeartbeatChecker) MsgID() uint32 {
	return h.msgID
}

func (h *HeartbeatChecker) Router() ziface.IRouter {
	return h.router
}
