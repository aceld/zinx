package router

import (
	"bytes"
	"fmt"
	"github.com/aceld/zinx/zdecoder"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

type Data0x15Router struct {
	znet.BaseRouter
}

func (this *Data0x15Router) Handle(request ziface.IRequest) {
	fmt.Println("Data0x15Router Handle", request.GetMessage().GetData())
	_response := request.GetResponse()
	if _response != nil {
		switch _response.(type) {
		case zdecoder.HtlvCrcDecoder:
			_data := _response.(zdecoder.HtlvCrcDecoder)
			fmt.Println("Data0x15Router", _data)
			buffer := pack15(_data)
			request.GetConnection().Send(buffer)
		}
	}
}

// 头码   功能码 数据长度      Body                         CRC
// A2      10     0E        0102030405060708091011121314 050B
func pack15(_data zdecoder.HtlvCrcDecoder) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte(0xA1)
	buffer.WriteByte(_data.Funcode)
	buffer.WriteByte(0x26)
	//3~9：设备代码	将IMEI码转换为16进制
	buffer.Write(_data.Body[:7])
	//10：机型代码	A8（即热式幼儿园机）
	buffer.WriteByte(0xA8)
	//11：主机状态1	Bit0：0-待机中，1-运行中
	//Bit1：0-非智控，1-智控【本设备按智控】
	//Bit2：0-不能饮用，1-可以饮用
	//Bit3：0-无人用水，1-有人用水
	//Bit4：0-上电进水中，1-正常工作中
	//Bit5：0-消毒未启动，1-消毒进行中
	//Bit6：0-低压开关断开（无水），1-低压开关接通（有水）
	//Bit7：0-主机不带RO，1-主机带RO
	buffer.WriteByte(0x01)
	//12：主机状态2	Bit0：0－RO机不允许启动水泵，1－RO机允许启动水泵
	//Bit1：0－开水无人用，1－开水有人用
	//Bit2：0－高压开关断开（满水），1－高压开关接通（缺水）
	//Bit3：0－冰水无人用，1－冰水有人用【本设备无意义】
	//Bit4：0－无漏水信号，1－有漏水信号
	//Bit5：0－紫外灯未启动，1－紫外线灯杀菌中
	//Bit6：预留
	//Bit7：预留
	buffer.WriteByte(0x01)
	//13：水位状态
	//（即热水位）	Bit0：低水位，0-代表无水，1-代表有水【本设备低水位有水即表示水满】
	//Bit1：中水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit2：高水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit3：溢出水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit4：预留
	//Bit5：预留
	//Bit6：预留
	//Bit7：预留
	buffer.WriteByte(0x01)
	//14：开水温度	0℃~100℃，表示当前开水温度
	buffer.WriteByte(0x1A)
	//15：当前系统的停止加热温度	30~98℃，实际数值
	buffer.WriteByte(0x27)
	//16：负载状态	Bit0：加热，0－未加热，1－加热中
	//Bit1：进水，0－未进水，1－进水中
	//Bit2：换水或消毒，0－未换水，1－换水或消毒
	//Bit3：冲洗，0－未冲洗，1－冲洗中
	//Bit4：增压泵和RO进水阀，0－未启动，1－启动增压泵和RO进水阀
	//Bit5：RO进水阀2，0－未启动，1－启动中【本设备无意义】
	//Bit6：开水出水阀1，0－未启动，1－启动中
	//Bit7：净化水出水阀1，0－未启动，1－启动中【本设备无意义】
	buffer.WriteByte(0x01)
	//17：负载状态2	预留，填0x00
	buffer.WriteByte(0x00)
	//18：故障状态	Bit0：故障1，0－无故障，1－有故障
	//Bit1：故障2，0－无故障，1－有故障
	//Bit2：障保3，0－无故障，1－有故障
	//Bit3：故障4 ，0－无故障，1－有故障
	//Bit4：故障5 ，0－无故障，1－有故障
	//Bit5：故障6，0－无故障，1－有故障
	//Bit6：故障7，0－无故障，1－有故障
	//Bit7：故障8，0－无故障，1－有故障
	buffer.WriteByte(0x00)
	//19：故障状态2	Bit0：故障A，0－无故障，1－有故障
	//Bit1：故障B，0－无故障，1－有故障
	//Bit2：障保C，0－无故障，1－有故障
	//Bit3：故障9，0－无故障，1－有故障
	//Bit4：故障D，0－无故障，1－有故障
	//Bit5：故障E，0－无故障，1－有故障
	//Bit6：预留
	//Bit7：预留
	buffer.WriteByte(0x00)
	//20：主板软件版本	实际数值1~255
	buffer.WriteByte(0x01)
	//21：水位状态2	Bit0：纯水箱低水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit1：纯水箱中水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit2：纯水箱高水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit3：保温箱低水位，0-代表无水，1-代表有水
	//Bit4：保温箱中水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit5：保温箱高水位，0-代表无水，1-代表有水
	//Bit6：保温箱溢水位，0-代表无水，1-代表有水
	//Bit7：预留
	buffer.WriteByte(0x01)
	//22：温开水温度	0~100℃
	buffer.WriteByte(0x30)
	//23~24：剩余滤芯寿命	单位：小时，实际数值
	buffer.Write([]byte{23, 24})
	//25~26：剩余紫外线灯寿命
	buffer.Write([]byte{25, 26})
	//27~28：源水TDS值	0x0000－无此功能
	//实际数值，单位，ppm
	buffer.Write([]byte{27, 28})
	//29：净水TDS值	0x00－无此功能
	//实际数值，单位，ppm
	buffer.WriteByte(0x00)
	//30~33：耗电量	0xFFFFFFFF－无此功能
	//实际数值，高位在前，低位在后，单位wh
	buffer.Write([]byte{30, 31, 32, 33})
	//34：信号强度	0x01~0x28
	//0x01~0x0A对应:-81~-90dbm=极差
	//0x0B~0x14对应：-71~-80dbm=差
	//0x15~0x1E对应-61~-70dbm=好
	//0x1F~0x28对应：-41以上~-50dbm=良好
	buffer.WriteByte(0x30)
	//35~40：预留	全填0x00
	buffer.Write([]byte{0x00, 0x00})
	crc := zdecoder.GetCrC(buffer.Bytes())
	buffer.Write(crc)
	return buffer.Bytes()

}
