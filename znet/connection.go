package znet

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/chnkenc/zinx-xiaoan/utils"
	"github.com/chnkenc/zinx-xiaoan/ziface"
)

// Connection 链接
type Connection struct {
	// 当前Conn属于哪个Server
	TCPServer ziface.IServer
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	// 消息管理MsgID和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandle
	// 告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	// 有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	sync.RWMutex
	// 链接属性
	property map[string]interface{}
	// 保护当前property的锁
	propertyLock sync.Mutex
	// 当前连接的关闭状态
	isClosed bool
}

// NewConnection 创建连接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	// 初始化Conn属性
	c := &Connection{
		TCPServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		MsgHandler:  msgHandler,
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:    nil,
	}

	// 将新创建的Conn添加到链接管理中
	c.TCPServer.GetConnMgr().Add(c)
	return c
}

// StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *Connection) StartWriter() {
	logger.Infof(
		"[Zinx][Connection][StartWriter]Writer Goroutine is Running, ConnID: %d, Remote Addr: %s",
		c.ConnID,
		c.RemoteAddr().String(),
	)
	defer logger.Infof(
		"[Zinx][Connection][StartWriter]Conn Writer Exit! ConnID: %d, Remote Addr: %s",
		c.ConnID,
		c.RemoteAddr().String(),
	)

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				// 有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					logger.Errorf(
						"[Zinx][Connection][StartWriter]Send Buff Data Error, Conn Writer Exit, ConnID: %d, Remote Addr: %s, Error: %v",
						c.ConnID,
						c.RemoteAddr().String(),
						err,
					)
					return
				}
			} else {
				logger.Errorf(
					"[Zinx][Connection][StartWriter]msgBuffChan is Closed, ConnID: %d, Remote Addr: %s",
					c.ConnID,
					c.RemoteAddr().String(),
				)
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	logger.Infof(
		"[Zinx][Connection][StartReader]Reader Goroutine is Running, ConnID: %d, Remote Addr: %s",
		c.ConnID,
		c.RemoteAddr().String(),
	)

	defer logger.Infof(
		"[Zinx][Connection][StartReader]Conn Reader Exit! ConnID: %d, Remote Addr: %s",
		c.ConnID,
		c.RemoteAddr().String(),
	)
	defer c.Stop()

	// 创建拆包解包的对象
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 读取客户端的Msg head
			headData := make([]byte, c.TCPServer.Packet().GetHeadLen())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				logger.Errorf(
					"[Zinx][Connection][StartReader]Read Msg Head Error, ConnID: %d, Remote Addr: %s, Error: %v",
					c.ConnID,
					c.RemoteAddr().String(),
					err,
				)
				return
			}

			// 拆包，得到msgID 和 datalen 放在msg中
			msg, err := c.TCPServer.Packet().Unpack(headData)
			if err != nil {
				logger.Errorf(
					"[Zinx][Connection][StartReader]Unpack Error, ConnID: %d, Remote Addr: %s, Error: %v",
					c.ConnID,
					c.RemoteAddr().String(),
					err,
				)
				return
			}

			// 根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					logger.Errorf(
						"[Zinx][Connection][StartReader]Read Msg Data Error, ConnID: %d, Remote Addr: %s, Error: %v",
						c.ConnID,
						c.RemoteAddr().String(),
						err,
					)
					return
				}
			}

			msg.SetHeaderData(headData)
			msg.SetData(data)
			// 得到当前客户端请求的Request数据
			req := Request{
				conn: c,
				msg:  msg,
			}

			if utils.GlobalObject.WorkerPoolSize > 0 {
				// 已经启动工作池机制，将消息交给Worker处理
				c.MsgHandler.SendMsgToTaskQueue(&req)
			} else {
				// 从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.MsgHandler.DoMsgHandler(&req)
			}
		}
	}
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	// 1：开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	// 2：开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	// 按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TCPServer.CallOnConnStart(c)

	<-c.ctx.Done()
	c.finalizer()
}

// Stop 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	c.cancel()
}

// GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgID, sn uint8, data []byte) error {
	_, err := c.SendMsgPackage(NewMsgPackage(msgID, sn, data))
	if err != nil {
		return err
	}

	return nil
}

// SendMsgPackage 发送Message包
func (c *Connection) SendMsgPackage(msgPackage ziface.IMessage) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed {
		return nil, errors.New("connection closed when send msg")
	}

	// 将data封包，并且发送
	dp := c.TCPServer.Packet()
	msg, err := dp.Pack(msgPackage)
	if err != nil {
		logger.Errorf(
			"[Zinx][Connection][SendMsgPackage]Pack Error, ConnID: %d, Remote Addr: %s, Msg Pkg: %+v, Error: %v",
			c.ConnID,
			c.RemoteAddr().String(),
			msgPackage,
			err,
		)
		return nil, err
	}

	_, err = c.Conn.Write(msg)
	if err != nil {
		logger.Errorf(
			"[Zinx][Connection][SendMsgPackage]Write Data Error, ConnID: %d, Remote Addr: %s, Msg Pkg: %+v, Error: %v",
			c.ConnID,
			c.RemoteAddr().String(),
			msgPackage,
			err,
		)
		return nil, err
	}

	return msg, nil
}

// SendBuffMsg 发送BuffMsg
func (c *Connection) SendBuffMsg(msgID, sn uint8, data []byte) error {
	_, err := c.SendBuffMsgPackage(NewMsgPackage(msgID, sn, data))
	if err != nil {
		return err
	}

	return nil
}

// SendBuffMsgPackage 发送BuffMsg包
func (c *Connection) SendBuffMsgPackage(msgPackage ziface.IMessage) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()
	if c.isClosed {
		return nil, errors.New("connection closed when send buff msg")
	}

	// 将data封包，并且发送
	dp := c.TCPServer.Packet()
	msg, err := dp.Pack(msgPackage)
	if err != nil {
		logger.Errorf(
			"[Zinx][Connection][SendBuffMsgPackage]Pack Error, ConnID: %d, Msg Pkg: %+v, Remote Addr: %s, Error: %v",
			c.ConnID,
			msgPackage,
			c.RemoteAddr().String(),
			err,
		)
		return nil, err
	}

	// 发送超时
	select {
	case <-idleTimeout.C:
		return nil, errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return msg, nil
	}
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// 返回ctx，用于用户自定义的go程获取连接退出状态
func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) finalizer() {
	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TCPServer.CallOnConnStop(c)

	c.Lock()
	defer c.Unlock()

	//如果当前链接已经关闭
	if c.isClosed {
		return
	}

	logger.Infof(
		"[Zinx][Connection][finalizer]Conn Stop()...ConnID: %d, Remote Addr: %s",
		c.ConnID,
		c.RemoteAddr().String(),
	)

	// 关闭socket链接
	_ = c.Conn.Close()

	//将链接从连接管理器中删除
	c.TCPServer.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.msgBuffChan)

	//设置标志位
	c.isClosed = true
}
