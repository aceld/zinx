package znet

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"

	"github.com/chnkenc/zinx-xiaoan/utils"
	"github.com/chnkenc/zinx-xiaoan/ziface"
)

var (
	// 默认数据包头长度
	defaultHeaderLen uint16 = 6

	// 默认魔数
	defualtMagicCode = "aa55"
)

//DataPack 封包拆包类实例，暂时不需要成员
type DataPack struct{}

//NewDataPack 封包拆包实例初始化方法
func NewDataPack() ziface.Packet {
	return &DataPack{}
}

//GetHeadLen 获取包头长度方法
func (dp *DataPack) GetHeadLen() uint16 {
	// 魔数 uint16(2字节) +  命令字 uint8(1字节) + 序列号 uint8（1字节）+ 长度 uint16（2字节）
	return defaultHeaderLen
}

//Pack 封包方法(压缩数据)
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 写魔数
	defaultMagicCodeByte, err := hex.DecodeString(defualtMagicCode)

	if err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.BigEndian, defaultMagicCodeByte); err != nil {
		return nil, err
	}

	// 写msgID（命令字）
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	// 写序列号
	serialSn := msg.GetSerialSn()
	if serialSn == 0 {
		serialSn = dp.GenerateSerialSn()
	}

	if err := binary.Write(dataBuff, binary.BigEndian, serialSn); err != nil {
		return nil, err
	}

	// 写dataLen
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写data数据
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//Unpack 拆包方法(解压数据)
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}

// GenerateSerialSn 生成序列号
func (dp *DataPack) GenerateSerialSn() uint8 {
	rand.Seed(time.Now().UnixNano())
	sn := rand.Intn(128)

	return uint8(sn)
}
