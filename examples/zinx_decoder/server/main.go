// LengthFieldFrameDecoder是一个基于长度字段的解码器，比较难理解的解码器，它主要有5个核心的参数配置：
// maxFrameLength：     数据包最大长度
// lengthFieldOffset：  长度字段偏移量
// lengthFieldLength：  长度字段所占的字节数
// lengthAdjustment：   长度的调整值
// initialBytesToStrip：解码后跳过的字节数

// 案例分析，以下7种案例足以满足所有协议，只处理断粘包，并不能处理错包，包的完整性需要依靠协议自身定义CRC来校验
// >>>> 案例1：
// lengthFieldOffset  =0 长度字段从0开始
// lengthFieldLength  =2 长度字段本身占2个字节
// lengthAdjustment   =0 需要调整0字节
// initialBytesToStrip=0 解码后跳过0字节
//
// 解码前 (14 bytes)                 解码后 (14 bytes)
// +--------+----------------+      +--------+----------------+
// | Length | Actual Content |----->| Length | Actual Content |
// | 0x000C | "HELLO, WORLD" |      | 0x000C | "HELLO, WORLD" |
// +--------+----------------+      +--------+----------------+
// Length为0x000C，这个是十六进制，0x000C转化十进制就是14

// >>>> 案例2：
// lengthFieldOffset  =0 长度字段从0开始
// lengthFieldLength  =2 长度字段本身占2个字节
// lengthAdjustment   =0 需要调整0字节
// initialBytesToStrip=2 解码后跳过2字节
//
// 解码前 (14 bytes)                 解码后 (12 bytes)
// +--------+----------------+      +----------------+
// | Length | Actual Content |----->| Actual Content |
// | 0x000C | "HELLO, WORLD" |      | "HELLO, WORLD" |
// +--------+----------------+      +----------------+
// 这时initialBytesToStrip字段起作用了，在解码后会将前面的2字节跳过，所以解码后就只剩余了数据部分。

// >>>> 案例3：
// lengthFieldOffset  =0 长度字段从0开始
// lengthFieldLength  =2 长度字段本身占2个字节
// lengthAdjustment   =-2 需要调整 -2 字节
// initialBytesToStrip=0 解码后跳过2字节
//
// 解码前 (14 bytes)                 解码后 (14 bytes)
// +--------+----------------+      +--------+----------------+
// | Length | Actual Content |----->| Length | Actual Content |
// | 0x000E | "HELLO, WORLD" |      | 0x000E | "HELLO, WORLD" |
// +--------+----------------+      +--------+----------------+
// 这时lengthAdjustment起作用了，因为长度字段的值包含了长度字段本身的2字节，
// 如果要获取数据的字节数，需要加上lengthAdjustment的值，就是 14+（-2）=12，这样才算出来数据的长度。

// >>>> 案例4：
// lengthFieldOffset  =2 长度字段从第2个字节开始
// lengthFieldLength  =3 长度字段本身占3个字节
// lengthAdjustment   =0 需要调整0字节
// initialBytesToStrip=0 解码后跳过0字节
//
// 解码前 (17 bytes)                              解码后 (17 bytes)
// +----------+----------+----------------+      +----------+----------+----------------+
// | Header 1 |  Length  | Actual Content |----->| Header 1 |  Length  | Actual Content |
// |  0xCAFE  | 0x00000C | "HELLO, WORLD" |      |  0xCAFE  | 0x00000C | "HELLO, WORLD" |
// +----------+----------+----------------+      +----------+----------+----------------+
// 由于数据包最前面加了2个字节的Header，所以lengthFieldOffset为2，
// 说明长度字段是从第2个字节开始的。然后lengthFieldLength为3，说明长度字段本身占了3个字节。

