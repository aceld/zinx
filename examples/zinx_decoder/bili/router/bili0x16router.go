package router

import (
	"bytes"
	"fmt"
	"github.com/aceld/zinx/zdecoder"
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
		case zdecoder.HtlvCrcDecoder:
			_data := _response.(zdecoder.HtlvCrcDecoder)
			fmt.Println("Data0x16Router", _data)
			buffer := pack16(_data)
			request.GetConnection().Send(buffer)
		}
	}
}

// Pack a complete 0x16 protocol data
// Format:
// HeadCode FuncCode DataLen Body                         CRC
// A2       10       0E      0102030405060708091011121314 050B
// 头码      功能码    数据长度  Body                         CRC
// A2       10        0E     0102030405060708091011121314 050B
func pack16(_data zdecoder.HtlvCrcDecoder) []byte {
	_data.Data[0] = 0xA1
	buffer := bytes.NewBuffer(_data.Data[:len(_data.Data)-2])
	crc := zdecoder.GetCrC(buffer.Bytes())
	buffer.Write(crc)
	return buffer.Bytes()

}
