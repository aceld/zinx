// Package ziface 主要提供zinx全部抽象层接口定义.
// 包括:
//		IServer 服务mod接口
//		IRouter 路由mod接口
//		IConnection 连接mod层接口
//      IMessage 消息mod接口
//		IDataPack 消息拆解接口
//      IMsgHandler 消息处理及协程池接口
//
// 当前文件描述:
// @Title  idatapack.go
// @Description  消息的打包和解包方法
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

/*
	封包数据和拆包数据
	直接面向TCP连接中的数据流,为传输数据添加头部信息，用于处理TCP粘包问题。
*/
type IDataPack interface {
	GetHeadLen() uint32                //获取包头长度方法
	Pack(msg IMessage) ([]byte, error) //封包方法
	Unpack([]byte) (IMessage, error)   //拆包方法
}
