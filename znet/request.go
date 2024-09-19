package znet

import "github.com/chnkenc/zinx-xiaoan/ziface"

// Request 请求
type Request struct {
	conn ziface.IConnection // 已经和客户端建立好的 链接
	msg  ziface.IMessage    // 客户端请求的数据
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetMessage 获取客户端请求数据
func (r *Request) GetMessage() ziface.IMessage {
	return r.msg
}

// GetMagicCode 获取消息魔数
func (r *Request) GetMagicCode() uint16 {
	return r.msg.GetMagicCode()
}

// GetExtendData 获取扩展数据
func (r *Request) GetExtendData() []byte {
	return r.msg.GetExtendData()
}

// GetHeaderData 获取消息头数据
func (r *Request) GetHeaderData() []byte {
	return r.msg.GetHeaderData()
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取请求的消息的ID
func (r *Request) GetMsgID() uint8 {
	return r.msg.GetMsgID()
}

// GetSn 获取请求的消息的序列号
func (r *Request) GetSn() uint8 {
	return r.msg.GetSn()
}
