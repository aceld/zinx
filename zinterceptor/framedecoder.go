/**
 * @author uuxia
 * @date 15:57 2023/3/10
 * @description 通用解码器
 **/

package zinterceptor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sync"

	"github.com/aceld/zinx/ziface"
)

// FrameDecoder
// A decoder that splits the received {@link ByteBuf}s dynamically by the
// value of the length field in the message.  It is particularly useful when you
// decode a binary message which has an integer header field that represents the
// length of the message body or the whole message.
//
// ziface.LengthField has many configuration parameters so
// that it can decode any message with a length field, which is often seen in
// proprietary client-server protocols. Here are some example that will give
// you the basic idea on which option does what.
//
// <I. 2 bytes length field at offset 0, do not strip header>
//
// The value of the length field in this example is <tt>12 (0x0C)</tt> which
// represents the length of "HELLO, WORLD".  By default, the decoder assumes
// that the length field represents the number of the bytes that follows the
// length field.  Therefore, it can be decoded with the simplistic parameter
// combination.
//
// LengthFieldOffset   = 0
// LengthFieldLength   = 2
// LengthAdjustment    = 0
// InitialBytesToStrip = 0 (= do not strip header)
//
// BEFORE DECODE (14 bytes)         AFTER DECODE (14 bytes)
// +--------+----------------+      +--------+----------------+
// | Length | Actual Content |----->| Length | Actual Content |
// | 0x000C | "HELLO, WORLD" |      | 0x000C | "HELLO, WORLD" |
// +--------+----------------+      +--------+----------------+
//
//
// <II. 2 bytes length field at offset 0, strip header>
//
// Because we can get the length of the content by calling
// {@link ByteBuf#readableBytes()}, you might want to strip the length
// field by specifying `InitialBytesToStrip`.  In this example, we
// specified `2`, that is same with the length of the length field, to
// strip the first two bytes.
//
// LengthFieldOffset   = 0
// LengthFieldLength   = 2
// LengthAdjustment    = 0
// InitialBytesToStrip = 2 (= the length of the Length field)
//
// BEFORE DECODE (14 bytes)         AFTER DECODE (12 bytes)
// +--------+----------------+      +----------------+
// | Length | Actual Content |----->| Actual Content |
// | 0x000C | "HELLO, WORLD" |      | "HELLO, WORLD" |
// +--------+----------------+      +----------------+
//
//  <III. 2 bytes length field at offset 0, do not strip header, the length field>
//
//	represents the length of the whole message</h3>
//
// In most cases, the length field represents the length of the message body
// only, as shown in the previous examples.  However, in some protocols, the
// length field represents the length of the whole message, including the
// message header.  In such a case, we specify a non-zero
// `LengthAdjustment`.  Because the length value in this example message
//
// is always greater than the body length by `2`, we specify `-2`
// as `LengthAdjustment` for compensation.
//
// LengthFieldOffset   =  0
// LengthFieldLength   =  2
// LengthAdjustment    = -2 (= the length of the Length field)
// InitialBytesToStrip =  0
//
// BEFORE DECODE (14 bytes)         AFTER DECODE (14 bytes)
// +--------+----------------+      +--------+----------------+
// | Length | Actual Content |----->| Length | Actual Content |
// | 0x000E | "HELLO, WORLD" |      | 0x000E | "HELLO, WORLD" |
// +--------+----------------+      +--------+----------------+
//
// <IV. 3 bytes length field at the end of 5 bytes header, do not strip header>
//
// The following message is a simple variation of the first example.  An extra
// header value is prepended to the message.  <tt>LengthAdjustment</tt> is zero
// again because the decoder always takes the length of the prepended data into
// account during frame length calculation.
//
// LengthFieldOffset   = 2 (= the length of Header 1)
// LengthFieldLength   = 3
// LengthAdjustment    = 0
// InitialBytesToStrip = 0
//
// BEFORE DECODE (17 bytes)                      AFTER DECODE (17 bytes)
// +----------+----------+----------------+      +----------+----------+----------------+
// | Header 1 |  Length  | Actual Content |----->| Header 1 |  Length  | Actual Content |
// |  0xCAFE  | 0x00000C | "HELLO, WORLD" |      |  0xCAFE  | 0x00000C | "HELLO, WORLD" |
// +----------+----------+----------------+      +----------+----------+----------------+
//
//
// <V. 3 bytes length field at the beginning of 5 bytes header, do not strip header>
//
// This is an advanced example that shows the case where there is an extra
// header between the length field and the message body.  You have to specify a
// positive `LengthAdjustment` so that the decoder counts the extra
// header into the frame length calculation.
//
// LengthFieldOffset   = 0
// LengthFieldLength   = 3
// LengthAdjustment    = 2 (= the length of Header 1)
// InitialBytesToStrip = 0
//
// BEFORE DECODE (17 bytes)                      AFTER DECODE (17 bytes)
// +----------+----------+----------------+      +----------+----------+----------------+
// |  Length  | Header 1 | Actual Content |----->|  Length  | Header 1 | Actual Content |
// | 0x00000C |  0xCAFE  | "HELLO, WORLD" |      | 0x00000C |  0xCAFE  | "HELLO, WORLD" |
// +----------+----------+----------------+      +----------+----------+----------------+
//
//
//  <VI. 2 bytes length field at offset 1 in the middle of 4 bytes header,
//	strip the first header field and the length field>
//
// This is a combination of all the examples above.  There are the prepended
// header before the length field and the extra header after the length field.
// The prepended header affects the	`LengthFieldOffset` and the extra
// header affects the `LengthAdjustment`.  We also specified a non-zero
// `InitialBytesToStrip` to strip the length field and the prepended
// header from the frame.  If you don't want to strip the prepended header, you
// could specify `0` for `initialBytesToSkip`.
//
// LengthFieldOffset   = 1 (= the length of HDR1)
// LengthFieldLength   = 2
// LengthAdjustment    = 1 (= the length of HDR2)
// InitialBytesToStrip = 3 (= the length of HDR1 + LEN)
//
// BEFORE DECODE (16 bytes)                       AFTER DECODE (13 bytes)
// +------+--------+------+----------------+      +------+----------------+
// | HDR1 | Length | HDR2 | Actual Content |----->| HDR2 | Actual Content |
// | 0xCA | 0x000C | 0xFE | "HELLO, WORLD" |      | 0xFE | "HELLO, WORLD" |
// +------+--------+------+----------------+      +------+----------------+
//
//
//  <VII. 2 bytes length field at offset 1 in the middle of 4 bytes header,
//	strip the first header field and the length field, the length field
//	represents the length of the whole message>
//
// Let's give another twist to the previous example.  The only difference from
// the previous example is that the length field represents the length of the
// whole message instead of the message body, just like the third example.
// We have to count the length of HDR1 and Length into `LengthAdjustment`.
// Please note that we don't need to take the length of HDR2 into account
// because the length field already includes the whole header length.
//
// LengthFieldOffset   =  1
// LengthFieldLength   =  2
// LengthAdjustment    = -3 (= the length of HDR1 + LEN, negative)
// InitialBytesToStrip =  3
//
// BEFORE DECODE (16 bytes)                       AFTER DECODE (13 bytes)
// +------+--------+------+----------------+      +------+----------------+
// | HDR1 | Length | HDR2 | Actual Content |----->| HDR2 | Actual Content |
// | 0xCA | 0x0010 | 0xFE | "HELLO, WORLD" |      | 0xFE | "HELLO, WORLD" |
// +------+--------+------+----------------+      +------+----------------+

