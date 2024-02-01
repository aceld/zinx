package api

import (
	"fmt"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zinx_app_demo/mmo_game/core"
	"github.com/aceld/zinx/zinx_app_demo/mmo_game/pb"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
)

// WorldChatApi World chat route business
// (WorldChatApi 世界聊天 路由业务)
type WorldChatApi struct {
	znet.BaseRouter
}

func (*WorldChatApi) Handle(request ziface.IRequest) {
	// 1. Decode the incoming proto protocol from the client
	// (1. 将客户端传来的proto协议解码)
	msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Talk Unmarshal error ", err)
		return
	}

	// 2. Identify which player sent the current message, retrieve from the connection property pID
	// (2. 得知当前的消息是从哪个玩家传递来的,从连接属性pID中获取)
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}

	// 3. Get the player object based on pID
	// (3. 根据pID得到player对象)
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))

	// 4. Have the player object initiate the chat broadcast request
	// (4. 让player对象发起聊天广播请求)
	player.Talk(msg.Content)
}
