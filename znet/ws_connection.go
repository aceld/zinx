package znet

import (
	"context"
	"encoding/hex"
	"errors"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zinterceptor"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/zpack"
	"github.com/gorilla/websocket"
)

// WsConnection is a module for handling the read and write operations of a WebSocket connection.
// (Websocket连接模块, 用于处理 Websocket 连接的读写业务 一个连接对应一个Connection)
type WsConnection struct {
	// conn is the current connection's WebSocket socket TCP socket. (当前连接的socket TCP套接字)
	conn *websocket.Conn

	// connID is the current connection's ID, which can also be referred to as SessionID, and is globally unique.
	// uint64 range: 018,446,744,073,709,551,615
	// This is the maximum number of connIDs that the theory supports per process.
	// (当前连接的ID 也可以称作为SessionID，ID全局唯一 ，服务端Connection使用
	//  uint64 取值范围：0 ~ 18,446,744,073,709,551,615
	//  这个是理论支持的进程connID的最大数量)
	connID uint64

	// connection id for string
	// (字符串的连接id)
	connIdStr string

	// The workerid responsible for handling the link
	// 负责处理该链接的workerid
	workerID uint32

	// msgHandler is the message management module for MsgID and the corresponding message handling method.
	// (消息管理MsgID和对应处理方法的消息管理模块)
	msgHandler ziface.IMsgHandle

	// ctx and cancel are used to notify that the connection has exited/stopped.
	// (告知该链接已经退出/停止的channel)
	ctx    context.Context
	cancel context.CancelFunc

	// msgBuffChan is a buffered channel used for message communication between the read and write goroutines.
	// (有缓冲管道，用于读、写两个goroutine之间的消息通信)
	msgBuffChan chan []byte

	// msgLock is used for locking when users send and receive messages.
	// (用户收发消息的Lock)
	msgLock sync.RWMutex

	// property is the connection attribute. (链接属性)
	property map[string]interface{}

	// propertyLock protects the current property lock. (保护当前property的锁)
	propertyLock sync.Mutex

	// isClosed is the current connection's closed state. (当前连接的关闭状态)
	isClosed bool

	// connManager is the Connection Manager to which the current connection belongs. (当前链接是属于哪个Connection Manager的)
	connManager ziface.IConnManager

	// onConnStart is the Hook function when the current connection is created.
	// (当前连接创建时Hook函数)
	onConnStart func(conn ziface.IConnection)

	// onConnStop is the Hook function when the current connection is disconnected.
	// (当前连接断开时的Hook函数)
	onConnStop func(conn ziface.IConnection)

	// packet is the data packet format.
	// (数据报文封包方式)
	packet ziface.IDataPack

	// lastActivityTime is the last time the connection was active.
	// (最后一次活动时间)
	lastActivityTime time.Time

	// frameDecoder is the decoder for splitting or splicing data packets.
	// (断粘包解码器)
	frameDecoder ziface.IFrameDecoder

	// hc is the Heartbeat Checker. (心跳检测器)
	hc ziface.IHeartbeatChecker

	// name is the name of the connection and is the same as the Name of the Server/Client that created the connection.
	// (链接名称，默认与创建链接的Server/Client的Name一致)
	name string

	// localAddr is the local address of the current connection. (当前链接的本地地址)
	localAddr string

	// remoteAddr is the remote address of the current connection. (当前链接的远程地址)
	remoteAddr string

	// Close callback
	closeCallback callbacks

	// Close callback mutex
	closeCallbackMutex sync.RWMutex
}

// newServerConn: for Server, a method to create a connection with Server characteristics
// Note: The name has been changed from NewConnection
// (newServerConn :for Server, 创建一个Server服务端特性的连接的方法
// Note: 名字由 NewConnection 更变)
func newWebsocketConn(server ziface.IServer, conn *websocket.Conn, connID uint64) ziface.IConnection {
	// Initialize Conn properties (初始化Conn属性)
	c := &WsConnection{
		conn:        conn,
		connID:      connID,
		connIdStr:   strconv.FormatUint(connID, 10),
		isClosed:    false,
		msgBuffChan: nil,
		property:    nil,
		name:        server.ServerName(),
		localAddr:   conn.LocalAddr().String(),
		remoteAddr:  conn.RemoteAddr().String(),
	}

	lengthField := server.GetLengthField()
	if lengthField != nil {
		c.frameDecoder = zinterceptor.NewFrameDecoder(*lengthField)
	}

	// Inherited attributes from server (从server继承过来的属性)
	c.packet = server.GetPacket()
	c.onConnStart = server.GetOnConnStart()
	c.onConnStop = server.GetOnConnStop()
	c.msgHandler = server.GetMsgHandler()

	// Bind the current Connection to the Server's ConnManager (将当前的Connection与Server的ConnManager绑定)
	c.connManager = server.GetConnMgr()

	// Add the newly created Conn to the connection management (将新创建的Conn添加到链接管理中)
	server.GetConnMgr().Add(c)

	return c
}

