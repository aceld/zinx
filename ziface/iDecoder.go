package ziface

import "encoding/binary"

type IDecoder interface {
	Decode(buff []byte) [][]byte
}

type LengthField struct {
	//大小端排序
	//大端模式：是指数据的高字节保存在内存的低地址中，而数据的低字节保存在内存的高地址中，地址由小向大增加，而数据从高位往低位放；
	//小端模式：是指数据的高字节保存在内存的高地址中，而数据的低字节保存在内存的低地址中，高地址部分权值高，低地址部分权值低，和我们的日常逻辑方法一致。
	//不了解的自行查阅一下资料
	Order               binary.ByteOrder
	MaxFrameLength      int64 //最大帧长度
	LengthFieldOffset   int   //长度字段偏移量
	LengthFieldLength   int   //长度域字段的字节数
	LengthAdjustment    int   //长度调整
	InitialBytesToStrip int   //需要跳过的字节数
}
