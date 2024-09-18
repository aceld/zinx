// Package ziface 主要提供zinx全部抽象层接口定义.
// 包括:
//
//			IServer 服务mod接口
//			IRouter 路由mod接口
//			IConnection 连接mod层接口
//	     IMessage 消息mod接口
//			IDataPack 消息拆解接口
//	     IMsgHandler 消息处理及协程池接口
//
// 当前文件描述:
// @Title  imessage.go
// @Description  提供消息的基本方法
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

/*
	将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessage interface {
	GetMagicCode() uint16 // 获取魔数
	GetDataLen() uint16   // 获取消息数据段长度
	GetMsgID() uint8      // 获取消息ID（命令字）
	GetSn() uint8         // 获取序列号
	GetData() []byte      // 获取消息内容

	SetMagicCode(uint16) // 设置魔数
	SetMsgID(uint8)      // 设置消息ID
	SetData([]byte)      // 设置消息内容
	SetDataLen(uint16)   // 设置消息数据段长度
}