// << 中文含义 By Aceld >>
//
// FrameDecoder
// 一个解码器，根据消息中长度字段的值动态地拆分接收到的二进制数据帧,
// 当您解码具有表示消息正文或整个消息长度的整数头字段的二进制消息时，它特别有用。
//
// ziface.LengthField 有许多配置参数，因此它可以解码任何具有长度字段的消息，
// 这在专有的客户端-服务器协议中经常见到。
//
// 以下是一些示例，它们将为您提供基本的想法，了解每个选项的作用。
//
// 案例一. 在偏移量为0的位置使用2字节长度字段，不去掉消息头
//
// 在这个例子中，长度字段的值是 `12 (0x0C)`，表示"HELLO, WORLD"的长度。
// 默认情况下，解码器会假设长度字段代表跟在长度字段后面的字节数。因此，
// 可以使用简单的参数组合进行解码。
//
// LengthFieldOffset = 0
// LengthFieldLength = 2
// LengthAdjustment = 0
// InitialBytesToStrip = 0 （= 不去掉消息头）
//
// 解码前（14个字节）                  解码后（14个字节）
// +--------+----------------+      +--------+----------------+
// | Length | Actual Content |----->| Length | Actual Content |
// | 0x000C | "HELLO, WORLD" |      | 0x000C | "HELLO, WORLD" |
// +--------+----------------+      +--------+----------------+
//
// 案例二. 位于偏移量0的2字节长度字段，去掉消息头
//
// 由于我们可以通过调用{@link ByteBuf#readableBytes()}来获取内容的长度，
// 因此您可能希望通过指定"InitialBytesToStrip"来去掉长度字段。在此示例中，我们指定了"2"，
// 与长度字段的长度相同，以去掉前两个字节。
//
// LengthFieldOffset = 0
// LengthFieldLength = 2
// LengthAdjustment = 0
// InitialBytesToStrip = 2 （等于Length字段的长度）
//
// 解码前 (14 bytes)         		解码后 (12 bytes)
// +--------+----------------+      +----------------+
// | Length | Actual Content |----->| Actual Content |
// | 0x000C | "HELLO, WORLD" |      | "HELLO, WORLD" |
// +--------+----------------+      +----------------+
//
// 案例三. 位于偏移量0处的2字节长度字段，不剥离头部，该长度字段表示整个消息的长度
//
// 在大多数情况下，长度字段仅表示消息体的长度，就像前面的例子所示。然而，在一些协议中，
// 长度字段表示整个消息的长度，包括消息头部。在这种情况下，我们需要指定一个非零的LengthAdjustment。
// 因为这个例子消息中长度值总是比消息体长度大2，所以我们将LengthAdjustment设置为-2进行补偿。
//
// LengthFieldOffset = 0
// LengthFieldLength = 2
// LengthAdjustment = -2 (长度字段的长度)
// InitialBytesToStrip = 0
//
// 解码前 (14 bytes)         		解码后 (14 bytes)
// +--------+----------------+      +--------+----------------+
// | Length | Actual Content |----->| Length | Actual Content |
// | 0x000E | "HELLO, WORLD" |      | 0x000E | "HELLO, WORLD" |
// +--------+----------------+      +--------+----------------+
//
// 案例四. 5个字节的头部中包含3个字节的长度字段，不去除头部
//
// 下面的消息是第一个示例的简单变体。在消息前面添加了额外的头部值。
// LengthAdjustment 再次为零，因为解码器始终考虑预置数据的长度进行帧长度计算。
//
// LengthFieldOffset = 2（等于 Header 1 的长度）
// LengthFieldLength = 3
// LengthAdjustment = 0
// InitialBytesToStrip = 0
//
// 解码前 (17 bytes)                      		 解码后(17 bytes)
// +----------+----------+----------------+      +----------+----------+----------------+
// | Header 1 | Length   | Actual Content |----->| Header 1 | Length   | Actual Content |
// | 0xCAFE   | 0x00000C | "HELLO, WORLD" |      | 0xCAFE   | 0x00000C | "HELLO, WORLD" |
// +----------+----------+----------------+      +----------+----------+----------------+
//
// 案例五. 在 5 个字节的头部中有 3 个字节的长度字段，不剥离头部
//
// 这是一个高级的例子，展示了在长度字段和消息体之间有额外头部的情况。
// 您需要指定一个正的 LengthAdjustment，以便解码器在帧长度计算中计算额外的头部。
//
// LengthFieldOffset = 0
// LengthFieldLength = 3
// LengthAdjustment = 2 （即 Header 1 的长度）
// InitialBytesToStrip = 0
//
// 解码前 (17 bytes)                      		 解码后 (17 bytes)
// +----------+----------+----------------+      +----------+----------+----------------+
// | Length   | Header 1 | Actual Content |----->| Length   | Header 1 | Actual Content |
// | 0x00000C | 0xCAFE   | "HELLO, WORLD" |      | 0x00000C | 0xCAFE   | "HELLO, WORLD" |
// +----------+----------+----------------+      +----------+----------+----------------+
//
//
// 案例六. 4字节头部，其中2字节长度字段位于偏移量1的位置，剥离第一个头部字段和长度字段
//
// 这是以上所有示例的组合。在长度字段之前有预置的头部，
// 而在长度字段之后有额外的头部。预置头部影响LengthFieldOffset，
// 额外的头部影响LengthAdjustment。我们还指定了一个非零的
// InitialBytesToStrip，以从帧中剥离长度字段和预置头部。
// 如果您不想剥离预置头部，则可以将initialBytesToSkip指定为0。
//
// LengthFieldOffset = 1（HDR1的长度）
// LengthFieldLength = 2
// LengthAdjustment = 1（HDR2的长度）
// InitialBytesToStrip = 3（HDR1 + LEN的长度）
//
// BEFORE DECODE (16 bytes)                       AFTER DECODE (13 bytes)
// +------+--------+------+----------------+      +------+----------------+
// | HDR1 | Length | HDR2 | Actual Content |----->| HDR2 | Actual Content |
// | 0xCA | 0x000C | 0xFE | "HELLO, WORLD" |      | 0xFE | "HELLO, WORLD" |
// +------+--------+------+----------------+      +------+----------------+
//
// 案例七. 2字节长度字段在4字节头部的偏移量为1的位置，去除第一个头字段和长度字段，长度字段表示整个消息的长度
//
// 让我们对前面的示例进行一些变化。
// 与先前的示例唯一的区别在于，长度字段表示整个消息的长度，而不是消息正文，就像第三个示例一样。
// 我们必须将HDR1和Length的长度计入LengthAdjustment。
// 请注意，我们不需要考虑HDR2的长度，因为长度字段已经包含了整个头部的长度。
//
// LengthFieldOffset = 1
// LengthFieldLength = 2
// LengthAdjustment = -3 （HDR1 + LEN的长度，为负数）
// InitialBytesToStrip = 3
//
// 在解码之前（16个字节）                              解码之后（13个字节）
// +------+--------+------+----------------+       +------+----------------+
// | HDR1 | Length | HDR2 | Actual Content |-----> | HDR2 | Actual Content |
// | 0xCA | 0x0010 | 0xFE | "HELLO, WORLD" |       | 0xFE | "HELLO, WORLD" |
// +------+--------+------+----------------+       +------+----------------+

