package router

import (
	"bytes"
	"fmt"
	"github.com/aceld/zinx/zdecoder"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type Data0x13Router struct {
	znet.BaseRouter
}

func (this *Data0x13Router) Handle(request ziface.IRequest) {
	fmt.Println("Data0x13Router Handle", request.GetMessage().GetData())
	_response := request.GetResponse()
	if _response != nil {
		switch _response.(type) {
		case zdecoder.HtlvCrcDecoder:
			_data := _response.(zdecoder.HtlvCrcDecoder)
			fmt.Println("Data0x13Router", _data)
			buffer := pack13(_data)
			request.GetConnection().Send(buffer)
		}
	}
}

// 头码   功能码 数据长度      Body                         CRC
// A2      10     0E        0102030405060708091011121314 050B
func pack13(_data zdecoder.HtlvCrcDecoder) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte(0xA1)
	buffer.WriteByte(_data.Funcode)
	buffer.WriteByte(0x0E)
	//3~9:3~6：用户卡号	用户IC卡卡号
	buffer.Write(_data.Body[:4])
	//7：卡状态：	0x00-未绑定（如服务器未查询到该IC卡时）
	//0x01-已绑定
	//0x02-解除绑定（如服务器查询到该IC卡解除绑定时下发）
	buffer.WriteByte(0x01)
	//8~9：剩余使用天数	该用户的剩余流量天数
	buffer.Write([]byte{8, 9})
	//10~11：每次最大出水量	单位mL，实际出水量
	buffer.Write([]byte{10, 11})
	//12~16：预留	全填0x00
	buffer.Write([]byte{12, 13, 14, 15, 16})
	crc := zdecoder.GetCrC(buffer.Bytes())
	buffer.Write(crc)
	return buffer.Bytes()

}
