package znet

import (
	"context"
	"encoding/hex"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/zinterceptor"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/zpack"
	"github.com/gorilla/websocket"

	"github.com/aceld/zinx/ziface"
)

// Connection TCP connection module
// Used to handle the read and write business of TCP connections, one Connection corresponds to one connection
// (用于处理Tcp连接的读写业务 一个连接对应一个Connection)
type Connection struct {
	// // The socket TCP socket of the current connection(当前连接的socket TCP套接字)
	conn net.Conn

	// The ID of the current connection, also known as SessionID, globally unique, used by server Connection
	// uint64 range: 0~18,446,744,073,709,551,615
	// This is the maximum number of connID theoretically supported by the process
	// (当前连接的ID 也可以称作为SessionID，ID全局唯一 ，服务端Connection使用
	// uint64 取值范围：0 ~ 18,446,744,073,709,551,615
	// 这个是理论支持的进程connID的最大数量)
	connID uint64

	// The workerid responsible for handling the link
	// 负责处理该链接的workerid
	workerID uint32

	// The message management module that manages MsgID and the corresponding processing method
	// (消息管理MsgID和对应处理方法的消息管理模块)
	msgHandler ziface.IMsgHandle

	// Channel to notify that the connection has exited/stopped
	// (告知该链接已经退出/停止的channel)
	ctx    context.Context
	cancel context.CancelFunc

	// Buffered channel used for message communication between the read and write goroutines
	// (有缓冲管道，用于读、写两个goroutine之间的消息通信)
	msgBuffChan chan []byte

	// Lock for user message reception and transmission
	// (用户收发消息的Lock)
	msgLock sync.RWMutex

	// Connection properties
	// (链接属性)
	property map[string]interface{}

	// Lock to protect the current property
	// (保护当前property的锁)
	propertyLock sync.Mutex

	// The current connection's close state
	// (当前连接的关闭状态)
	isClosed bool

	// Which Connection Manager the current connection belongs to
	// (当前链接是属于哪个Connection Manager的)
	connManager ziface.IConnManager

	// Hook function when the current connection is created
	// (当前连接创建时Hook函数)
	onConnStart func(conn ziface.IConnection)

	// Hook function when the current connection is disconnected
	// (当前连接断开时的Hook函数)
	onConnStop func(conn ziface.IConnection)

	// Data packet packaging method
	// (数据报文封包方式)
	packet ziface.IDataPack

	// Last activity time
	// (最后一次活动时间)
	lastActivityTime time.Time

	// Framedecoder for solving fragmentation and packet sticking problems
	// (断粘包解码器)
	frameDecoder ziface.IFrameDecoder

	// Heartbeat checker
	// (心跳检测器)
	hc ziface.IHeartbeatChecker

	// Connection name, default to be the same as the name of the Server/Client that created the connection
	// (链接名称，默认与创建链接的Server/Client的Name一致)
	name string

	// Local address of the current connection
	// (当前链接的本地地址)
	localAddr string

	// Remote address of the current connection
	// (当前链接的远程地址)
	remoteAddr string
}