// FrameDecoder 基于LengthField模式的解码器
type FrameDecoder struct {
	ziface.LengthField //从ILengthField集成的基础属性

	LengthFieldEndOffset   int   //长度字段结束位置的偏移量  LengthFieldOffset+LengthFieldLength
	failFast               bool  //快速失败
	discardingTooLongFrame bool  //true 表示开启丢弃模式，false 正常工作模式
	tooLongFrameLength     int64 //当某个数据包的长度超过maxLength，则开启丢弃模式，此字段记录需要丢弃的数据长度
	bytesToDiscard         int64 //记录还剩余多少字节需要丢弃
	in                     []byte
	lock                   sync.Mutex
}

func NewFrameDecoder(lf ziface.LengthField) ziface.IFrameDecoder {

	frameDecoder := new(FrameDecoder)

	//基础属性赋值
	if lf.Order == nil {
		frameDecoder.Order = binary.BigEndian
	} else {
		frameDecoder.Order = lf.Order
	}
	frameDecoder.MaxFrameLength = lf.MaxFrameLength
	frameDecoder.LengthFieldOffset = lf.LengthFieldOffset
	frameDecoder.LengthFieldLength = lf.LengthFieldLength
	frameDecoder.LengthAdjustment = lf.LengthAdjustment
	frameDecoder.InitialBytesToStrip = lf.InitialBytesToStrip

	//self
	frameDecoder.LengthFieldEndOffset = lf.LengthFieldOffset + lf.LengthFieldLength
	frameDecoder.in = make([]byte, 0)

	return frameDecoder
}

