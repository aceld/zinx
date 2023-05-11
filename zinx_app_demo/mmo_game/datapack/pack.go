package datapack

import (
	"fmt"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zpack"
	"github.com/golang/protobuf/proto"
)

// Replace here with your custom packer (if one exists)
// 将这里替换为你自定义的打包器（如果存在的话）
var DefaultPack = zpack.Factory().NewPack(ziface.ZinxDataPack)

// Serialize message
// 将消息序列化
func SerializeMsg2Bytes(msgId uint32, msgData proto.Message) []byte {

	// Marshal message
	// 将proto Message结构体序列化
	msg, err := proto.Marshal(msgData)
	if err != nil {
		fmt.Printf("SerializeMsg2Bytes, marshal failed, msgId:%v, msgData:%v, err: %v\n", msgId, msgData, err)
		return nil
	}

	// Pack data
	// 封包
	data, err := DefaultPack.Pack(zpack.NewMsgPackage(msgId, msg))
	if err != nil {
		fmt.Printf("SerializeMsg2Bytes, pack failed, msgId:%v, msg:%v, err: %v\n", msgId, msg, err)
		return nil
	}

	return data
}