// newClientConn :for Client, creates a connection with Client-side features
// (newClientConn :for Client, 创建一个Client服务端特性的连接的方法)
func newWsClientConn(client ziface.IClient, conn *websocket.Conn) ziface.IConnection {
	c := &WsConnection{
		conn:        conn,
		connID:      0,  // client ignore
		connIdStr:   "", // client ignore
		isClosed:    false,
		msgBuffChan: nil,
		property:    nil,
		name:        client.GetName(),
		localAddr:   conn.LocalAddr().String(),
		remoteAddr:  conn.RemoteAddr().String(),
	}

	lengthField := client.GetLengthField()
	if lengthField != nil {
		c.frameDecoder = zinterceptor.NewFrameDecoder(*lengthField)
	}

	// Inherit properties from client (从client继承过来的属性)
	c.packet = client.GetPacket()
	c.onConnStart = client.GetOnConnStart()
	c.onConnStop = client.GetOnConnStop()
	c.msgHandler = client.GetMsgHandler()

	return c
}

// StartWriter is a Goroutine that sends messages to the client
// (StartWriter 写消息Goroutine， 用户将数据发送给客户端)
func (c *WsConnection) StartWriter() {
	zlog.Ins().InfoF("Writer Goroutine is running")
	defer zlog.Ins().InfoF("%s [conn Writer exit!]", c.RemoteAddr().String())

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				if err := c.Send(data); err != nil {
					zlog.Ins().ErrorF("Send Buff Data error:, %s Conn Writer exit", err)
					break
				}

			} else {
				zlog.Ins().ErrorF("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// StartReader is a Goroutine that reads messages from the client.
// (StartReader 读消息Goroutine，用于从客户端中读取数据)
func (c *WsConnection) StartReader() {
	zlog.Ins().InfoF("[Reader Goroutine is running]")
	defer zlog.Ins().InfoF("%s [conn Reader exit!]", c.RemoteAddr().String())
	defer c.Stop()

	// Create a pack-unpack object. (创建拆包解包的对象)
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// add by uuxia 2023-02-03
			// Read data from the conn's IO to the memory buffer.
			// (从conn的IO中读取数据到内存缓冲buffer中)
			messageType, buffer, err := c.conn.ReadMessage()
			if err != nil {
				c.cancel()
				return
			}
			if messageType == websocket.PingMessage {
				c.updateActivity()
				continue
			}
			n := len(buffer)
			if err != nil {
				zlog.Ins().ErrorF("read msg head [read datalen=%d], error = %s", n, err.Error())
				return
			}
			zlog.Ins().DebugF("read buffer %s \n", hex.EncodeToString(buffer[0:n]))

			// Update the Active status of heartbeat detection normally after reading data from the peer.
			// (正常读取到对端数据，更新心跳检测Active状态)
			if n > 0 && c.hc != nil {
				c.updateActivity()
			}

			// Handle custom protocol fragmentation and packet sticking issues add by uuxia 2023-03-21
			// (处理自定义协议断粘包问题)
			if c.frameDecoder != nil {
				// Decode the 0-n bytes of data read.
				// (为读取到的0-n个字节的数据进行解码)
				bufArrays := c.frameDecoder.Decode(buffer)
				if bufArrays == nil {
					continue
				}
				for _, bytes := range bufArrays {
					zlog.Ins().DebugF("read buffer %s \n", hex.EncodeToString(bytes))
					msg := zpack.NewMessage(uint32(len(bytes)), bytes)
					// Get the Request data requested by the current client.
					// (得到当前客户端请求的Request数据)
					req := GetRequest(c, msg)
					c.msgHandler.Execute(req)
				}
			} else {
				msg := zpack.NewMessage(uint32(n), buffer[0:n])
				// Get the Request data requested by the current client.
				// (得到当前客户端请求的Request数据)
				req := GetRequest(c, msg)
				c.msgHandler.Execute(req)
			}
		}
	}
}

