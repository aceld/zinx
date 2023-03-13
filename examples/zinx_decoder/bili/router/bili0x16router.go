package router

import (
	"bytes"
	"fmt"
	"github.com/aceld/zinx/examples/zinx_decoder/bili/utils"
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
		case BiliData:
			_data := _response.(BiliData)
			fmt.Println("Data0x16Router", _data)
			buffer := pack16(_data)
			request.GetConnection().Send(buffer)
		}
	}
}

// 头码   功能码 数据长度      Body                         CRC
// A2      10     0E        0102030405060708091011121314 050B
func pack16(_data BiliData) []byte {
	_data.data[0] = 0xA1
	buffer := bytes.NewBuffer(_data.data[:len(_data.data)-2])
	crc := utils.GetCrC(buffer.Bytes())
	buffer.Write(crc)
	return buffer.Bytes()

}
