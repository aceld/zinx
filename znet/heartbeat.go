package znet

import (
	"time"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
)

type HeartbeatChecker struct {
	interval time.Duration //  Heartbeat detection interval(心跳检测时间间隔)
	quitChan chan bool     // Quit signal(退出信号)

	onRemoteNotAlive ziface.OnRemoteNotAlive //  User-defined method for handling remote connections that are not alive (用户自定义的远程连接不存活时的处理方法)

	conn ziface.IConnection // Bound connection(绑定的链接)
}

func notAliveDefaultFunc(conn ziface.IConnection) {
	zlog.Ins().InfoF("Remote connection %s is not alive, stop it", conn.RemoteAddr())
	conn.Stop()
}

func NewHeartbeatChecker(interval time.Duration) ziface.IHeartbeatChecker {
	heartbeat := &HeartbeatChecker{
		interval: interval,
		quitChan: make(chan bool),
		// Use default heartbeat message generation function and remote connection not alive handling method
		// (均使用默认的心跳消息生成函数和远程连接不存活时的处理方法)
		onRemoteNotAlive: notAliveDefaultFunc,
	}

	return heartbeat
}

func (h *HeartbeatChecker) SetOnRemoteNotAlive(f ziface.OnRemoteNotAlive) {
	if f != nil {
		h.onRemoteNotAlive = f
	}
}

func (h *HeartbeatChecker) start() {
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

func (h *HeartbeatChecker) Start() {
	go h.start()
}

func (h *HeartbeatChecker) Stop() {
	zlog.Ins().InfoF("heartbeat checker stop, connID=%+v", h.conn.GetConnID())
	h.quitChan <- true
}

func (h *HeartbeatChecker) check() {
	if h.conn == nil {
		return
	}

	if !h.conn.IsAlive() {
		h.onRemoteNotAlive(h.conn)
	}
}

func (h *HeartbeatChecker) BindConn(conn ziface.IConnection) {
	h.conn = conn
	conn.SetHeartBeat(h)
}

// Clone clones to a specified connection
// (克隆到一个指定的链接上)
func (h *HeartbeatChecker) Clone() ziface.IHeartbeatChecker {
	heartbeat := &HeartbeatChecker{
		interval:         h.interval,
		quitChan:         make(chan bool),
		onRemoteNotAlive: h.onRemoteNotAlive,
		conn:             nil, // The bound connection needs to be reassigned
	}
	return heartbeat
}
