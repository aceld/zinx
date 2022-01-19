package znet

//Message 消息
type Message struct {
	MagicCode uint16 // 魔数（小安必须）
	ID        uint8  // 消息的ID（命令字，小安必须）
	SerialSn  uint8  // 序列号（小安必须）
	DataLen   uint16 // 消息的长度（小安必须）
	Data      []byte // 消息的内容
}

//NewMsgPackage 创建一个Message消息包
func NewMsgPackage(ID uint8, data []byte) *Message {
	return &Message{
		DataLen: uint16(len(data)),
		ID:      ID,
		Data:    data,
	}
}

//GetDataLen 获取消息数据段长度
func (msg *Message) GetDataLen() uint16 {
	return msg.DataLen
}

// GetSerialSn 获取头部序列号
func (msg *Message) GetSerialSn() uint8 {
	return msg.SerialSn
}

//GetMsgID 获取消息ID
func (msg *Message) GetMsgID() uint8 {
	return msg.ID
}

//GetData 获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

//SetDataLen 设置消息数据段长度
func (msg *Message) SetDataLen(len uint16) {
	msg.DataLen = len
}

//SetMsgID 设计消息ID
func (msg *Message) SetMsgID(msgID uint8) {
	msg.ID = msgID
}

//SetData 设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
