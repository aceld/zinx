package main

import (
	"fmt"

	"github.com/aceld/zinx/zdecoder"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zinx_app_demo/mmo_game/api"
	"github.com/aceld/zinx/zinx_app_demo/mmo_game/core"
	"github.com/aceld/zinx/znet"
	"github.com/aceld/zinx/zpack"
)

// OnConnectionAdd is a hook function called when a client establishes a connection
// 当客户端建立连接的时候的hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	fmt.Println("=====> OnConnectionAdd is Called ...")
	// Create a new player
	// 创建一个玩家
	player := core.NewPlayer(conn)

	// Synchronize the current player's ID to the client using MsgID:1 message
	// 同步当前的PlayerID给客户端， 走MsgID:1 消息
	player.SyncPID()

	// Synchronize the initial coordinate information of the current player to the client using MsgID:200 message
	// 同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	player.BroadCastStartPosition()

	// Add the newly online player to the WorldManager
	// 将当前新上线玩家添加到worldManager中
	core.WorldMgrObj.AddPlayer(player)

	// Bind the property "pID" to the connection
	// 将该连接绑定属性PID
	conn.SetProperty("pID", player.PID)

	// Synchronize online player information and display surrounding player information
	// 同步周边玩家上线信息，与现实周边玩家信息
	player.SyncSurrounding()

	fmt.Println("=====> Player pIDID = ", player.PID, " arrived ====")
}

// OnConnectionLost Hook function called when a client disconnects
// 当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	// Get the "pID" property of the current connection
	// 获取当前连接的PID属性
	pID, _ := conn.GetProperty("pID")
	var playerID int32
	if pID != nil {
		playerID = pID.(int32)
	}

	// Get the corresponding player object based on the player ID
	// 根据pID获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPID(playerID)

	// Trigger the player's disconnection business logic
	// 触发玩家下线业务
	if player != nil {
		player.LostConnection()
	}

	fmt.Println("====> Player ", playerID, " left =====")

}

func main() {
	// Create a server instance
	// 创建服务器句柄
	s := znet.NewServer()

	// Register functions for client connection establishment and loss
	// 注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// Register routers
	s.AddRouter(2, &api.WorldChatApi{})
	s.AddRouter(3, &api.MoveApi{})

	// Add LTV data format Decoder
	s.SetDecoder(zdecoder.NewLTV_Little_Decoder())
	// Add LTV data format Pack packet Encoder
	s.SetPacket(zpack.NewDataPackLtv())

	// Start the server
	s.Serve()
}
