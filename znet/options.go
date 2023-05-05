package znet

import "github.com/aceld/zinx/ziface"

// Options for Server
// (Server的服务Option)
type Option func(s *Server)

// Implement custom data packet format by implementing the Packet interface,
// otherwise use the default data packet format
// (只要实现Packet 接口可自由实现数据包解析格式，如果没有则使用默认解析格式)
func WithPacket(pack ziface.IDataPack) Option {
	return func(s *Server) {
		s.SetPacket(pack)
	}
}

// Options for Client
type ClientOption func(c ziface.IClient)

// Implement custom data packet format by implementing the Packet interface for client,
// otherwise use the default data packet format
func WithPacketClient(pack ziface.IDataPack) ClientOption {
	return func(c ziface.IClient) {
		c.SetPacket(pack)
	}
}

// Set client name
func WithNameClient(name string) ClientOption {
	return func(c ziface.IClient) {
		c.SetName(name)
	}
}
