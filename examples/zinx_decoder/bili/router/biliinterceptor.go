// HTLV+CRC，H头码，T功能码，L数据长度，V数据内容
// +------+-------+---------+--------+--------+
// | 头码  | 功能码 | 数据长度 | 数据内容 | CRC校验 |
// | 1字节 | 1字节  | 1字节   | N字节   |  2字节  |
// +------+-------+---------+--------+--------+

//    头码   功能码 数据长度      Body                         CRC
//    A2      10     0E        0102030405060708091011121314 050B
//
//
//    说明：
//    1.数据长度len是14(0E),这里的len仅仅指Body长度;
//
//
//    lengthFieldOffset   = 2   (len的索引下标是2，下标从0开始) 长度字段的偏差
//    lengthFieldLength   = 1   (len是1个byte) 长度字段占的字节数
//    lengthAdjustment    = 2   (len只表示Body长度，程序只会读取len个字节就结束，但是CRC还有2byte没读呢，所以为2)
//    initialBytesToStrip = 0   (这个0表示完整的协议内容，如果不想要A2，那么这里就是1) 从解码帧中第一次去除的字节数
//    maxFrameLength      = 255 + 4(起始码、功能码、CRC) (len是1个byte，所以最大长度是无符号1个byte的最大值)

package router

import (
	"encoding/hex"
	"fmt"
	"github.com/aceld/zinx/examples/zinx_decoder/bili/utils"
	"github.com/aceld/zinx/zcode"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"math"
)

const HEADER_SIZE = 5

type BiliData struct {
	data    []byte //数据内容
	head    byte   //头码
	funcode byte   //功能码
	length  byte   //数据长度
	body    []byte //数据内容
	crc     []byte //CRC校验
}

type HtlvcrcInterceptor struct {
}

func (this *HtlvcrcInterceptor) GetDecoder() ziface.Interceptor {
	return zcode.NewLengthFieldFrameInterceptor(math.MaxUint8+4, 2, 1, 2, 0)
}

func (this *HtlvcrcInterceptor) Intercept(chain ziface.Chain) ziface.Response {
	request := chain.Request()
	if request != nil {
		switch request.(type) {
		case ziface.IRequest:
			iRequest := request.(ziface.IRequest)
			iMessage := iRequest.GetMessage()
			if iMessage != nil {
				data := iMessage.GetData()
				zlog.Ins().DebugF("HtlvcrcInterceptor %s \n", hex.EncodeToString(data))
				datasize := len(data)
				htlvData := BiliData{
					data: data,
				}
				if datasize >= HEADER_SIZE {
					htlvData.head = data[0]
					htlvData.funcode = data[1]
					htlvData.length = data[2]
					htlvData.body = data[3 : datasize-2]
					htlvData.crc = data[datasize-2 : datasize]
					if !utils.CheckCRC(data[:datasize-2], htlvData.crc) {
						fmt.Println("crc校验失败", hex.EncodeToString(data), hex.EncodeToString(htlvData.crc))
						return nil
					}
					iMessage.SetMsgID(uint32(htlvData.funcode))
					iRequest.SetResponse(htlvData)
					//zlog.Ins().DebugF("2htlvData %s \n", hex.EncodeToString(htlvData.data))
				}
			}
		}
	}
	return chain.Proceed(chain.Request())
}