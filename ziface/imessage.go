// @Title imessage.go
// @Description Provides basic methods for messages
// @Author Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

// IMessage Package ziface defines an abstract interface for encapsulating a request message into a message
type IMessage interface {
	GetDataLen() uint32 // Gets the length of the message data segment(获取消息数据段长度)
	GetMsgID() uint32   // Gets the ID of the message(获取消息ID)
	GetData() []byte    // Gets the content of the message(获取消息内容)
	GetRawData() []byte // Gets the raw data of the message(获取原始数据)

	SetMsgID(uint32)   // Sets the ID of the message(设计消息ID)
	SetData([]byte)    // Sets the content of the message(设计消息内容)
	SetDataLen(uint32) // Sets the length of the message data segment(设置消息数据段长度)
}
