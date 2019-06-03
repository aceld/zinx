package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"
	"zinx/zinx_app_demo/mmo_game/pb"
)

type Message struct {
	Len   uint32
	MsgId uint32
	Data  []byte
}

type TcpClient struct {
	conn     net.Conn
	X        float32
	Y        float32
	Z        float32
	V        float32
	Pid      int32
	isOnline chan bool
}

func (this *TcpClient) Unpack(headdata []byte) (head *Message, err error) {
	headbuf := bytes.NewReader(headdata)

	head = &Message{}

	// 读取Len
	if err = binary.Read(headbuf, binary.LittleEndian, &head.Len); err != nil {
		return nil, err
	}

	// 读取MsgId
	if err = binary.Read(headbuf, binary.LittleEndian, &head.MsgId); err != nil {
		return nil, err
	}

	// 封包太大
	//if head.Len > MaxPacketSize {
	//	return nil, packageTooBig
	//}

	return head, nil
}

func (this *TcpClient) Pack(msgId uint32, dataBytes []byte) (out []byte, err error) {
	outbuff := bytes.NewBuffer([]byte{})
	// 写Len
	if err = binary.Write(outbuff, binary.LittleEndian, uint32(len(dataBytes))); err != nil {
		return
	}
	// 写MsgId
	if err = binary.Write(outbuff, binary.LittleEndian, msgId); err != nil {
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
	binary_data, err := proto.Marshal(data)
	if err != nil {
		fmt.Println(fmt.Sprintf("marshaling error:  %s", err))
		return
	}

	sendData, err := this.Pack(msgID, binary_data)
	if err == nil {
		this.conn.Write(sendData)
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
		content := fmt.Sprintf("hello 我是player %d, 你是谁?", this.Pid)
		msg := &pb.Talk{
			Content: content,
		}
		this.SendMsg(2, msg)
	} else {
		//移动
		x := this.X
		z := this.Z

		randpos := rand.Intn(2)
		if randpos == 0 {
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
		randv := rand.Intn(2)
		v := this.V
		if randv == 0 {
			v = 25
		} else {
			v = 335
		}
		//封装Postsition消息
		msg := &pb.Position{
			X: x,
			Y: this.Y,
			Z: z,
			V: v,
		}

		fmt.Println(fmt.Sprintf("player ID: %d. Walking...", this.Pid))
		//发送移动MsgID:3的指令
		this.SendMsg(3, msg)
	}
}

/*
	处理一个回执业务
*/
func (this *TcpClient) DoMsg(msg *Message) {
	//处理消息
	//fmt.Println(fmt.Sprintf("msg id :%d, data len: %d", msg.MsgId, msg.Len))
	if msg.MsgId == 1 {
		//服务器回执给客户端 分配ID

		//解析proto
		syncpid := &pb.SyncPid{}
		proto.Unmarshal(msg.Data, syncpid)

		//给当前客户端ID进行赋值
		this.Pid = syncpid.Pid
	} else if msg.MsgId == 200 {
		//服务器回执客户端广播数据

		//解析proto
		bdata := &pb.BroadCast{}
		proto.Unmarshal(msg.Data, bdata)

		//初次玩家上线 广播位置消息
		if bdata.Tp == 2 && bdata.Pid == this.Pid {
			//本人
			//更新客户端坐标
			this.X = bdata.GetP().X
			this.Y = bdata.GetP().Y
			this.Z = bdata.GetP().Z
			this.V = bdata.GetP().V
			fmt.Println(fmt.Sprintf("player ID: %d online.. at(%f,%f,%f,%f)", bdata.Pid, this.X, this.Y, this.Z, this.V))

			//玩家已经成功上线
			this.isOnline <- true

		} else if bdata.Tp == 1 {
			fmt.Println(fmt.Sprintf("世界聊天,玩家%d说的话是: %s", bdata.Pid, bdata.GetContent()))
		}
	}
}

func (this *TcpClient) Start() {
	go func() {
		for {
			//read per head data
			headdata := make([]byte, 8)

			if _, err := io.ReadFull(this.conn, headdata); err != nil {
				fmt.Println(err)
				return
			}
			pkgHead, err := this.Unpack(headdata)
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

	select {
	case <-this.isOnline:
		//自动AI业务
		go func() {
			for {
				this.AIRobotAction()
				time.Sleep(3 * time.Second)
			}
		}()
	}
}

func NewTcpClient(ip string, port int) *TcpClient {
	addrStr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addrStr)
	if err == nil {
		client := &TcpClient{
			conn:     conn,
			Pid:      0,
			X:        0,
			Y:        0,
			Z:        0,
			V:        0,
			isOnline: make(chan bool),
		}
		return client
	} else {
		panic(err)
	}
}

func main() {
	for i := 0; i < 1000; i++ {
		client := NewTcpClient("127.0.0.1", 8999)
		client.Start()
		time.Sleep(1 * time.Second)
	}

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("=======", sig)
}