// >>>> 案例5：
// lengthFieldOffset  =0 长度字段从第0个字节开始
// lengthFieldLength  =3 长度字段本身占3个字节
// lengthAdjustment   =2 需要调整2字节
// initialBytesToStrip=0 解码后跳过0字节
//
// 解码前 (17 bytes)                              解码后 (17 bytes)
// +----------+----------+----------------+      +----------+----------+----------------+
// |  Length  | Header 1 | Actual Content |----->|  Length  | Header 1 | Actual Content |
// | 0x00000C |  0xCAFE  | "HELLO, WORLD" |      | 0x00000C |  0xCAFE  | "HELLO, WORLD" |
// +----------+----------+----------------+      +----------+----------+----------------+
// lengthFieldOffset为0，所以长度字段从0字节开始。lengthFieldLength为3，长度总共占3字节。
// 因为长度字段后面还剩余14字节的总数据，但是长度字段的值为12，只表示了数据的长度，不包含头的长度，
// 所以lengthAdjustment为2，就是12+2=14，计算出Header+Content的总长度。

// >>>> 案例6：
// lengthFieldOffset  =1 长度字段从第1个字节开始
// lengthFieldLength  =2 长度字段本身占2个字节
// lengthAdjustment   =1 需要调整1字节
// initialBytesToStrip=3 解码后跳过3字节
//
// 解码前 (16 bytes)                               解码后 (13 bytes)
// +------+--------+------+----------------+      +------+----------------+
// | HDR1 | Length | HDR2 | Actual Content |----->| HDR2 | Actual Content |
// | 0xCA | 0x000C | 0xFE | "HELLO, WORLD" |      | 0xFE | "HELLO, WORLD" |
// +------+--------+------+----------------+      +------+----------------+
//这一次将Header分为了两个1字节的部分，lengthFieldOffset为1表示长度从第1个字节开始，lengthFieldLength为2表示长度字段占2个字节。
//因为长度字段的值为12，只表示了数据的长度，所以lengthAdjustment为1，12+1=13，
//表示Header的第二部分加上数据的总长度为13。因为initialBytesToStrip为3，所以解码后跳过前3个字节。

// >>>> 案例7：
// lengthFieldOffset  =1 长度字段从第1个字节开始
// lengthFieldLength  =2 长度字段本身占2个字节
// lengthAdjustment   =-3 需要调整 -3 字节
// initialBytesToStrip=3 解码后跳过3字节
//
// 解码前 (16 bytes)                               解码后 (13 bytes)
// +------+--------+------+----------------+      +------+----------------+
// | HDR1 | Length | HDR2 | Actual Content |----->| HDR2 | Actual Content |
// | 0xCA | 0x0010 | 0xFE | "HELLO, WORLD" |      | 0xFE | "HELLO, WORLD" |
// +------+--------+------+----------------+      +------+----------------+
// 这一次长度字段的值为16，表示包的总长度，所以lengthAdjustment为 -3 ，16+ (-3)=13，
// 表示Header的第二部分加数据部分的总长度为13字节。initialBytesToStrip为3，解码后跳过前3个字节。

package main

import (
	"github.com/aceld/zinx/examples/zinx_decoder/router"
	"github.com/aceld/zinx/examples/zinx_decoder/server/interceptor"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

// 创建连接的时候执行
func DoConnectionBegin(conn ziface.IConnection) {
	zlog.Ins().InfoF("DoConnecionBegin is Called ...")

	//设置两个链接属性，在连接创建之后
	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://www.kancloud.cn/@aceld")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		zlog.Error(err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn ziface.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		zlog.Ins().InfoF("Conn Property Name = %v", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		zlog.Ins().InfoF("Conn Property Home = %v", home)
	}

	zlog.Ins().InfoF("Conn is Lost")
}

func main() {
	//创建一个server句柄
	s := znet.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//_interceptor := interceptor.HtlvcrcInterceptor{} //HTLV+CRC 断粘包
	_interceptor := interceptor.TLVInterceptor{} //TLV 断粘包
	//配置路由
	s.AddRouter(0x00000001, &router.TLVRouter{})
	//s.AddRouter(0x10, &router.HTLVCRCRouter{})

	//断粘包解码器
	s.AddInterceptor(_interceptor.GetDecoder())
	//解码后数据处理器
	s.AddInterceptor(&_interceptor)
	//开启服务
	s.Serve()
}
