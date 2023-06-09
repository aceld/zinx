// TLV, which stands for Tag(Type)-Length-Value, is a simple and practical data transmission scheme.
// In the definition of TLV, it can be seen that it consists of three fields: the tag field (Tag), the length
// field (Length), and the value field (Value).
// The value of the length field is actually the length of the content field.
//
// Before decoding (20 bytes) After decoding (20 bytes)
// +------------+------------+-----------------+       +------------+------------+-----------------+
// | Tag        | Length     | Value           |-----> | Tag        | Length     | Value           |
// | 0x00000001 | 0x0000000C | "HELLO, WORLD"  |       | 0x00000001 | 0x0000000C | "HELLO, WORLD"  |
// +------------+------------+-----------------+       +------------+------------+-----------------+
// Tag: uint32 type, occupies 4 bytes, Tag is set as MsgId, temporarily set to 1
// Length: uint32 type, occupies 4 bytes, Length marks the length of Value, which is 12(hex:0x0000000C)
// Value: 12 characters in total, occupies 12 bytes
//
// Explanation:
// lengthFieldOffset = 4 (The byte index of Length is 4) Length field offset
// lengthFieldLength = 4 (Length is 4 bytes) Length field length in bytes
// lengthAdjustment = 0 (Length only represents the length of Value. The program will read only Length bytes and end.
// 	                     If there is a crc of 2 bytes after Value, then this is 2. If Length marks the total length of
//	                     Tag+Length+Value, then this is -8)
// initialBytesToStrip = 0 (This 0 means that the complete protocol content Tag+Length+Value is returned. If you only
//                          want to return the Value content, remove the 4 bytes of Tag and 4 bytes of Length, and
//                          this is 8) Number of bytes to strip from the decoded frame
// maxFrameLength = 2^32 + 4 + 4 (Since Length is of uint type, 2^32 represents the maximum length of Value. In addition,
//                                Tag and Length each occupy 4 bytes.)

// [简体中文]
// TLV，即Tag(Type)—Length—Value，是一种简单实用的数据传输方案。
//在TLV的定义中，可以知道它包括三个域，分别为：标签域（Tag），长度域（Length），内容域（Value）。这里的长度域的值实际上就是内容域的长度。
//
//解码前 (20 bytes)                                   解码后 (20 bytes)
//+------------+------------+-----------------+      +------------+------------+-----------------+
//|     Tag    |   Length   |     Value       |----->|     Tag    |   Length   |     Value       |
//| 0x00000001 | 0x0000000C | "HELLO, WORLD"  |      | 0x00000001 | 0x0000000C | "HELLO, WORLD"  |
//+------------+------------+-----------------+      +------------+------------+-----------------+
// Tag：   uint32类型，占4字节，Tag作为MsgId，暂定为1
// Length：uint32类型，占4字节，Length标记Value长度12(hex:0x0000000C)
// Value： 共12个字符，占12字节
//
//   说明：
//   lengthFieldOffset   = 4            (Length的字节位索引下标是4) 长度字段的偏差
//   lengthFieldLength   = 4            (Length是4个byte) 长度字段占的字节数
//   lengthAdjustment    = 0            (Length只表示Value长度，程序只会读取Length个字节就结束，后面没有来，故为0，
//                                       若Value后面还有crc占2字节的话，那么此处就是2。若Length标记的是Tag+Length+Value总长度，
//                                       那么此处是-8)
//   initialBytesToStrip = 0            (这个0表示返回完整的协议内容Tag+Length+Value，如果只想返回Value内容，去掉Tag的4字节和Length
//                                       的4字节，此处就是8) 从解码帧中第一次去除的字节数
//   maxFrameLength      = 2^32 + 4 + 4 (Length为uint类型，故2^32次方表示Value最大长度，此外Tag和Length各占4字节)

package zdecoder

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/aceld/zinx/ziface"
)

const TLV_HEADER_SIZE = 8 //表示TLV空包长度

