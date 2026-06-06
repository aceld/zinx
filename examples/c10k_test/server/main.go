package main

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

// Packet 数据结构
type Packet struct {
	A int32
	B int32
	C int32
}

// Encode 将 Packet 编码为字节数组
func (p *Packet) Encode() []byte {
	buf := make([]byte, 12) // 3 * 4 bytes
	binary.BigEndian.PutUint32(buf[0:4], uint32(p.A))
	binary.BigEndian.PutUint32(buf[4:8], uint32(p.B))
	binary.BigEndian.PutUint32(buf[8:12], uint32(p.C))
	return buf
}

// Decode 将字节数组解码为 Packet
func Decode(buf []byte) *Packet {
	if len(buf) < 12 {
		return nil
	}
	return &Packet{
		A: int32(binary.BigEndian.Uint32(buf[0:4])),
		B: int32(binary.BigEndian.Uint32(buf[4:8])),
		C: int32(binary.BigEndian.Uint32(buf[8:12])),
	}
}

var (
	connCount int32
	mutex     sync.Mutex
)

func OnConnStart(conn ziface.IConnection) {
	mutex.Lock()
	connCount++
	mutex.Unlock()
	//logger.Info("Client connected", "conn_id", conn.GetConnID(), "addr", conn.RemoteAddrString(), "total", connCount)
}

func OnConnStop(conn ziface.IConnection) {
	mutex.Lock()
	connCount--
	mutex.Unlock()
	//logger.Info("Client disconnected", "conn_id", conn.GetConnID(), "addr", conn.RemoteAddrString(), "total", connCount)
}

// CalculateRouter 计算路由
type CalculateRouter struct{}

func (c *CalculateRouter) PreHandle(request ziface.IRequest) {
}

func (c *CalculateRouter) Handle(request ziface.IRequest) {
	// 获取消息数据
	data := request.GetData()
	conn := request.GetConnection()

	// 解码 Packet
	packet := Decode(data)
	if packet == nil {
		zlog.Ins().ErrorF("Invalid packet data from %s", conn.RemoteAddr().String())
		return
	}

	// 计算 C = A + B
	packet.C = packet.A + packet.B
	// 使用 SendBuffMsg 回复（异步方式）
	err := conn.SendBuffMsg(1, packet.Encode())
	if err != nil {
		zlog.Ins().ErrorF("SendBuffMsg error (first): %v, retrying...", err)
		// err = conn.SendBuffMsg(1, packet.Encode())
		// if err != nil {
		// 	zlog.Ins().ErrorF("SendBuffMsg error (retry): %v", err)
		// }
	}
}

func (c *CalculateRouter) PostHandle(request ziface.IRequest) {
}

func main() {

	// 创建服务器
	s := znet.NewServer()

	s.SetOnConnStart(OnConnStart)
	s.SetOnConnStop(OnConnStop)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			mutex.Lock()
			count := connCount
			mutex.Unlock()
			fmt.Println("Connection stats", "clients", count, "goroutines", runtime.NumGoroutine())
		}
	}()

	// 注册路由
	s.AddRouter(1001, &CalculateRouter{})

	fmt.Printf("C10K Test Server starting on :%d\n", zconf.GlobalObject.TCPPort)

	// 启动服务
	s.Serve()
}