func NewFrameDecoderByParams(maxFrameLength uint64, lengthFieldOffset, lengthFieldLength, lengthAdjustment, initialBytesToStrip int) ziface.IFrameDecoder {
	return NewFrameDecoder(ziface.LengthField{
		MaxFrameLength:      maxFrameLength,
		LengthFieldOffset:   lengthFieldOffset,
		LengthFieldLength:   lengthFieldLength,
		LengthAdjustment:    lengthAdjustment,
		InitialBytesToStrip: initialBytesToStrip,
		Order:               binary.BigEndian,
	})
}

func (d *FrameDecoder) fail(frameLength int64) {
	//丢弃完成或未完成都抛异常
	//if frameLength > 0 {
	//	msg := fmt.Sprintf("Adjusted frame length exceeds %d : %d - discarded", this.MaxFrameLength, frameLength)
	//	panic(msg)
	//} else {
	//	msg := fmt.Sprintf("Adjusted frame length exceeds %d - discarded", this.MaxFrameLength)
	//	panic(msg)
	//}
}

func (d *FrameDecoder) discardingTooLongFrameFunc(buffer *bytes.Buffer) {
	//保存还需丢弃多少字节
	bytesToDiscard := d.bytesToDiscard
	//获取当前可以丢弃的字节数，有可能出现半包
	localBytesToDiscard := math.Min(float64(bytesToDiscard), float64(buffer.Len()))
	//fmt.Println("--->", bytesToDiscard, buffer.Len(), localBytesToDiscard)
	//localBytesToDiscard = 2
	//丢弃
	buffer.Next(int(localBytesToDiscard))
	//更新还需丢弃的字节数
	bytesToDiscard -= int64(localBytesToDiscard)
	d.bytesToDiscard = bytesToDiscard
	//是否需要快速失败，回到上面的逻辑
	d.failIfNecessary(false)
}