// Start starts the connection and makes it work.
// (Start 启动连接，让当前连接开始工作)
func (c *WsConnection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	// Execute the hook method according to the business needs of creating the connection passed in by the user.
	// (按照用户传递进来的创建连接时需要处理的业务，执行钩子方法)
	c.callOnConnStart()

	// Start the heartbeat check
	// (启动心跳检测)
	if c.hc != nil {
		c.hc.Start()
		c.updateActivity()
	}

	// 占用workerid
	c.workerID = useWorker(c)

	// Start the Goroutine for users to read data from the client.
	// (开启用户从客户端读取数据流程的Goroutine)
	go c.StartReader()

	select {
	case <-c.ctx.Done():
		c.finalizer()

		// 归还workerid
		freeWorker(c)
		return
	}
}

// Stop stops the connection and ends its current state.
// (停止连接，结束当前连接状态)
func (c *WsConnection) Stop() {
	c.cancel()
}

func (c *WsConnection) GetConnection() net.Conn {
	return nil
}

func (c *WsConnection) GetWsConn() *websocket.Conn {
	return c.conn
}

// Deprecated: use GetConnection instead
func (c *WsConnection) GetTCPConnection() net.Conn {
	return nil
}

func (c *WsConnection) GetConnID() uint64 {
	return c.connID
}

func (c *WsConnection) GetConnIdStr() string {
	return c.connIdStr
}

func (c *WsConnection) GetWorkerID() uint32 {
	return c.workerID
}

func (c *WsConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *WsConnection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *WsConnection) Send(data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()
	if c.isClosed == true {
		return errors.New("WsConnection closed when send msg")
	}

	err := c.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		zlog.Ins().ErrorF("SendMsg err data = %+v, err = %+v", data, err)
		return err
	}

	return nil
}

func (c *WsConnection) SendToQueue(data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()

	if c.msgBuffChan == nil {
		c.msgBuffChan = make(chan []byte, zconf.GlobalObject.MaxMsgChanLen)
		// Start a goroutine for writing data back to the client,
		// which only reads data from MsgBuffChan and hasn't allocated memory or started the coroutine until SendBuffMsg is called
		// (开启用于写回客户端数据流程的Goroutine
		// 此方法只读取MsgBuffChan中的数据没调用SendBuffMsg可以分配内存和启用协程)
		go c.StartWriter()
	}

	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed == true {
		return errors.New("WsConnection closed when send buff msg")
	}

	if data == nil {
		zlog.Ins().ErrorF("Pack data is nil")
		return errors.New("Pack data is nil ")
	}

	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- data:
		return nil
	}
}

// SendMsg directly sends the Message data to the remote TCP client.
// (直接将Message数据发送数据给远程的TCP客户端)
func (c *WsConnection) SendMsg(msgID uint32, data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()
	if c.isClosed == true {
		return errors.New("WsConnection closed when send msg")
	}

	// Package data and send
	// (将data封包，并且发送)
	msg, err := c.packet.Pack(zpack.NewMsgPackage(msgID, data))
	if err != nil {
		zlog.Ins().ErrorF("Pack error msg ID = %d", msgID)
		return errors.New("Pack error msg ")
	}

	// Write back to the client
	err = c.conn.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		zlog.Ins().ErrorF("SendMsg err msg ID = %d, data = %+v, err = %+v", msgID, string(msg), err)
		return err
	}

	return nil
}

// SendBuffMsg sends BuffMsg
func (c *WsConnection) SendBuffMsg(msgID uint32, data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()

	if c.msgBuffChan == nil {
		c.msgBuffChan = make(chan []byte, zconf.GlobalObject.MaxMsgChanLen)
		// Start the Goroutine for writing back to the client data stream
		// This method only reads data from MsgBuffChan, allocating memory and starting Goroutine without calling SendBuffMsg
		// (开启用于写回客户端数据流程的Goroutine
		// 此方法只读取MsgBuffChan中的数据没调用SendBuffMsg可以分配内存和启用协程)
		go c.StartWriter()
	}

	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed == true {
		return errors.New("WsConnection closed when send buff msg")
	}

	// Package data and send
	// (将data封包，并且发送)
	msg, err := c.packet.Pack(zpack.NewMsgPackage(msgID, data))
	if err != nil {
		zlog.Ins().ErrorF("Pack error msg ID = %d", msgID)
		return errors.New("Pack error msg ")
	}

	// Send timeout
	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}
}

