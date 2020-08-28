package zd

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

//zinx node集群相互通信的消息体
type ZdMessage struct {
	CmdId   int32
	Unit    ZinxUnit
	DataLen uint32
	Data    []byte
}

//新建一个ZdMessage结构体
func NewZdMessage(cmdId int, unit *ZinxUnit, data []byte) *ZdMessage {
	return &ZdMessage{
		CmdId:   int32(cmdId),
		Unit:    *unit,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

//将传输的消息打包
func ZdMsgPack(msg *ZdMessage) ([]byte, error) {

	//创建一个存放bytes字节的缓冲
	buf := new(bytes.Buffer)

	//写cmdID (4个字节)
	if err := binary.Write(buf, binary.LittleEndian, (int32)(msg.CmdId)); err != nil {
		fmt.Println("ZdMsgPack encode cmdId error, ", err)
		return nil, err
	}

	unitStr, err := json.Marshal(msg.Unit)
	if err != nil {
		fmt.Println("json Marshal unit err, unitStr = ", unitStr)
		return nil, err
	}
	unitLen := len(unitStr)

	//写unit 数据长度  (4个字节)
	if err := binary.Write(buf, binary.LittleEndian, unitLen); err != nil {
		fmt.Println("ZdMsgPack encode unitLen error,", err)
		return nil, err
	}

	//写unitStr  数据长度(unitLen个字节)
	if err := binary.Write(buf, binary.LittleEndian, unitStr); err != nil {
		fmt.Println("ZdMsgPack encode unitStr error,", err)
		return nil, err
	}

	//写dataLen(4个字节)
	if err := binary.Write(buf, binary.LittleEndian, msg.DataLen); err != nil {
		fmt.Println("ZdMsgPack encode dataLen error, ", err)
		return nil, err
	}

	//写data
	if err := binary.Write(buf, binary.LittleEndian, msg.Data); err != nil {
		fmt.Println("ZdMsgPack encode data error, ", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

//解压包裹，到当前的message中
func ZdMsgUnpack(binaryData []byte) (*ZdMessage, error) {

	//创建一个从输入二进制数据的ioReader
	buf := bytes.NewReader(binaryData)

	msg := &ZdMessage{}

	//读cmdID(4个字节)
	if err := binary.Read(buf, binary.LittleEndian, &msg.CmdId); err != nil {
		fmt.Println("ZdMsgUnpack decode cmdId error, ", err)
		return nil, err
	}

	//读unitLen(4个字节)
	var unitLen int32
	if err := binary.Read(buf, binary.LittleEndian, &unitLen); err != nil {
		fmt.Println("ZdMsgUnpack decode unitLen error, ", err)
		return nil, err
	}

	//读取unitData(unitLen个字节)
	var jsonBuf []byte = make([]byte, unitLen)
	if err := binary.Read(buf, binary.LittleEndian, jsonBuf); err != nil {
		fmt.Println("ZdMsgUnpack decode unitData error, ", err)
		return nil, err
	}

	if err := json.Unmarshal(jsonBuf, &msg.Unit); err != nil {
		fmt.Println("ZdMsgUnpack Unmarshal json unit error, ", err)
		return nil, err
	}

	//读取dataLen(4个字节)
	if err := binary.Read(buf, binary.LittleEndian, &msg.DataLen); err != nil {
		fmt.Println("ZdMsgUnpack decode DataLen error, ", err)
		return nil, err
	}

	//读取DataBuf(dataLen个字节)
	var dataBuf []byte = make([]byte, msg.DataLen)
	if err := binary.Read(buf, binary.LittleEndian, dataBuf); err != nil {
		fmt.Println("ZdMsgUnpack decode Data error, ", err)
		return nil, err
	}

	msg.Data = dataBuf[:]

	return msg, nil
}
