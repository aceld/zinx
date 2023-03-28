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
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"math"
	"unsafe"
)

const HEADER_SIZE = 5

type HtlvCrcDecoder struct {
	Head    byte   //头码
	Funcode byte   //功能码
	Length  byte   //数据长度
	Body    []byte //数据内容
	Crc     []byte //CRC校验
	Data    []byte //数据内容
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

func (hcd *HtlvCrcDecoder) Intercept(chain ziface.IChain) ziface.IcResp {
	request := chain.Request()
	if request == nil {
		return chain.Proceed(chain.Request())
	}

	switch request.(type) {
	case ziface.IRequest:
		iRequest := request.(ziface.IRequest)
		iMessage := iRequest.GetMessage()
		if iMessage == nil {
			break
		}

		data := iMessage.GetData()
		zlog.Ins().DebugF("HTLVCRC-RawData size:%d data:%s\n", len(data), hex.EncodeToString(data))
		datasize := len(data)

		htlvData := HtlvCrcDecoder{
			Data: data,
		}

		//数据大于消息头长度，进行解析
		if datasize >= HEADER_SIZE {
			//解析头
			htlvData.Head = data[0]
			htlvData.Funcode = data[1]
			htlvData.Length = data[2]
			htlvData.Body = data[3 : datasize-2]
			htlvData.Crc = data[datasize-2 : datasize]

			//CRC校验
			if !CheckCRC(data[:datasize-2], htlvData.Crc) {
				zlog.Ins().DebugF("crc校验失败 %s %s\n", hex.EncodeToString(data), hex.EncodeToString(htlvData.Crc))
				return nil
			}

			//设置ZinxMessage消息ID
			iMessage.SetMsgID(uint32(htlvData.Funcode))
			//设置ZinxMessage消息内容
			iRequest.SetResponse(htlvData)

			//zlog.Ins().DebugF("2htlvData %s \n", hex.EncodeToString(htlvData.data))
			zlog.Ins().DebugF("HTLVCRC-DecodeData size:%d data:%+v\n", unsafe.Sizeof(htlvData), htlvData)
		}
	}

	return chain.Proceed(chain.Request())
}