type TLVDecoder struct {
	Tag    uint32 //T
	Length uint32 //L
	Value  []byte //V
}

func NewTLVDecoder() ziface.IDecoder {
	return &TLVDecoder{}
}

func (tlv *TLVDecoder) GetLengthField() *ziface.LengthField {
	// +---------------+---------------+---------------+
	// |    Tag        |     Length    |     Value     |
	// | uint32(4byte) | uint32(4byte) |     n byte    |
	// +---------------+---------------+---------------+
	// Length：uint32类型，占4字节，Length标记Value长度
	// Tag：   uint32类型，占4字节
	// Value： 占n字节
	//
	//说明:
	//    lengthFieldOffset   = 4            (Length的字节位索引下标是4) 长度字段的偏差
	//    lengthFieldLength   = 4            (Length是4个byte) 长度字段占的字节数
	//    lengthAdjustment    = 0            (Length只表示Value长度，程序只会读取Length个字节就结束，后面没有来，故为0，若Value后面还有crc占2字节的话，那么此处就是2。若Length标记的是Tag+Length+Value总长度，那么此处是-8)
	//    initialBytesToStrip = 0            (这个0表示返回完整的协议内容Tag+Length+Value，如果只想返回Value内容，去掉Tag的4字节和Length的4字节，此处就是8) 从解码帧中第一次去除的字节数
	//    maxFrameLength      = 2^32 + 4 + 4 (Length为uint32类型，故2^32次方表示Value最大长度，此外Tag和Length各占4字节)
	//默认使用TLV封包方式
	return &ziface.LengthField{
		MaxFrameLength:      math.MaxUint32 + 4 + 4,
		LengthFieldOffset:   4,
		LengthFieldLength:   4,
		LengthAdjustment:    0,
		InitialBytesToStrip: 0,
	}
}

func (tlv *TLVDecoder) decode(data []byte) *TLVDecoder {
	tlvData := TLVDecoder{}
	//Get T
	tlvData.Tag = binary.BigEndian.Uint32(data[0:4])
	//Get L
	tlvData.Length = binary.BigEndian.Uint32(data[4:8])
	//Determine the length of V. (确定V的长度)
	tlvData.Value = make([]byte, tlvData.Length)

	//Get V
	binary.Read(bytes.NewBuffer(data[8:8+tlvData.Length]), binary.BigEndian, tlvData.Value)

	//zlog.Ins().DebugF("TLV-DecodeData size:%d data:%+v\n", unsafe.Sizeof(data), tlvData)
	return &tlvData
}

func (tlv *TLVDecoder) Intercept(chain ziface.IChain) ziface.IcResp {

	//1. Get the IMessage of zinx
	iMessage := chain.GetIMessage()
	if iMessage == nil {
		// Go to the next layer in the chain of responsibility
		return chain.ProceedWithIMessage(iMessage, nil)
	}

	//2. Get Data
	data := iMessage.GetData()
	//zlog.Ins().DebugF("TLV-RawData size:%d data:%s\n", len(data), hex.EncodeToString(data))

	//3. If the amount of data read is less than the length of the header, proceed to the next layer directly.
	// (读取的数据不超过包头，直接进入下一层)
	if len(data) < TLV_HEADER_SIZE {
		return chain.ProceedWithIMessage(iMessage, nil)
	}

	//4. TLV Decode
	tlvData := tlv.decode(data)

	//5. Set the decoded data back to the IMessage, the Zinx Router needs MsgID for addressing
	// (将解码后的数据重新设置到IMessage中, Zinx的Router需要MsgID来寻址)
	iMessage.SetMsgID(tlvData.Tag)
	iMessage.SetData(tlvData.Value)
	iMessage.SetDataLen(tlvData.Length)

	//6. Pass the decoded data to the next layer.
	// (将解码后的数据进入下一层)
	return chain.ProceedWithIMessage(iMessage, *tlvData)
}
