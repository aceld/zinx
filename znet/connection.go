/*
	服务端Server的链接模块
*/
package znet

import (
	"context"
	"errors"
	"fmt"
	"github.com/aceld/zinx/zpack"
	"io"
	"net"
	"sync"
	"time"

	"github.com/aceld/zinx/utils"
	"github.com/aceld/zinx/ziface"
)

//Connection Tcp连接模块
//用于处理Tcp连接的读写业务 一个连接对应一个Connection
type Connection struct {
	//当前连接的socket TCP套接字
	conn net.Conn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一 ，服务端Connection使用
	connID uint32
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
	////保护当前property的锁
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
}

//newServerConn :for Server, 创建一个Server服务端特性的连接的方法
//Note: 名字由 NewConnection 更变
func newServerConn(server ziface.IServer, conn net.Conn, connID uint32) *Connection {
	//初始化Conn属性
	c := &Connection{
		conn:        conn,
		connID:      connID,
		isClosed:    false,
		msgBuffChan: nil,
		property:    nil,
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

//newClientConn :for Client, 创建一个Client服务端特性的连接的方法
func newClientConn(client ziface.IClient, conn net.Conn) *Connection {
	c := &Connection{
		conn:        conn,
		connID:      0, //client ignore
		isClosed:    false,
		msgBuffChan: nil,
		property:    nil,
	}

	//从client继承过来的属性
	c.packet = client.GetPacket()
	c.onConnStart = client.GetOnConnStart()
	c.onConnStop = client.GetOnConnStop()
	c.msgHandler = client.GetMsgHandler()

	return c
}

//StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

//StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Stop()

	// 创建拆包解包的对象
	for {
		select {
		case <-c.ctx.Done():
			return
		default:

			//读取客户端的Msg head
			headData := make([]byte, c.packet.GetHeadLen())
			if _, err := io.ReadFull(c.conn, headData); err != nil {
				fmt.Println("read msg head error ", err)
				return
			}
			//fmt.Printf("read headData %+v\n", headData)

			//拆包，得到msgID 和 datalen 放在msg中
			msg, err := c.packet.Unpack(headData)
			if err != nil {
				fmt.Println("unpack error ", err)
				return
			}
			//fmt.Printf("read msg %+v\n", msg)

			//根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.conn, data); err != nil {
					fmt.Println("read msg data error ", err)
					return
				}
			}
			msg.SetData(data)

			//得到当前客户端请求的Request数据
			req := NewRequest(c, msg)

			if utils.GlobalObject.WorkerPoolSize > 0 {
				//已经启动工作池机制，将消息交给Worker处理
				c.msgHandler.SendMsgToTaskQueue(req)
			} else {
				//从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.msgHandler.DoMsgHandler(req)
			}
		}
	}
}

//Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.callOnConnStart()
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()

	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}

//Stop 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	c.cancel()
}

func (c *Connection) GetConnection() net.Conn {
	return c.conn
}

//GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.connID
}

//RemoteAddr 获取链接远程地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

//LocalAddr 获取链接本地地址信息
func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

//SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	//将data封包，并且发送
	msg, err := c.packet.Pack(zpack.NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	_, err = c.conn.Write(msg)
	return err
}

//SendBuffMsg  发生BuffMsg
func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()
	if c.msgBuffChan == nil {
		c.msgBuffChan = make(chan []byte, utils.GlobalObject.MaxMsgChanLen)
		//开启用于写回客户端数据流程的Goroutine
		//此方法只读取MsgBuffChan中的数据没调用SendBuffMsg可以分配内存和启用协程
		go c.StartWriter()
	}
	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}

	//将data封包，并且发送
	msg, err := c.packet.Pack(zpack.NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}

	// 发送超时
	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}

	return nil
}

//SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

//GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

//RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

//返回ctx，用于用户自定义的go程获取连接退出状态
func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) finalizer() {
	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.callOnConnStop()

	c.msgLock.Lock()
	defer c.msgLock.Unlock()

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	fmt.Println("Conn Stop()...ConnID = ", c.connID)

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
}

//callOnConnStart 调用连接OnConnStart Hook函数
func (c *Connection) callOnConnStart() {
	if c.onConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		c.onConnStart(c)
	}
}

//callOnConnStart 调用连接OnConnStop Hook函数
func (c *Connection) callOnConnStop() {
	if c.onConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		c.onConnStop(c)
	}
}

/*
func (c *Connection) IsAlive() bool {
	if c.isClosed {
		return false
	}
	// 检查连接最后一次活动时间，如果超过心跳间隔，则认为连接已经死亡
	return time.Now().Unix()-c.GetLastActivityTime().Unix() < int64(utils.GlobalObject.HeartbeatInterval.Seconds())
}

// GetLastActivityTime方法用于获取最后一次活动时间
func (c *Connection) GetLastActivityTime() time.Time {
	return time.Now()
}

*/
