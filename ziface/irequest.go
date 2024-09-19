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
// @Title  irequest.go
// @Description  提供连接请求全部接口声明
// @Author  Ken.C - Thu Jan 20 16:31:18 CST 2022
package ziface

/*
	IRequest 接口：
	实际上是把客户端请求的链接信息 和 请求的数据 包装到了 Request里
*/
type IRequest interface {
	GetConnection() IConnection // 获取请求连接信息
	GetMessage() IMessage       // 获取客户端请求数据
	GetMagicCode() uint16       // 获取消息魔数
	GetExtendData() []byte      // 获取扩展数据
	GetHeaderData() []byte      // 获取消息头数据
	GetData() []byte            // 获取请求消息的数据
	GetMsgID() uint8            // 获取请求的消息ID
	GetSn() uint8               // 获取请求的消息的序列号
}