func (c *WsConnection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

func (c *WsConnection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

func (c *WsConnection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// Context returns the context for the connection, which can be used by user-defined goroutines to get the connection exit status.
// (返回ctx，用于用户自定义的go程获取连接退出状态)
func (c *WsConnection) Context() context.Context {
	return c.ctx
}

func (c *WsConnection) finalizer() {
	// If the user has registered a close callback for the connection, it should be called explicitly at this moment.
	// (如果用户注册了该链接的	关闭回调业务，那么在此刻应该显示调用)
	c.callOnConnStop()

	c.msgLock.Lock()
	defer c.msgLock.Unlock()

	// If the current connection is already closed.
	// (如果当前链接已经关闭)
	if c.isClosed == true {
		return
	}

	// Stop the heartbeat detector bound to the connection.
	// (关闭链接绑定的心跳检测器)
	if c.hc != nil {
		c.hc.Stop()
	}

	// Close the socket connection.
	// (关闭socket链接)
	_ = c.conn.Close()

	// Remove the connection from the connection manager.
	// (将链接从连接管理器中删除)
	if c.connManager != nil {
		c.connManager.Remove(c)
	}

	// Close all channels associated with this connection.
	// (关闭该链接全部管道)
	if c.msgBuffChan != nil {
		close(c.msgBuffChan)
	}

	// Set the flag to indicate that the connection is closed. (设置标志位)
	c.isClosed = true

	go func() {
		defer func() {
			if err := recover(); err != nil {
				zlog.Ins().ErrorF("Conn finalizer panic: %v", err)
			}
		}()

		c.InvokeCloseCallbacks()
	}()

	zlog.Ins().InfoF("Conn Stop()...ConnID = %d", c.connID)
}

func (c *WsConnection) callOnConnStart() {
	if c.onConnStart != nil {
		zlog.Ins().InfoF("ZINX CallOnConnStart....")
		c.onConnStart(c)
	}
}

func (c *WsConnection) callOnConnStop() {
	if c.onConnStop != nil {
		zlog.Ins().InfoF("ZINX CallOnConnStop....")
		c.onConnStop(c)
	}
}

func (c *WsConnection) IsAlive() bool {
	if c.isClosed {
		return false
	}
	// Check the time duration since the last activity of the connection, if it exceeds the maximum heartbeat interval,
	// then the connection is considered dead
	// (检查连接最后一次活动时间，如果超过心跳间隔，则认为连接已经死亡)
	return time.Now().Sub(c.lastActivityTime) < zconf.GlobalObject.HeartbeatMaxDuration()
}

func (c *WsConnection) updateActivity() {
	c.lastActivityTime = time.Now()
}

func (c *WsConnection) SetHeartBeat(checker ziface.IHeartbeatChecker) {
	c.hc = checker
}

func (c *WsConnection) LocalAddrString() string {
	return c.localAddr
}

func (c *WsConnection) RemoteAddrString() string {
	return c.remoteAddr
}

func (c *WsConnection) GetName() string {
	return c.name
}

func (c *WsConnection) GetMsgHandler() ziface.IMsgHandle {
	return c.msgHandler
}

func (s *WsConnection) AddCloseCallback(handler, key interface{}, f func()) {
	if s.isClosed {
		return
	}
	s.closeCallbackMutex.Lock()
	defer s.closeCallbackMutex.Unlock()
	s.closeCallback.Add(handler, key, f)
}

func (s *WsConnection) RemoveCloseCallback(handler, key interface{}) {
	if s.isClosed {
		return
	}
	s.closeCallbackMutex.Lock()
	defer s.closeCallbackMutex.Unlock()
	s.closeCallback.Remove(handler, key)
}

func (s *WsConnection) InvokeCloseCallbacks() {
	s.closeCallbackMutex.RLock()
	defer s.closeCallbackMutex.RUnlock()
	s.closeCallback.Invoke()
}
