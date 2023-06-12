// @Title idatapack.go
// @Description Message packing and unpacking methods
// @Author Aceld - Thu Mar 11 10:32:29 CST 2019
package ziface

/*
IDataPack Package and unpack data.
Operating on the data stream of TCP connections, add header information to transfer data, and solve TCP sticky packets.
(封包数据和拆包数据
直接面向TCP连接中的数据流,为传输数据添加头部信息，用于处理TCP粘包问题。)
*/
type IDataPack interface {
	GetHeadLen() uint32                // Get the length of the message header(获取包头长度方法)
	Pack(msg IMessage) ([]byte, error) // Package message (封包方法)
	Unpack([]byte) (IMessage, error)   // Unpackage message(拆包方法)
}

const (
	// Zinx standard packing and unpacking method (Zinx 标准封包和拆包方式)
	ZinxDataPack    string = "zinx_pack_tlv_big_endian"
	ZinxDataPackOld string = "zinx_pack_ltv_little_endian"

	//...(+)
	//// Custom packing method can be added here(自定义封包方式在此添加)
)

const (
	// Zinx default standard message protocol format(Zinx 默认标准报文协议格式)
	ZinxMessage string = "zinx_message"
)
