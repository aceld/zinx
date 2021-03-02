package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/aceld/zinx/zinx_app_demo/mmo_game/pb"
	"github.com/golang/protobuf/proto"
)

type Message struct {
	Len   uint32
	MsgID uint32
	Data  []byte
}

type TcpClient struct {
	conn     net.Conn
	X        float32
	Y        float32
	Z        float32
	V        float32
	PID      int32
	isOnline chan bool
}

func (this *TcpClient) Unpack(headdata []byte) (head *Message, err error) {
	headBuf := bytes.NewReader(headdata)

	head = &Message{}

	// 读取Len
	if err = binary.Read(headBuf, binary.LittleEndian, &head.Len); err != nil {
		return nil, err
	}

	// 读取MsgID
	if err = binary.Read(headBuf, binary.LittleEndian, &head.MsgID); err != nil {
		return nil, err
	}

	// 封包太大
	//if head.Len > MaxPacketSize {
	//	return nil, packageTooBig
	//}

	return head, nil
}

func (this *TcpClient) Pack(msgID uint32, dataBytes []byte) (out []byte, err error) {
	outbuff := bytes.NewBuffer([]byte{})
	// 写Len
	if err = binary.Write(outbuff, binary.LittleEndian, uint32(len(dataBytes))); err != nil {
		return
	}
	// 写MsgID
	if err = binary.Write(outbuff, binary.LittleEndian, msgID); err != nil {
		return
	}

	//all pkg data
	if err = binary.Write(outbuff, binary.LittleEndian, dataBytes); err != nil {
		return
	}

	out = outbuff.Bytes()

	return
}

func (this *TcpClient) SendMsg(msgID uint32, data proto.Message) {

	// 进行编码
	binaryData, err := proto.Marshal(data)
	if err != nil {
		fmt.Println(fmt.Sprintf("marshaling error:  %s", err))
		return
	}

	sendData, err := this.Pack(msgID, binaryData)
	if err == nil {
		_, _ = this.conn.Write(sendData)
	} else {
		fmt.Println(err)
	}

	return
}

func (this *TcpClient) AIRobotAction() {
	//聊天或者移动

	//随机获得动作
	tp := rand.Intn(2)
	if tp == 0 {
		content := fmt.Sprintf("hello 我是player %d, 你是谁?", this.PID)
		msg := &pb.Talk{
			Content: content,
		}
		this.SendMsg(2, msg)
	} else {
		//移动
		x := this.X
		z := this.Z

		randPos := rand.Intn(2)
		if randPos == 0 {
			x -= float32(rand.Intn(10))
			z -= float32(rand.Intn(10))
		} else {
			x += float32(rand.Intn(10))
			z += float32(rand.Intn(10))
		}

		//纠正坐标位置
		if x > 410 {
			x = 410
		} else if x < 85 {
			x = 85
		}

		if z > 400 {
			z = 400
		} else if z < 75 {
			z = 75
		}

		//移动方向角度
		randV := rand.Intn(2)
		v := this.V
		if randV == 0 {
			v = 25
		} else {
			v = 335
		}
		//封装Position消息
		msg := &pb.Position{
			X: x,
			Y: this.Y,
			Z: z,
			V: v,
		}

		fmt.Println(fmt.Sprintf("player ID: %d. Walking...", this.PID))
		//发送移动MsgID:3的指令
		this.SendMsg(3, msg)
	}
}

/*
	处理一个回执业务
*/
func (this *TcpClient) DoMsg(msg *Message) {
	//处理消息
	//fmt.Println(fmt.Sprintf("msg ID :%d, data len: %d", msg.MsgID, msg.Len))
	if msg.MsgID == 1 {
		//服务器回执给客户端 分配ID

		//解析proto
		syncpID := &pb.SyncPID{}
		_ = proto.Unmarshal(msg.Data, syncpID)

		//给当前客户端ID进行赋值
		this.PID = syncpID.PID
	} else if msg.MsgID == 200 {
		//服务器回执客户端广播数据

		//解析proto
		bdata := &pb.BroadCast{}
		_ = proto.Unmarshal(msg.Data, bdata)

		//初次玩家上线 广播位置消息
		if bdata.Tp == 2 && bdata.PID == this.PID {
			//本人
			//更新客户端坐标
			this.X = bdata.GetP().X
			this.Y = bdata.GetP().Y
			this.Z = bdata.GetP().Z
			this.V = bdata.GetP().V
			fmt.Println(fmt.Sprintf("player ID: %d online.. at(%f,%f,%f,%f)", bdata.PID, this.X, this.Y, this.Z, this.V))

			//玩家已经成功上线
			this.isOnline <- true

		} else if bdata.Tp == 1 {
			fmt.Println(fmt.Sprintf("世界聊天,玩家%d说的话是: %s", bdata.PID, bdata.GetContent()))
		}
	}
}

func (this *TcpClient) Start() {
	go func() {
		for {
			//读取服务端发来的数据 ==》 SyncPID
			//1.读取8字节
			//第一次读取，读取数据头
			headData := make([]byte, 8)

			if _, err := io.ReadFull(this.conn, headData); err != nil {
				fmt.Println(err)
				return
			}
			pkgHead, err := this.Unpack(headData)
			if err != nil {
				return
			}
			//data
			if pkgHead.Len > 0 {
				pkgHead.Data = make([]byte, pkgHead.Len)
				if _, err := io.ReadFull(this.conn, pkgHead.Data); err != nil {
					return
				}
			}

			//处理服务器回执业务
			this.DoMsg(pkgHead)
		}
	}()

	// 10s后，断开连接
	for {
		select {
		case <-this.isOnline:
			go func() {
				for {
					this.AIRobotAction()
					time.Sleep(time.Second)
				}
			}()
		case <-time.After(time.Second * 10):
			_ = this.conn.Close()
			return
		}
	}
}

func NewTcpClient(ip string, port int) *TcpClient {
	addrStr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addrStr)
	if err != nil {
		panic(err)
	}

	client := &TcpClient{
		conn:     conn,
		PID:      0,
		X:        0,
		Y:        0,
		Z:        0,
		V:        0,
		isOnline: make(chan bool),
	}
	return client
}

func main() {
	// 开启一个waitgroup，同时运行3个goroutine

	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			client := NewTcpClient("127.0.0.1", 8999)
			client.Start()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			client := NewTcpClient("127.0.0.1", 8999)
			client.Start()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			client := NewTcpClient("127.0.0.1", 8999)
			client.Start()
		}
	}()

	fmt.Println("AI robot start")
	wg.Wait()
	fmt.Println("AI robot exit")
}
