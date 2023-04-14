// LTV，即Length-Tag(Type)-Value，是一种简单实用的数据传输方案。
//在LTV的定义中，可以知道它包括三个域，分别为：标签域（Tag），长度域（Length），内容域（Value）。这里的长度域的值实际上就是内容域的长度。
//
//解码前 (20 bytes)                                   解码后 (20 bytes)
//+------------+------------+-----------------+      +------------+------------+-----------------+
//|   Length   |     Tag    |     Value       |----->|  Length    |     Tag    |     Value       |
//| 0x0000000C | 0x00000001 | "HELLO, WORLD"  |      | 0x0000000C | 0x00000001 | "HELLO, WORLD"  |
//+------------+------------+-----------------+      +------------+------------+-----------------+
// Length：uint32类型，占4字节，Length标记Value长度12(hex:0x0000000C)
// Tag：   uint32类型，占4字节，Tag作为MsgId，暂定为1
// Value： 共12个字符，占12字节
//
//   说明：
//   lengthFieldOffset   = 0            (Length的字节位索引下标是0) 长度字段的偏差
//   lengthFieldLength   = 4            (Length是4个byte) 长度字段占的字节数
//   lengthAdjustment    = 4            (Length只表示Value长度，程序只会读取Length个字节就结束，后面没有来，故为0，若Value后面还有crc占2字节的话，那么此处就是2。若Length标记的是Tag+Length+Value总长度，那么此处是-8)
//   initialBytesToStrip = 0            (这个0表示返回完整的协议内容Tag+Length+Value，如果只想返回Value内容，去掉Tag的4字节和Length的4字节，此处就是8) 从解码帧中第一次去除的字节数
//   maxFrameLength      = 2^32 + 4 + 4 (Length为uint类型，故2^32次方表示Value最大长度，此外Tag和Length各占4字节)

package zdecoder

import (
	"bytes"
	"encoding/binary"
	"github.com/aceld/zinx/ziface"
	"math"
)

const LTV_HEADER_SIZE = 8 //表示TLV空包长度

type LTV_Little_Decoder struct {
	Length uint32 //消息长度
	Tag    uint32 //消息类型
	Value  []byte //消息内容
}

func NewLTV_Little_Decoder() ziface.IDecoder {
	return &LTV_Little_Decoder{}
}

func (ltv *LTV_Little_Decoder) GetLengthField() *ziface.LengthField {
	// +---------------+---------------+---------------+
	// |    Length     |     Tag       |     Value     |
	// | uint32(4byte) | uint32(4byte) |     n byte    |
	// +---------------+---------------+---------------+
	// Length：uint32类型，占4字节，Length标记Value长度
	// Tag：   uint32类型，占4字节
	// Value： 占n字节
	//
	//说明:
	//    lengthFieldOffset   = 0            (Length的字节位索引下标是4) 长度字段的偏差
	//    lengthFieldLength   = 4            (Length是4个byte) 长度字段占的字节数
	//    lengthAdjustment    = 4            (Length只表示Value长度，程序只会读取Length个字节就结束，后面没有来，故为0，若Value后面还有crc占2字节的话，那么此处就是2。若Length标记的是Tag+Length+Value总长度，那么此处是-8)
	//    initialBytesToStrip = 0            (这个0表示返回完整的协议内容Tag+Length+Value，如果只想返回Value内容，去掉Tag的4字节和Length的4字节，此处就是8) 从解码帧中第一次去除的字节数
	//    maxFrameLength      = 2^32 + 4 + 4 (Length为uint32类型，故2^32次方表示Value最大长度，此外Tag和Length各占4字节)
	//默认使用TLV封包方式
	return &ziface.LengthField{
		MaxFrameLength:      math.MaxUint32 + 4 + 4,
		LengthFieldOffset:   0,
		LengthFieldLength:   4,
		LengthAdjustment:    4,
		InitialBytesToStrip: 0,
	}
}

func (ltv *LTV_Little_Decoder) decode(data []byte) *LTV_Little_Decoder {
	ltvData := LTV_Little_Decoder{}

	//获取L
	ltvData.Length = binary.LittleEndian.Uint32(data[0:4])
	//获取T
	ltvData.Tag = binary.LittleEndian.Uint32(data[4:8])
	//确定V的长度
	ltvData.Value = make([]byte, ltvData.Length)

	//5. 获取V
	binary.Read(bytes.NewBuffer(data[8:8+ltvData.Length]), binary.LittleEndian, ltvData.Value)

	return &ltvData
}

func (ltv *LTV_Little_Decoder) Intercept(chain ziface.IChain) ziface.IcResp {
	//1.获取zinx的IMessage
	iMessage := chain.GetIMessage()
	if iMessage == nil {
		//进入责任链下一层
		return chain.ProceedWithIMessage(iMessage, nil)
	}

	//2. 获取数据
	data := iMessage.GetData()
	//zlog.Ins().DebugF("LTV-RawData size:%d data:%s\n", len(data), hex.EncodeToString(data))

	//3. 读取的数据不超过包头，直接进入下一层
	if len(data) < LTV_HEADER_SIZE {
		//进入责任链下一层
		return chain.ProceedWithIMessage(iMessage, nil)
	}

	//4. ltv解码
	ltvData := ltv.decode(data)

	//5. 将解码后的数据重新设置到IMessage中, Zinx的Router需要MsgID来寻址
	iMessage.SetDataLen(ltvData.Length)
	iMessage.SetMsgID(ltvData.Tag)
	iMessage.SetData(ltvData.Value)

	//6. 将解码后的数据进入下一层
	return chain.ProceedWithIMessage(iMessage, *ltvData)
}
