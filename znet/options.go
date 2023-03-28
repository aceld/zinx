package znet

import "github.com/aceld/zinx/ziface"

//Server的服务Option
type Option func(s *Server)

// 只要实现Packet 接口可自由实现数据包解析格式，如果没有则使用默认解析格式
func WithPacket(pack ziface.IDataPack) Option {
	return func(s *Server) {
		s.SetPacket(pack)
	}
}

//Client的客户端Option
type ClientOption func(c ziface.IClient)

//Client的客户端Option
func WithPacketClient(pack ziface.IDataPack) ClientOption {
	return func(c ziface.IClient) {
		c.SetPacket(pack)
	}
}
