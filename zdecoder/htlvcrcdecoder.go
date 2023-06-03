// HTLV+CRC, H header code, T function code, L data length, V data content
//+------+-------+-------+--------+--------+
//| H    | T     | L     | V      | CRC    |
//| 1Byte| 1Byte | 1Byte | NBytes | 2Bytes |
//+------+-------+-------+--------+--------+

// HeaderCode FunctionCode DataLength Body                         CRC
// A2         10           0E         0102030405060708091011121314 050B
//
//
// Explanation:
// 1. The data length len is 14 (0E), where len only refers to the length of the Body.
//
//
// lengthFieldOffset = 2 (the index of len is 2, starting from 0) The offset of the length field
// lengthFieldLength = 1 (len is 1 byte) The length of the length field in bytes
// lengthAdjustment = 2 (len only represents the length of the Body, the program will only read len bytes and end, but there are still 2 bytes of CRC to read, so it's 2)
// initialBytesToStrip = 0 (this 0 represents the complete protocol content. If you don't want A2, then it's 1) The number of bytes to strip from the decoding frame for the first time
// maxFrameLength = 255 + 4 (starting code, function code, CRC) (len is 1 byte, so the maximum length is the maximum value of an unsigned byte plus 4 bytes)

// [简体中文]
//
// HTLV+CRC，H头码，T功能码，L数据长度，V数据内容
//+------+-------+---------+--------+--------+
//| 头码  | 功能码 | 数据长度 | 数据内容 | CRC校验 |
//| 1字节 | 1字节  | 1字节   | N字节   |  2字节  |
//+------+-------+---------+--------+--------+

//头码   功能码 数据长度      Body                         CRC
//A2      10     0E        0102030405060708091011121314 050B
//
//
//   说明：
//   1.数据长度len是14(0E),这里的len仅仅指Body长度;
//
//
//   lengthFieldOffset   = 2   (len的索引下标是2，下标从0开始) 长度字段的偏差
//   lengthFieldLength   = 1   (len是1个byte) 长度字段占的字节数
//   lengthAdjustment    = 2   (len只表示Body长度，程序只会读取len个字节就结束，但是CRC还有2byte没读呢，所以为2)
//   initialBytesToStrip = 0   (这个0表示完整的协议内容，如果不想要A2，那么这里就是1) 从解码帧中第一次去除的字节数
//   maxFrameLength      = 255 + 4(起始码、功能码、CRC) (len是1个byte，所以最大长度是无符号1个byte的最大值)

package zdecoder

import (
	"encoding/hex"
	"math"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
)

const HEADER_SIZE = 5

type HtlvCrcDecoder struct {
	Head    byte   //HeaderCode(头码)
	Funcode byte   //FunctionCode(功能码)
	Length  byte   //DataLength(数据长度)
	Body    []byte //BodyData(数据内容)
	Crc     []byte //CRC校验
	Data    []byte //// Original data content(原始数据内容)
}

func NewHTLVCRCDecoder() ziface.IDecoder {
	return &HtlvCrcDecoder{}
}

func (hcd *HtlvCrcDecoder) GetLengthField() *ziface.LengthField {
	//+------+-------+---------+--------+--------+
	//| 头码  | 功能码 | 数据长度 | 数据内容 | CRC校验 |
	//| 1字节 | 1字节  | 1字节   | N字节   |  2字节  |
	//+------+-------+---------+--------+--------+
	//头码   功能码 数据长度      Body                         CRC
	//A2      10     0E        0102030405060708091011121314 050B
	//说明：
	//   1.数据长度len是14(0E),这里的len仅仅指Body长度;
	//
	//   lengthFieldOffset   = 2   (len的索引下标是2，下标从0开始) 长度字段的偏差
	//   lengthFieldLength   = 1   (len是1个byte) 长度字段占的字节数
	//   lengthAdjustment    = 2   (len只表示Body长度，程序只会读取len个字节就结束，但是CRC还有2byte没读呢，所以为2)
	//   initialBytesToStrip = 0   (这个0表示完整的协议内容，如果不想要A2，那么这里就是1) 从解码帧中第一次去除的字节数
	//   maxFrameLength      = 255 + 4(起始码、功能码、CRC) (len是1个byte，所以最大长度是无符号1个byte的最大值)
	return &ziface.LengthField{
		MaxFrameLength:      math.MaxInt8 + 4,
		LengthFieldOffset:   2,
		LengthFieldLength:   1,
		LengthAdjustment:    2,
		InitialBytesToStrip: 0,
	}
}

func (hcd *HtlvCrcDecoder) decode(data []byte) *HtlvCrcDecoder {
	datasize := len(data)

	htlvData := HtlvCrcDecoder{
		Data: data,
	}

	// Parse the header
	htlvData.Head = data[0]
	htlvData.Funcode = data[1]
	htlvData.Length = data[2]
	htlvData.Body = data[3 : datasize-2]
	htlvData.Crc = data[datasize-2 : datasize]

	// CRC
	if !CheckCRC(data[:datasize-2], htlvData.Crc) {
		zlog.Ins().DebugF("crc check error %s %s\n", hex.EncodeToString(data), hex.EncodeToString(htlvData.Crc))
		return nil
	}

	//zlog.Ins().DebugF("2htlvData %s \n", hex.EncodeToString(htlvData.data))
	//zlog.Ins().DebugF("HTLVCRC-DecodeData size:%d data:%+v\n", unsafe.Sizeof(htlvData), htlvData)

	return &htlvData
}

func (hcd *HtlvCrcDecoder) Intercept(chain ziface.IChain) ziface.IcResp {
	//1. Get the IMessage of zinx
	iMessage := chain.GetIMessage()
	if iMessage == nil {
		// Go to the next layer in the chain of responsibility
		return chain.ProceedWithIMessage(iMessage, nil)
	}

	//2. Get Data
	data := iMessage.GetData()
	//zlog.Ins().DebugF("HTLVCRC-RawData size:%d data:%s\n", len(data), hex.EncodeToString(data))

	//3. If the amount of data read is less than the length of the header, proceed to the next layer directly.
	// (读取的数据不超过包头，直接进入下一层)
	if len(data) < HEADER_SIZE {
		return chain.ProceedWithIMessage(iMessage, nil)
	}

	//4. HTLV+CRC Decode
	htlvData := hcd.decode(data)

	//5. Set the decoded data back to the IMessage, the Zinx Router needs MsgID for addressing
	// (将解码后的数据重新设置到IMessage中, Zinx的Router需要MsgID来寻址)
	iMessage.SetMsgID(uint32(htlvData.Funcode))

	//6. Pass the decoded data to the next layer.
	// (将解码后的数据进入下一层)
	return chain.ProceedWithIMessage(iMessage, *htlvData)
}
