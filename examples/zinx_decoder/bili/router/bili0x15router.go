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
// HeaderCode FunctionCode DataLength Body                         CRC
// A2         10           0E         0102030405060708091011121314 050B
func pack15(_data zdecoder.HtlvCrcDecoder) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte(0xA1)
	buffer.WriteByte(_data.Funcode)
	buffer.WriteByte(0x26)
	//3~9: Device Code Convert IMEI code to hex (3~9：设备代码, 将IMEI码转换为16进制)
	buffer.Write(_data.Body[:7])
	//10: Model Code A8 (Hot Type Kindergarten Machine) , 10: 机型代码	A8（即热式幼儿园机）
	buffer.WriteByte(0xA8)
	//11：主机状态1
	//Bit0：0-待机中，1-运行中
	//Bit1：0-非智控，1-智控【本设备按智控】
	//Bit2：0-不能饮用，1-可以饮用
	//Bit3：0-无人用水，1-有人用水
	//Bit4：0-上电进水中，1-正常工作中
	//Bit5：0-消毒未启动，1-消毒进行中
	//Bit6：0-低压开关断开（无水），1-低压开关接通（有水）
	//Bit7：0-主机不带RO，1-主机带RO
	// 11: Host Status 1
	// Bit0: 0-Standby, 1-Running
	// Bit1: 0-Non-intelligent control, 1-Intelligent control [This device follows intelligent control]
	// Bit2: 0-Cannot be drunk, 1-Can be drunk
	// Bit3: 0-No one uses water, 1-Someone uses water
	// Bit4: 0-Entering water at power-on, 1-Working normally
	// Bit5: 0-Disinfection not started, 1-Disinfection in progress
	// Bit6: 0-Low voltage switch off (no water), 1-Low voltage switch on (with water)
	// Bit7: 0-Master without RO, 1-Master with RO
	buffer.WriteByte(0x01)
	//12：主机状态2	Bit0：0－RO机不允许启动水泵，1－RO机允许启动水泵
	//Bit1：0－开水无人用，1－开水有人用
	//Bit2：0－高压开关断开（满水），1－高压开关接通（缺水）
	//Bit3：0－冰水无人用，1－冰水有人用【本设备无意义】
	//Bit4：0－无漏水信号，1－有漏水信号
	//Bit5：0－紫外灯未启动，1－紫外线灯杀菌中
	//Bit6：预留
	//Bit7：预留
	// 12: Host Status 2
	// Bit0: 0-RO machine does not allow pump to start, 1-RO machine allows pump to start
	// Bit1: 0-No one uses hot water, 1-Someone uses hot water
	// Bit2: 0-High-voltage switch off (full water), 1-High-voltage switch on (lack of water)
	// Bit3: 0-No one uses ice water, 1-Someone uses ice water [meaningless for this device]
	// Bit4: 0-No water leakage signal, 1-Water leakage signal
	// Bit5: 0-Ultraviolet lamp not started, 1-Ultraviolet lamp sterilizing
	// Bit6: Reserved
	// Bit7: Reserved
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
	// 13: Water level status (i.e. hot water level)
	// Bit0: Low water level, 0-represents no water, 1-represents water [Low water level with water indicates full water for this device]
	// Bit1: Medium water level, 0-represents no water, 1-represents water [Meaningless for this device]
	// Bit2: High water level, 0-represents no water, 1-represents water [Meaningless for this device]
	// Bit3: Overflow water level, 0-represents no water, 1-represents water [Meaningless for this device]
	// Bit4: Reserved
	// Bit5: Reserved
	// Bit6: Reserved
	// Bit7: Reserved
	buffer.WriteByte(0x01)
	buffer.WriteByte(0x01)
	// 14: Hot water temperature 0℃100℃, indicating the current hot water temperature
	//14：开水温度	0℃~100℃，表示当前开水温度
	buffer.WriteByte(0x1A)
	// 15: Stop heating temperature of current system 30~98℃, actual value
	// 15：当前系统的停止加热温度	30~98℃，实际数值
	buffer.WriteByte(0x27)
	//16：负载状态
	//Bit0：加热，0－未加热，1－加热中
	//Bit1：进水，0－未进水，1－进水中
	//Bit2：换水或消毒，0－未换水，1－换水或消毒
	//Bit3：冲洗，0－未冲洗，1－冲洗中
	//Bit4：增压泵和RO进水阀，0－未启动，1－启动增压泵和RO进水阀
	//Bit5：RO进水阀2，0－未启动，1－启动中【本设备无意义】
	//Bit6：开水出水阀1，0－未启动，1－启动中
	//Bit7：净化水出水阀1，0－未启动，1－启动中【本设备无意义】
	// 16: Load status
	// Bit0: Heating, 0 - Not heating, 1 - Heating
	// Bit1: Inlet water, 0 - No inlet water, 1 - Inlet water
	// Bit2: Water change or disinfection, 0 - No change, 1 - Water change or disinfection
	// Bit3: Flushing, 0 - Not flushing, 1 - Flushing
	// Bit4: Booster pump and RO inlet valve, 0 - Not started, 1 - Started booster pump and RO inlet valve
	// Bit5: RO inlet valve 2, 0 - Not started, 1 - Started 【Irrelevant to this device】
	// Bit6: Hot water outlet valve 1, 0 - Not started, 1 - Started
	// Bit7: Purified water outlet valve 1, 0 - Not started, 1 - Started 【Irrelevant to this device】
	buffer.WriteByte(0x01)
	//17：负载状态2	预留，填0x00
	//17: Load State 2, reserved, fill with 0x00.
	buffer.WriteByte(0x00)
	//18：故障状态
	//Bit0：故障1，0－无故障，1－有故障
	//Bit1：故障2，0－无故障，1－有故障
	//Bit2：障保3，0－无故障，1－有故障
	//Bit3：故障4 ，0－无故障，1－有故障
	//Bit4：故障5 ，0－无故障，1－有故障
	//Bit5：故障6，0－无故障，1－有故障
	//Bit6：故障7，0－无故障，1－有故障
	//Bit7：故障8，0－无故障，1－有故障
	//18: Fault status
	//Bit0: Fault 1, 0-no fault, 1-fault
	//Bit1: Fault 2, 0-no fault, 1-fault
	//Bit2: Fault 3, 0-no fault, 1-fault
	//Bit3: Fault 4, 0-no fault, 1-fault
	//Bit4: Fault 5, 0-no fault, 1-fault
	//Bit5: Fault 6, 0-no fault, 1-fault
	//Bit6: Fault 7, 0-no fault, 1-fault
	//Bit7: Fault 8, 0-no fault, 1-fault
	buffer.WriteByte(0x00)
	//19：故障状态2
	//Bit0：故障A，0－无故障，1－有故障
	//Bit1：故障B，0－无故障，1－有故障
	//Bit2：障保C，0－无故障，1－有故障
	//Bit3：故障9，0－无故障，1－有故障
	//Bit4：故障D，0－无故障，1－有故障
	//Bit5：故障E，0－无故障，1－有故障
	//Bit6：预留
	//Bit7：预留
	//19: Fault status 2
	//Bit0: Fault A, 0 - no fault, 1 - there is a fault
	//Bit1: Fault B, 0 - no fault, 1 - there is a fault
	//Bit2: Fault C, 0 - no fault, 1 - there is a fault
	//Bit3: Fault 9, 0 - no fault, 1 - there is a fault
	//Bit4: Fault D, 0 - no fault, 1 - there is a fault
	//Bit5: Fault E, 0 - no fault, 1 - there is a fault
	//Bit6: Reserved
	//Bit7: Reserved
	buffer.WriteByte(0x00)
	//20：主板软件版本	实际数值1~255
	//20: Mainboard software version, actual value 1~255.
	buffer.WriteByte(0x01)
	//21：水位状态2
	//Bit0：纯水箱低水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit1：纯水箱中水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit2：纯水箱高水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit3：保温箱低水位，0-代表无水，1-代表有水
	//Bit4：保温箱中水位，0-代表无水，1-代表有水【本设备无意义】
	//Bit5：保温箱高水位，0-代表无水，1-代表有水
	//Bit6：保温箱溢水位，0-代表无水，1-代表有水
	//Bit7：预留
	//21: Water level status 2
	//Bit0: Low water level in pure water tank, 0 - no water, 1 - water present (not applicable for this device)
	//Bit1: Medium water level in pure water tank, 0 - no water, 1 - water present (not applicable for this device)
	//Bit2: High water level in pure water tank, 0 - no water, 1 - water present (not applicable for this device)
	//Bit3: Low water level in insulation box, 0 - no water, 1 - water present
	//Bit4: Medium water level in insulation box, 0 - no water, 1 - water present (not applicable for this device)
	//Bit5: High water level in insulation box, 0 - no water, 1 - water present
	//Bit6: Overflow water level in insulation box, 0 - no water, 1 - water present
	//Bit7: Reserved
	buffer.WriteByte(0x01)
	//22：温开水温度	0~100℃
	//22: Hot water temperature in Celsius degree, range from 0 to 100℃.
	buffer.WriteByte(0x30)
	//23~24：剩余滤芯寿命	单位：小时，实际数值
	//23~24: Remaining filter life, unit: hours, actual numerical value
	buffer.Write([]byte{23, 24})
	//25~26：剩余紫外线灯寿命
	//25~26: Remaining UV lamp life, measured in hours, actual numerical value.
	buffer.Write([]byte{25, 26})
	//27~28：源水TDS值	0x0000－无此功能 ,实际数值，单位，ppm
	//27~28: Source water TDS value, 0x0000 - no such function, actual value in ppm
	buffer.Write([]byte{27, 28})
	//29：净水TDS值	0x00－无此功能, 实际数值，单位，ppm
	//29: TDS value of purified water 0x00 - No such function, actual value, unit: ppm
	buffer.WriteByte(0x00)
	//30~33：耗电量	0xFFFFFFFF－无此功能, 实际数值，高位在前，低位在后，单位wh
	//30~33: Power consumption, 0xFFFFFFFF - not supported, actual value, high byte first, low byte last, unit is Wh.
	buffer.Write([]byte{30, 31, 32, 33})
	//34：信号强度	0x01~0x28
	//0x01~0x0A对应:-81~-90dbm=极差
	//0x0B~0x14对应：-71~-80dbm=差
	//0x15~0x1E对应-61~-70dbm=好
	//0x1F~0x28对应：-41以上~-50dbm=良好
	//34: Signal strength 0x01~0x28
	//0x010x0A correspond to -81~-90dbm=poor
	//0x0B0x14 correspond to -71-80dbm=weak
	//0x150x1E correspond to -61-70dbm=good
	//0x1F0x28 correspond to -41 and above-50dbm=excellent
	buffer.WriteByte(0x30)
	//35~40：预留	全填0x00
	//35~40: Reserved. All filled with 0x00.
	buffer.Write([]byte{0x00, 0x00})
	crc := zdecoder.GetCrC(buffer.Bytes())
	buffer.Write(crc)
	return buffer.Bytes()

}
