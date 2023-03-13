package router

import (
	"bytes"
	"fmt"
	"github.com/aceld/zinx/examples/zinx_decoder/bili/utils"
	"github.com/aceld/zinx/examples/zinx_decoder/decode"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type Data0x16Router struct {
	znet.BaseRouter
}

func (this *Data0x16Router) Handle(request ziface.IRequest) {
	fmt.Println("Data0x16Router Handle", request.GetMessage().GetData())
	_response := request.GetResponse()
	if _response != nil {
		switch _response.(type) {
		case decode.HtlvCrcData:
			_data := _response.(decode.HtlvCrcData)
			fmt.Println("Data0x16Router", _data)
			buffer := pack16(_data)
			request.GetConnection().Send(buffer)
		}
	}
}

// 头码   功能码 数据长度      Body                         CRC
// A2      10     0E        0102030405060708091011121314 050B
func pack16(_data decode.HtlvCrcData) []byte {
	_data.Data[0] = 0xA1
	buffer := bytes.NewBuffer(_data.Data[:len(_data.Data)-2])
	crc := utils.GetCrC(buffer.Bytes())
	buffer.Write(crc)
	return buffer.Bytes()

}