// newServerConn :for Server, method to create a Server-side connection with Server-specific properties
// (创建一个Server服务端特性的连接的方法)
func newServerConn(server ziface.IServer, conn net.Conn, connID uint64) ziface.IConnection {

	// Initialize Conn properties
	c := &Connection{
		conn:        conn,
		connID:      connID,
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

	// Inherited properties from server (从server继承过来的属性)
	c.packet = server.GetPacket()
	c.onConnStart = server.GetOnConnStart()
	c.onConnStop = server.GetOnConnStop()
	c.msgHandler = server.GetMsgHandler()

	// Bind the current Connection with the Server's ConnManager
	// (将当前的Connection与Server的ConnManager绑定)
	c.connManager = server.GetConnMgr()

	// Add the newly created Conn to the connection manager
	// (将新创建的Conn添加到链接管理中)
	server.GetConnMgr().Add(c)

	return c
}

// newServerConn :for Server, method to create a Client-side connection with Client-specific properties
// (创建一个Client服务端特性的连接的方法)
func newClientConn(client ziface.IClient, conn net.Conn) ziface.IConnection {
	c := &Connection{
		conn:        conn,
		connID:      0, // client ignore
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

	// Inherited properties from server (从client继承过来的属性)
	c.packet = client.GetPacket()
	c.onConnStart = client.GetOnConnStart()
	c.onConnStop = client.GetOnConnStop()
	c.msgHandler = client.GetMsgHandler()

	return c
}

// StartWriter is the goroutine that writes messages to the client
// (写消息Goroutine， 用户将数据发送给客户端)
func (c *Connection) StartWriter() {
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

// StartReader is a goroutine that reads data from the client
// (读消息Goroutine，用于从客户端中读取数据)
func (c *Connection) StartReader() {
	zlog.Ins().InfoF("[Reader Goroutine is running]")
	defer zlog.Ins().InfoF("%s [conn Reader exit!]", c.RemoteAddr().String())
	defer c.Stop()
	defer func() {
		if err := recover(); err != nil {
			zlog.Ins().ErrorF("connID=%d, panic err=%v", c.GetConnID(), err)
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// add by uuxia 2023-02-03
			buffer := make([]byte, zconf.GlobalObject.IOReadBuffSize)

			// read data from the connection's IO into the memory buffer
			// (从conn的IO中读取数据到内存缓冲buffer中)
			n, err := c.conn.Read(buffer)
			if err != nil {
				zlog.Ins().ErrorF("read msg head [read datalen=%d], error = %s", n, err)
				return
			}
			fmt.Println("-------------------------------------------buffer----------------------------", n)
			zlog.Ins().DebugF("read buffer %s \n", hex.EncodeToString(buffer[0:n]))

			// If normal data is read from the peer, update the heartbeat detection Active state
			// (正常读取到对端数据，更新心跳检测Active状态)
			if n > 0 && c.hc != nil {
				c.updateActivity()
			}

			// Deal with the custom protocol fragmentation problem, added by uuxia 2023-03-21
			// (处理自定义协议断粘包问题)
			if c.frameDecoder != nil {
				fmt.Println("----------------------------------------------------Came Here---------------------------------------------------")
				// Decode the 0-n bytes of data read
				// (为读取到的0-n个字节的数据进行解码)
				data := hex.EncodeToString(buffer[0:n])
				fmt.Println("before the data is", data)
				contains := strings.HasPrefix(data, "00000001000000b1")
				if !contains {
					data = "00000001000000b1" + data
				}
				fmt.Println("after the data is", data)
				byteDate, err := hex.DecodeString(data)
				if err != nil {
					fmt.Println("error while byte conversion", err)
					continue
				}
				bufArrays := c.frameDecoder.Decode(byteDate)
				if bufArrays == nil {
					fmt.Println("----------------------------bufarray is empty--------------------------")
					continue
				}
				for _, bytes := range bufArrays {
					zlog.Ins().DebugF("read buffer ------------------------------------>%s \n", hex.EncodeToString(bytes))
					msg := zpack.NewMessage(uint32(len(bytes)), bytes)
					// Get the current client's Request data
					// (得到当前客户端请求的Request数据)
					req := NewRequest(c, msg)
					c.msgHandler.Execute(req)
				}
			} else {
				msg := zpack.NewMessage(uint32(n), buffer[0:n])
				// Get the current client's Request data
				// (得到当前客户端请求的Request数据)
				req := NewRequest(c, msg)
				c.msgHandler.Execute(req)
			}
		}
	}
}

// Start starts the connection and makes the current connection work.
// (启动连接，让当前连接开始工作)
func (c *Connection) Start() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Ins().ErrorF("Connection Start() error: %v", err)
		}
	}()
	c.ctx, c.cancel = context.WithCancel(context.Background())

	// Execute the hook method for processing business logic when creating a connection
	// (按照用户传递进来的创建连接时需要处理的业务，执行钩子方法)
	c.callOnConnStart()

	// Start heartbeating detection
	if c.hc != nil {
		c.hc.Start()
		c.updateActivity()
	}

	// 占用workerid
	c.workerID = useWorker(c)

	// Start the Goroutine for reading data from the client
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

// Stop stops the connection and ends the current connection state.
// (停止连接，结束当前连接状态)
func (c *Connection) Stop() {
	c.cancel()
}

func (c *Connection) GetConnection() net.Conn {
	return c.conn
}

func (c *Connection) GetWsConn() *websocket.Conn {
	return nil
}

// Deprecated: use GetConnection instead
func (c *Connection) GetTCPConnection() net.Conn {
	return c.conn
}

func (c *Connection) GetConnID() uint64 {
	return c.connID
}

func (c *Connection) GetWorkerID() uint32 {
	return c.workerID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Connection) Send(data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	_, err := c.conn.Write(data)
	if err != nil {
		zlog.Ins().ErrorF("SendMsg err data = %+v, err = %+v", data, err)
		return err
	}

	return nil
}

func (c *Connection) SendToQueue(data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()

	if c.msgBuffChan == nil {
		c.msgBuffChan = make(chan []byte, zconf.GlobalObject.MaxMsgChanLen)
		// Start a Goroutine to write data back to the client
		// This method only reads data from the MsgBuffChan without allocating memory or starting a Goroutine
		// (开启用于写回客户端数据流程的Goroutine
		// 此方法只读取MsgBuffChan中的数据没调用SendBuffMsg可以分配内存和启用协程)
		go c.StartWriter()
	}

	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}

	if data == nil {
		zlog.Ins().ErrorF("Pack data is nil")
		return errors.New("Pack data is nil")
	}

	// Send timeout
	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- data:
		return nil
	}
}

