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

// SmallPacket 小数据包 (16字节)
type SmallPacket struct {
	ID   int32
	Data [12]byte
}

func (p *SmallPacket) Encode() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint32(buf[0:4], uint32(p.ID))
	copy(buf[4:16], p.Data[:])
	return buf
}

// MediumPacket 中等数据包 (64字节)
type MediumPacket struct {
	ID   int32
	Seq  int32
	Data [56]byte
}

func (p *MediumPacket) Encode() []byte {
	buf := make([]byte, 64)
	binary.BigEndian.PutUint32(buf[0:4], uint32(p.ID))
	binary.BigEndian.PutUint32(buf[4:8], uint32(p.Seq))
	copy(buf[8:64], p.Data[:])
	return buf
}

// LargePacket 大数据包 (256字节)
type LargePacket struct {
	ID      int32
	Seq     int32
	Counter int32
	Data    [244]byte
}

func (p *LargePacket) Encode() []byte {
	buf := make([]byte, 256)
	binary.BigEndian.PutUint32(buf[0:4], uint32(p.ID))
	binary.BigEndian.PutUint32(buf[4:8], uint32(p.Seq))
	binary.BigEndian.PutUint32(buf[8:12], uint32(p.Counter))
	copy(buf[12:256], p.Data[:])
	return buf
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

// CalculateRouter 计算路由 (命令号 1001)
type CalculateRouter struct{}

func (c *CalculateRouter) PreHandle(request ziface.IRequest) {
}

func (c *CalculateRouter) Handle(request ziface.IRequest) {
	data := request.GetData()
	conn := request.GetConnection()
	packet := Decode(data)
	if packet == nil {
		zlog.Ins().ErrorF("Invalid packet data from %s", conn.RemoteAddr().String())
		return
	}
	packet.C = packet.A + packet.B
	err := conn.SendBuffMsg(1001, packet.Encode())
	if err != nil {
		zlog.Ins().ErrorF("SendBuffMsg error: %v", err)
	}
}

func (c *CalculateRouter) PostHandle(request ziface.IRequest) {
}

// SmallRouter 小数据路由 (命令号 2001)
type SmallRouter struct{}

func (s *SmallRouter) PreHandle(request ziface.IRequest) {
}

func (s *SmallRouter) Handle(request ziface.IRequest) {
	data := request.GetData()
	conn := request.GetConnection()
	if len(data) < 16 {
		zlog.Ins().ErrorF("Invalid small packet from %s", conn.RemoteAddr().String())
		return
	}
	pkt := &SmallPacket{
		ID: int32(binary.BigEndian.Uint32(data[0:4])),
	}
	copy(pkt.Data[:], data[4:16])
	pkt.ID++
	err := conn.SendBuffMsg(2001, pkt.Encode())
	if err != nil {
		zlog.Ins().ErrorF("SendBuffMsg error: %v", err)
	}
}

func (s *SmallRouter) PostHandle(request ziface.IRequest) {
}

// MediumRouter 中等数据路由 (命令号 3001)
type MediumRouter struct{}

func (m *MediumRouter) PreHandle(request ziface.IRequest) {
}

func (m *MediumRouter) Handle(request ziface.IRequest) {
	data := request.GetData()
	conn := request.GetConnection()
	if len(data) < 64 {
		zlog.Ins().ErrorF("Invalid medium packet from %s", conn.RemoteAddr().String())
		return
	}
	pkt := &MediumPacket{
		ID:  int32(binary.BigEndian.Uint32(data[0:4])),
		Seq: int32(binary.BigEndian.Uint32(data[4:8])),
	}
	copy(pkt.Data[:], data[8:64])
	pkt.Seq++
	err := conn.SendBuffMsg(3001, pkt.Encode())
	if err != nil {
		zlog.Ins().ErrorF("SendBuffMsg error: %v", err)
	}
}

func (m *MediumRouter) PostHandle(request ziface.IRequest) {
}

// LargeRouter 大数据路由 (命令号 4001)
type LargeRouter struct{}

func (l *LargeRouter) PreHandle(request ziface.IRequest) {
}

func (l *LargeRouter) Handle(request ziface.IRequest) {
	data := request.GetData()
	conn := request.GetConnection()
	if len(data) < 256 {
		zlog.Ins().ErrorF("Invalid large packet from %s", conn.RemoteAddr().String())
		return
	}
	pkt := &LargePacket{
		ID:      int32(binary.BigEndian.Uint32(data[0:4])),
		Seq:     int32(binary.BigEndian.Uint32(data[4:8])),
		Counter: int32(binary.BigEndian.Uint32(data[8:12])),
	}
	copy(pkt.Data[:], data[12:256])
	pkt.Counter++
	err := conn.SendBuffMsg(4001, pkt.Encode())
	if err != nil {
		zlog.Ins().ErrorF("SendBuffMsg error: %v", err)
	}
}

func (l *LargeRouter) PostHandle(request ziface.IRequest) {
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

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Println("Connection stats", "time", time.Now().Format("2006-01-02 15:04:05"), "clients", count, "goroutines", runtime.NumGoroutine(), "MEM:", m.Alloc/1024/1024, "MB")
		}
	}()

	// 注册路由
	s.AddRouter(1001, &CalculateRouter{})
	s.AddRouter(2001, &SmallRouter{})
	s.AddRouter(3001, &MediumRouter{})
	s.AddRouter(4001, &LargeRouter{})

	fmt.Printf("C10K Test Server starting on :%d\n", zconf.GlobalObject.TCPPort)

	// 启动服务
	s.Serve()
}
