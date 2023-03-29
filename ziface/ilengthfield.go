package ziface

import "encoding/binary"

type IFrameDecoder interface {
	Decode(buff []byte) [][]byte
}

//ILengthField 具备的基础属性
type LengthField struct {
	/*
		Note:
		大端模式：是指数据的高字节保存在内存的低地址中，而数据的低字节保存在内存的高地址中，地址由小向大增加，而数据从高位往低位放；
		小端模式：是指数据的高字节保存在内存的高地址中，而数据的低字节保存在内存的低地址中，高地址部分权值高，低地址部分权值低，和我们的日常逻辑方法一致。
	*/
	Order               binary.ByteOrder //大小端
	MaxFrameLength      uint64           //最大帧长度
	LengthFieldOffset   int              //长度字段偏移量
	LengthFieldLength   int              //长度域字段的字节数
	LengthAdjustment    int              //长度调整
	InitialBytesToStrip int              //需要跳过的字节数
}
