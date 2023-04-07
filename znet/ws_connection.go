package znet

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zinterceptor"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/zpack"
	"github.com/gorilla/websocket"
	"net"
	"sync"
	"time"
)

// WsConnection Websocket连接模块
// 用于处理 Websocket 连接的读写业务 一个连接对应一个Connection
type WsConnection struct {
	//当前连接的socket TCP套接字
	conn *websocket.Conn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一 ，服务端Connection使用
	//uint64 取值范围：0 ~ 18,446,744,073,709,551,615
	//这个是理论支持的进程connID的最大数量
	connID uint64
	//消息管理MsgID和对应处理方法的消息管理模块
	msgHandler ziface.IMsgHandle
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
	//用户收发消息的Lock
	msgLock sync.RWMutex
	//链接属性
	property map[string]interface{}
	//保护当前property的锁
	propertyLock sync.Mutex
	//当前连接的关闭状态
	isClosed bool
	//当前链接是属于哪个Connection Manager的
	connManager ziface.IConnManager
	//当前连接创建时Hook函数
	onConnStart func(conn ziface.IConnection)
	//当前连接断开时的Hook函数
	onConnStop func(conn ziface.IConnection)
	//数据报文封包方式
	packet ziface.IDataPack
	//最后一次活动时间
	lastActivityTime time.Time
	//断粘包解码器
	frameDecoder ziface.IFrameDecoder
	//心跳检测器
	hc ziface.IHeartbeatChecker
}

// newServerConn :for Server, 创建一个Server服务端特性的连接的方法
// Note: 名字由 NewConnection 更变
func newWebsocketConn(server ziface.IServer, conn *websocket.Conn, connID uint64) ziface.IConnection {
	//初始化Conn属性
	c := &WsConnection{
		conn:        conn,
		connID:      connID,
		isClosed:    false,
		msgBuffChan: nil,
		property:    nil,
	}

	lengthField := server.GetLengthField()
	if lengthField != nil {
		c.frameDecoder = zinterceptor.NewFrameDecoder(*lengthField)
	}

	//从server继承过来的属性
	c.packet = server.GetPacket()
	c.onConnStart = server.GetOnConnStart()
	c.onConnStop = server.GetOnConnStop()
	c.msgHandler = server.GetMsgHandler()

	//将当前的Connection与Server的ConnManager绑定
	c.connManager = server.GetConnMgr()

	//将新创建的Conn添加到链接管理中
	server.GetConnMgr().Add(c)

	return c
}

// newClientConn :for Client, 创建一个Client服务端特性的连接的方法
func newWsClientConn(client ziface.IClient, conn *websocket.Conn) ziface.IConnection {
	c := &WsConnection{
		conn:        conn,
		connID:      0, //client ignore
		isClosed:    false,
		msgBuffChan: nil,
		property:    nil,
	}

	lengthField := client.GetLengthField()
	if lengthField != nil {
		c.frameDecoder = zinterceptor.NewFrameDecoder(*lengthField)
	}

	//从client继承过来的属性
	c.packet = client.GetPacket()
	c.onConnStart = client.GetOnConnStart()
	c.onConnStop = client.GetOnConnStop()
	c.msgHandler = client.GetMsgHandler()

	return c
}

// StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *WsConnection) StartWriter() {
	zlog.Ins().InfoF("Writer Goroutine is running")
	defer zlog.Ins().InfoF("%s [conn Writer exit!]", c.RemoteAddr().String())

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给对端
				if err := c.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
					zlog.Ins().ErrorF("Send Buff Data error:, %s Conn Writer exit", err)
					break
				}

				//写对端成功, 更新链接活动时间
				//c.updateActivity()
			} else {
				zlog.Ins().ErrorF("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *WsConnection) StartReader() {
	zlog.Ins().InfoF("[Reader Goroutine is running]")
	defer zlog.Ins().InfoF("%s [conn Reader exit!]", c.RemoteAddr().String())
	defer c.Stop()

	// 创建拆包解包的对象
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			//add by uuxia 2023-02-03
			//从conn的IO中读取数据到内存缓冲buffer中
			messageType, buffer, err := c.conn.ReadMessage()
			if err != nil {
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

			//正常读取到对端数据，更新心跳检测Active状态
			if n > 0 && c.hc != nil {
				c.updateActivity()
			}

			//处理自定义协议断粘包问题 add by uuxia 2023-03-21
			if c.frameDecoder != nil {
				//为读取到的0-n个字节的数据进行解码
				bufArrays := c.frameDecoder.Decode(buffer)
				if bufArrays == nil {
					continue
				}
				for _, bytes := range bufArrays {
					zlog.Ins().DebugF("read buffer %s \n", hex.EncodeToString(bytes))
					msg := zpack.NewMessage(uint32(len(bytes)), bytes)
					//得到当前客户端请求的Request数据
					req := NewRequest(c, msg)
					c.msgHandler.Execute(req)
				}
			} else {
				msg := zpack.NewMessage(uint32(n), buffer[0:n])
				//得到当前客户端请求的Request数据
				req := NewRequest(c, msg)
				c.msgHandler.Execute(req)
			}
		}
	}
}