func (d *FrameDecoder) getUnadjustedFrameLength(buf *bytes.Buffer, offset int, length int, order binary.ByteOrder) int64 {
	//长度字段的值
	var frameLength int64
	arr := buf.Bytes()
	arr = arr[offset : offset+length]
	buffer := bytes.NewBuffer(arr)
	switch length {
	case 1:
		//byte
		var value uint8
		binary.Read(buffer, order, &value)
		frameLength = int64(value)
	case 2:
		//short
		var value uint16
		binary.Read(buffer, order, &value)
		frameLength = int64(value)
	case 3:
		//int占32位，这里取出后24位，返回int类型
		if order == binary.LittleEndian {
			n := uint(arr[0]) | uint(arr[1])<<8 | uint(arr[2])<<16
			frameLength = int64(n)
		} else {
			n := uint(arr[2]) | uint(arr[1])<<8 | uint(arr[0])<<16
			frameLength = int64(n)
		}
	case 4:
		//int
		var value uint32
		binary.Read(buffer, order, &value)
		frameLength = int64(value)
	case 8:
		//long
		binary.Read(buffer, order, &frameLength)
	default:
		panic(fmt.Sprintf("unsupported LengthFieldLength: %d (expected: 1, 2, 3, 4, or 8)", d.LengthFieldLength))
	}
	return frameLength
}

func (d *FrameDecoder) failOnNegativeLengthField(in *bytes.Buffer, frameLength int64, lengthFieldEndOffset int) {
	in.Next(lengthFieldEndOffset)
	panic(fmt.Sprintf("negative pre-adjustment length field: %d", frameLength))
}

func (d *FrameDecoder) failIfNecessary(firstDetectionOfTooLongFrame bool) {
	if d.bytesToDiscard == 0 {
		//说明需要丢弃的数据已经丢弃完成
		//保存一下被丢弃的数据包长度
		tooLongFrameLength := d.tooLongFrameLength
		d.tooLongFrameLength = 0
		//关闭丢弃模式
		d.discardingTooLongFrame = false
		//failFast：默认true
		//firstDetectionOfTooLongFrame：传入true
		if !d.failFast || firstDetectionOfTooLongFrame {
			//快速失败
			d.fail(tooLongFrameLength)
		}
	} else {
		//说明还未丢弃完成
		if d.failFast && firstDetectionOfTooLongFrame {
			//快速失败
			d.fail(d.tooLongFrameLength)
		}
	}
}

