package ziface

import "encoding/binary"

type IFrameDecoder interface {
	Decode(buff []byte) [][]byte
}

// ILengthField Basic attributes possessed by ILengthField
// (具备的基础属性)
type LengthField struct {
	/*
		Note:
		   Big-endian: the most significant byte (the "big end") of a word is placed at the byte with the lowest address;
		   the rest of the bytes are placed in order of decreasing significance towards the byte with the highest address.
		   Little-endian: the least significant byte (the "little end") of a word is placed at the byte with the lowest address;
		   the rest of the bytes are placed in order of increasing significance towards the byte with the highest address.
		(大端模式：是指数据的高字节保存在内存的低地址中，而数据的低字节保存在内存的高地址中，地址由小向大增加，而数据从高位往低位放；
		小端模式：是指数据的高字节保存在内存的高地址中，而数据的低字节保存在内存的低地址中，高地址部分权值高，低地址部分权值低，和我们的日常逻辑方法一致。)
	*/
	Order               binary.ByteOrder //The byte order: BigEndian or LittleEndian(大小端)
	MaxFrameLength      uint64           //The maximum length of a frame(最大帧长度)
	LengthFieldOffset   int              //The offset of the length field(长度字段偏移量)
	LengthFieldLength   int              //The length of the length field in bytes(长度域字段的字节数)
	LengthAdjustment    int              //The length adjustment(长度调整)
	InitialBytesToStrip int              //The number of bytes to strip from the decoded frame(需要跳过的字节数)
}
