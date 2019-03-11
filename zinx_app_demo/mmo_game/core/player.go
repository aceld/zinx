package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"zinx/ziface"
	"zinx/zinx_app_demo/mmo_game/pb"
)

//玩家对象
type Player struct {
	Pid int32  	//玩家ID
	Conn ziface.IConnection //当前玩家的连接
	X 	float32 //平面x坐标
	Y   float32 //高度
	Z   float32 //平面y坐标 (注意不是Y)
	V   float32 //旋转0-360度
}

//创建一个玩家对象
func NewPlayer(conn ziface.IConnection, pid int32) *Player {
	p := &Player{
		Pid : pid,
		Conn:conn,
		X:float32(160 + rand.Intn(10)),//随机在160坐标点 基于X轴偏移若干坐标
		Y:0, //高度为0
		Z:float32(134 + rand.Intn(17)), //随机在134坐标点 基于Y轴偏移若干坐标
		V:0, //角度为0，尚未实现
	}

	return p
}

//告知客户端pid,同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	//组建MsgId0 proto数据
	data := &pb.SyncPid{
		Pid:p.Pid,
	}

	//发送数据给客户端
	p.SendMsg(1, data)
}

//广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {

	msg := &pb.BroadCast{
		Pid:p.Pid,
		Tp:2,//TP2 代表广播坐标
		Data: &pb.BroadCast_P{
			&pb.Position{
				X:p.X,
				Y:p.Y,
				Z:p.Z,
				V:p.V,
			},
		},
	}

	p.SendMsg(200, msg)
}

//广播玩家的自身地理位置信息
func (p *Player) SyncSurrounding() {
	//TODO 根据自己的位子，获取周围九宫格内的玩家PID

	//TODO 根据获取的PID集合， 获取所有的玩家信息

	//TODO 构建自己的上线地点坐标

	//TODO
}


/*
	发送消息给客户端，
	主要是将pb的protobuf数据序列化之后发送
 */
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	fmt.Printf("before Marshal data = %+v\n", data)
	//将proto Message结构体序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}
	fmt.Printf("after Marshal data = %+v\n", msg)

	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	//调用Zinx框架的SendMsg发包
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error !")
		return
	}

	return
}