// frameLength：数据包的长度
func (d *FrameDecoder) exceededFrameLength(in *bytes.Buffer, frameLength int64) {
	//数据包长度-可读的字节数  两种情况
	//1. 数据包总长度为100，可读的字节数为50，说明还剩余50个字节需要丢弃但还未接收到
	//2. 数据包总长度为100，可读的字节数为150，说明缓冲区已经包含了整个数据包
	discard := frameLength - int64(in.Len())
	//记录一下最大的数据包的长度
	d.tooLongFrameLength = frameLength
	if discard < 0 {
		//说明是第二种情况，直接丢弃当前数据包
		in.Next(int(frameLength))
	} else {
		//说明是第一种情况，还有部分数据未接收到
		//开启丢弃模式
		d.discardingTooLongFrame = true
		//记录下次还需丢弃多少字节
		d.bytesToDiscard = discard
		//丢弃缓冲区所有数据
		in.Next(in.Len())
	}
	//跟进去
	d.failIfNecessary(true)
}

func (d *FrameDecoder) failOnFrameLengthLessThanInitialBytesToStrip(in *bytes.Buffer, frameLength int64, initialBytesToStrip int) {
	in.Next(int(frameLength))
	panic(fmt.Sprintf("Adjusted frame length (%d) is less  than InitialBytesToStrip: %d", frameLength, initialBytesToStrip))
}

func (d *FrameDecoder) decode(buf []byte) []byte {
	in := bytes.NewBuffer(buf)
	//丢弃模式
	if d.discardingTooLongFrame {
		d.discardingTooLongFrameFunc(in)
	}
	////判断缓冲区中可读的字节数是否小于长度字段的偏移量
	if in.Len() < d.LengthFieldEndOffset {
		//说明长度字段的包都还不完整，半包
		return nil
	}
	//执行到这，说明可以解析出长度字段的值了

	//计算出长度字段的开始偏移量
	actualLengthFieldOffset := d.LengthFieldOffset
	//获取长度字段的值，不包括lengthAdjustment的调整值
	frameLength := d.getUnadjustedFrameLength(in, actualLengthFieldOffset, d.LengthFieldLength, d.Order)

	//如果数据帧长度小于0，说明是个错误的数据包
	if frameLength < 0 {
		//内部会跳过这个数据包的字节数，并抛异常
		d.failOnNegativeLengthField(in, frameLength, d.LengthFieldEndOffset)
	}

	//套用前面的公式：长度字段后的数据字节数=长度字段的值+lengthAdjustment
	//frameLength就是长度字段的值，加上lengthAdjustment等于长度字段后的数据字节数
	//lengthFieldEndOffset为lengthFieldOffset+lengthFieldLength
	//那说明最后计算出的framLength就是整个数据包的长度
	frameLength += int64(d.LengthAdjustment) + int64(d.LengthFieldEndOffset)
	//丢弃模式就是在这开启的
	//如果数据包长度大于最大长度
	if uint64(frameLength) > d.MaxFrameLength {
		//对超过的部分进行处理
		d.exceededFrameLength(in, frameLength)
		return nil
	}

	//执行到这说明是正常模式
	//数据包的大小
	frameLengthInt := int(frameLength)
	//判断缓冲区可读字节数是否小于数据包的字节数
	if in.Len() < frameLengthInt {
		//半包，等会再来解析
		return nil
	}

	//执行到这说明缓冲区的数据已经包含了数据包

	//跳过的字节数是否大于数据包长度
	if d.InitialBytesToStrip > frameLengthInt {
		d.failOnFrameLengthLessThanInitialBytesToStrip(in, frameLength, d.InitialBytesToStrip)
	}
	//跳过initialBytesToStrip个字节
	in.Next(d.InitialBytesToStrip)
	//解码
	//获取跳过后的真实数据长度
	actualFrameLength := frameLengthInt - d.InitialBytesToStrip
	//提取真实的数据
	buff := make([]byte, actualFrameLength)
	in.Read(buff)
	//bytes.NewBuffer([]byte{})
	//_in := bytes.NewBuffer(buff)
	return buff
}

func (d *FrameDecoder) Decode(buff []byte) [][]byte {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.in = append(d.in, buff...)
	resp := make([][]byte, 0)

	for {
		arr := d.decode(d.in)

		if arr != nil {
			//证明已经解析出一个完整包
			resp = append(resp, arr)
			_size := len(arr) + d.InitialBytesToStrip
			//_len := len(this.in)
			//fmt.Println(_len)
			if _size > 0 {
				d.in = d.in[_size:]
			}
		} else {
			return resp
		}
	}
}