// Start 启动连接，让当前连接开始工作
func (c *WsConnection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.callOnConnStart()

	//启动心跳检测
	if c.hc != nil {
		c.hc.Start()
		c.updateActivity()
	}

	//开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()

	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}

// Stop 停止连接，结束当前连接状态M
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

// GetConnID 获取当前连接ID
func (c *WsConnection) GetConnID() uint64 {
	return c.connID
}

// RemoteAddr 获取链接远程地址信息
func (c *WsConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// LocalAddr 获取链接本地地址信息
func (c *WsConnection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *WsConnection) Send(data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()
	if c.isClosed == true {
		return errors.New("WsConnection closed when send msg")
	}

	//写回客户端
	err := c.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		zlog.Ins().ErrorF("SendMsg err data = %+v, err = %+v", data, err)
		return err
	}

	//写对端成功, 更新链接活动时间
	//c.updateActivity()

	return nil
}

func (c *WsConnection) SendToQueue(data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()

	if c.msgBuffChan == nil {
		c.msgBuffChan = make(chan []byte, zconf.GlobalObject.MaxMsgChanLen)
		//开启用于写回客户端数据流程的Goroutine
		//此方法只读取MsgBuffChan中的数据没调用SendBuffMsg可以分配内存和启用协程
		go c.StartWriter()
	}

	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed == true {
		return errors.New("WsConnection closed when send buff msg")
	}

	if data == nil {
		zlog.Ins().ErrorF("Pack data is nil")
		return errors.New("Pack data is nil")
	}

	// 发送超时
	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- data:
		return nil
	}
}

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *WsConnection) SendMsg(msgID uint32, data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()
	if c.isClosed == true {
		return errors.New("WsConnection closed when send msg")
	}

	//将data封包，并且发送
	msg, err := c.packet.Pack(zpack.NewMsgPackage(msgID, data))
	if err != nil {
		zlog.Ins().ErrorF("Pack error msg ID = %d", msgID)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	err = c.conn.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		zlog.Ins().ErrorF("SendMsg err msg ID = %d, data = %+v, err = %+v", msgID, string(msg), err)
		return err
	}

	//写对端成功, 更新链接活动时间
	//c.updateActivity()

	return nil
}

// SendBuffMsg  发生BuffMsg
func (c *WsConnection) SendBuffMsg(msgID uint32, data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()

	if c.msgBuffChan == nil {
		c.msgBuffChan = make(chan []byte, zconf.GlobalObject.MaxMsgChanLen)
		//开启用于写回客户端数据流程的Goroutine
		//此方法只读取MsgBuffChan中的数据没调用SendBuffMsg可以分配内存和启用协程
		go c.StartWriter()
	}

	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed == true {
		return errors.New("WsConnection closed when send buff msg")
	}

	//将data封包，并且发送
	msg, err := c.packet.Pack(zpack.NewMsgPackage(msgID, data))
	if err != nil {
		zlog.Ins().ErrorF("Pack error msg ID = %d", msgID)
		return errors.New("Pack error msg ")
	}

	// 发送超时
	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}
}

// SetProperty 设置链接属性
func (c *WsConnection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *WsConnection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

// RemoveProperty 移除链接属性
func (c *WsConnection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// 返回ctx，用于用户自定义的go程获取连接退出状态
func (c *WsConnection) Context() context.Context {
	return c.ctx
}

func (c *WsConnection) finalizer() {
	//如果用户注册了该链接的	关闭回调业务，那么在此刻应该显示调用
	c.callOnConnStop()

	c.msgLock.Lock()
	defer c.msgLock.Unlock()

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	//关闭链接绑定的心跳检测器
	if c.hc != nil {
		c.hc.Stop()
	}

	// 关闭socket链接
	_ = c.conn.Close()

	//将链接从连接管理器中删除
	if c.connManager != nil {
		c.connManager.Remove(c)
	}

	//关闭该链接全部管道
	if c.msgBuffChan != nil {
		close(c.msgBuffChan)
	}
	//设置标志位
	c.isClosed = true

	zlog.Ins().InfoF("Conn Stop()...ConnID = %d", c.connID)
}

// callOnConnStart 调用连接OnConnStart Hook函数
func (c *WsConnection) callOnConnStart() {
	if c.onConnStart != nil {
		zlog.Ins().InfoF("ZINX CallOnConnStart....")
		c.onConnStart(c)
	}
}

// callOnConnStart 调用连接OnConnStop Hook函数
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
	// 检查连接最后一次活动时间，如果超过心跳间隔，则认为连接已经死亡
	return time.Now().Sub(c.lastActivityTime) < zconf.GlobalObject.HeartbeatMaxDuration()
}

func (c *WsConnection) updateActivity() {
	c.lastActivityTime = time.Now()
}

func (c *WsConnection) SetHeartBeat(checker ziface.IHeartbeatChecker) {
	c.hc = checker
}
