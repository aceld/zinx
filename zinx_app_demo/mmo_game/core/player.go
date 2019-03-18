package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
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

/*
	Player ID 生成器
 */
var PidGen int32 = 1   //用来生成玩家ID的计数器
var IdLock sync.Mutex   //保护PidGen的互斥机制

//创建一个玩家对象
func NewPlayer(conn ziface.IConnection) *Player {
	//生成一个PID
	IdLock.Lock()
	id := PidGen
	PidGen ++
	IdLock.Unlock()

	p := &Player{
		Pid : id,
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

	//组建MsgId200 proto数据
	msg := &pb.BroadCast{
		Pid:p.Pid,
		Tp:2,//TP2 代表广播坐标
		Data: &pb.BroadCast_P{
			P:&pb.Position{
				X:p.X,
				Y:p.Y,
				Z:p.Z,
				V:p.V,
			},
		},
	}

	//发送数据给客户端
	p.SendMsg(200, msg)
}

//给当前玩家周边的(九宫格内)玩家广播自己的位置，让他们显示自己
func (p *Player) SyncSurrounding() {
	//1 根据自己的位置，获取周围九宫格内的玩家pid
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	//2 根据pid得到所有玩家对象
	players := make([]*Player, 0, len(pids))
	//3 给这些玩家发送MsgID:200消息，让自己出现在对方视野中
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	//3.1 组建MsgId200 proto数据
	msg := &pb.BroadCast{
		Pid:p.Pid,
		Tp:2,//TP2 代表广播坐标
		Data: &pb.BroadCast_P{
			P:&pb.Position{
				X:p.X,
				Y:p.Y,
				Z:p.Z,
				V:p.V,
			},
		},
	}
	//3.2 每个玩家分别给对应的客户端发送200消息，显示人物
	for _, player := range players {
		player.SendMsg(200, msg)
	}
	//4 让周围九宫格内的玩家出现在自己的视野中
	//4.1 制作Message SyncPlayers 数据
	playersData := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			Pid:player.Pid,
			P:&pb.Position{
				X:player.X,
				Y:player.Y,
				Z:player.Z,
				V:player.V,
			},
		}
		playersData = append(playersData, p)
	}

	//4.2 封装SyncPlayer protobuf数据
	SyncPlayersMsg := &pb.SyncPlayers{
		Ps:playersData[:],
	}

	//4.3 给当前玩家发送需要显示周围的全部玩家数据
	p.SendMsg(202, SyncPlayersMsg)
}

//广播玩家聊天
func (p *Player) Talk(content string) {
	//1. 组建MsgId200 proto数据
	msg := &pb.BroadCast{
		Pid:p.Pid,
		Tp:1,//TP 1 代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//2. 得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//3. 向所有的玩家发送MsgId:200消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}


//广播玩家位置移动
func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {
	//更新玩家的位置信息
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	//组装protobuf协议，发送位置给周围玩家
	msg := &pb.BroadCast{
		Pid:p.Pid,
		Tp:4,     //4- 移动之后的坐标信息
		Data: &pb.BroadCast_P{
			P:&pb.Position{
				X:p.X,
				Y:p.Y,
				Z:p.Z,
				V:p.V,
			},
		},
	}

	//获取当前玩家周边全部玩家
	players := p.GetSurroundingPlayers()
	//向周边的每个玩家发送MsgID:200消息，移动位置更新消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

//获得当前玩家的AOI周边玩家信息
func (p *Player) GetSurroundingPlayers() []*Player {
	//得到当前AOI区域的所有pid
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)

	//将所有pid对应的Player放到Player切片中
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

//玩家下线
func (p *Player) LostConnection() {
	//1 获取周围AOI九宫格内的玩家
	players := p.GetSurroundingPlayers()

	//2 封装MsgID:201消息
	msg := &pb.SyncPid{
		Pid:p.Pid,
	}

	//3 向周围玩家发送消息
	for _, player := range players {
		player.SendMsg(201, msg)
	}

	//4 世界管理器将当前玩家从AOI中摘除
	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayerByPid(p.Pid)
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