// SendMsg directly sends Message data to the remote TCP client.
// (直接将Message数据发送数据给远程的TCP客户端)
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	// Pack data and send it
	msg, err := c.packet.Pack(zpack.NewMsgPackage(msgID, data))
	if err != nil {
		zlog.Ins().ErrorF("Pack error msg ID = %d", msgID)
		return errors.New("Pack error msg ")
	}

	err = c.Send(msg)
	if err != nil {
		zlog.Ins().ErrorF("SendMsg err msg ID = %d, data = %+v, err = %+v", msgID, string(msg), err)
		return err
	}

	return nil
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send buff msg")
	}
	if c.msgBuffChan == nil {
		c.msgBuffChan = make(chan []byte, zconf.GlobalObject.MaxMsgChanLen)
		// Start a Goroutine to write data back to the client
		// This method only reads data from the MsgBuffChan without allocating memory or starting a Goroutine
		// (开启用于写回客户端数据流程的Goroutine
		// 此方法只读取MsgBuffChan中的数据没调用SendBuffMsg可以分配内存和启用协程)
		go c.StartWriter()
	}

	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	msg, err := c.packet.Pack(zpack.NewMsgPackage(msgID, data))
	if err != nil {
		zlog.Ins().ErrorF("Pack error msg ID = %d", msgID)
		return errors.New("Pack error msg ")
	}

	// send timeout
	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) finalizer() {
	// Call the callback function registered by the user when closing the connection if it exists
	// (如果用户注册了该链接的	关闭回调业务，那么在此刻应该显示调用)
	c.callOnConnStop()

	c.msgLock.Lock()
	defer c.msgLock.Unlock()

	// If the connection has already been closed
	if c.isClosed == true {
		return
	}

	// Stop the heartbeat detector associated with the connection
	if c.hc != nil {
		c.hc.Stop()
	}

	// Close the socket connection
	_ = c.conn.Close()

	// Remove the connection from the connection manager
	if c.connManager != nil {
		c.connManager.Remove(c)
	}

	// Close all channels associated with the connection
	if c.msgBuffChan != nil {
		close(c.msgBuffChan)
	}

	c.isClosed = true

	zlog.Ins().InfoF("Conn Stop()...ConnID = %d", c.connID)
}

func (c *Connection) callOnConnStart() {
	if c.onConnStart != nil {
		zlog.Ins().InfoF("ZINX CallOnConnStart....")
		c.onConnStart(c)
	}
}

func (c *Connection) callOnConnStop() {
	if c.onConnStop != nil {
		zlog.Ins().InfoF("ZINX CallOnConnStop....")
		c.onConnStop(c)
	}
}

func (c *Connection) IsAlive() bool {
	if c.isClosed {
		return false
	}
	// Check the last activity time of the connection. If it's beyond the heartbeat interval,
	// then the connection is considered dead.
	// (检查连接最后一次活动时间，如果超过心跳间隔，则认为连接已经死亡)
	return time.Now().Sub(c.lastActivityTime) < zconf.GlobalObject.HeartbeatMaxDuration()
}

func (c *Connection) updateActivity() {
	c.lastActivityTime = time.Now()
}

func (c *Connection) SetHeartBeat(checker ziface.IHeartbeatChecker) {
	c.hc = checker
}

func (c *Connection) LocalAddrString() string {
	return c.localAddr
}

func (c *Connection) RemoteAddrString() string {
	return c.remoteAddr
}

func (c *Connection) GetName() string {
	return c.name
}

func (c *Connection) GetMsgHandler() ziface.IMsgHandle {
	return c.msgHandler
}
