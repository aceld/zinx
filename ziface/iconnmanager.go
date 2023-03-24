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
// @Title  iconnmanager.go
// @Description    连接管理相关,包括添加、删除、通过一个连接ID获得连接对象，当前连接数量、清空全部连接等方法
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

/*
	连接管理抽象层
*/
type IConnManager interface {
	Add(IConnection)                                                       //添加链接
	Remove(IConnection)                                                    //删除连接
	Get(uint64) (IConnection, error)                                       //利用ConnID获取链接
	Len() int                                                              //获取当前连接
	ClearConn()                                                            //删除并停止所有链接
	GetAllConnID() []uint64                                                //获取所有连接ID
	Range(func(uint64, IConnection, interface{}) error, interface{}) error //遍历所有连接
}
