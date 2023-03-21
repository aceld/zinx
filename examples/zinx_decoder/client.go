package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
	"time"
)

// 使用该方法生成模拟数据
func getTLVPackData() []byte {
	msgID := 1
	tag := make([]byte, 4)
	binary.BigEndian.PutUint32(tag, uint32(msgID))

	str := "HELLO, WORLD"
	var value = []byte(str)

	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(value)))

	_data := make([]byte, 0)
	_data = append(_data, tag...)
	_data = append(_data, length...)
	_data = append(_data, value...)
	fmt.Println("--->", len(_data), hex.EncodeToString(_data))
	return _data
}

func getTLVData(index int) []byte {
	//通过 getTLVPackData()方法，获得一段完整的TLV模拟数据包:000000010000000c48454c4c4f2c20574f524c44
	tlvPackData := []string{
		"000000010000000c48454c4c4f2c20574f524c44000000010000000c",                         //一包半
		"48454c4c4f2c20574f524c44",                                                         //剩下的半包
		"000000010000000c48454c4c4f2c20574f524c44000000010000000c48454c4c4f2c20574f524c44", //两包
	}
	//此处模拟顺序如:两包一包半剩下的半包
	index = index % 3
	if index == 0 {
		fmt.Println("模拟-粘包")
		index = 2 //模拟粘包情况，两包数据一起
	} else {
		index = index / 2 % 2 //模拟断包情况，一包半+剩下的半包
		fmt.Println("模拟-断包")
	}
	arr, _ := hex.DecodeString(tlvPackData[index])
	return arr
}

func getHTLVCRCData(index int) []byte {
	//一段完整的HTLVCRC模拟数据包:A2100E0102030405060708091011121314050B
	tlvPackData := []string{
		"a21018686574000004d30000000000000000000000000000000000e7a2a2130e686574000004d300000001", //一包半
		"00000040c3", //剩下的半包
		"a21018686574000004d30000000000000000000000000000000000e7a2a2130e686574000004d30000000100000040c3", //两包
	}
	//此处模拟顺序如:两包一包半剩下的半包
	index = index % 3
	if index == 0 {
		fmt.Println("模拟-粘包")
		index = 2 //模拟粘包情况，两包数据一起
	} else {
		index = index / 2 % 2 //模拟断包情况，一包半+剩下的半包
		fmt.Println("模拟-断包")
	}
	arr, _ := hex.DecodeString(tlvPackData[index])
	return arr
}

// 客户端自定义业务
func business(conn ziface.IConnection) {
	var i int
	for {
		//buffer := getTLVData(i)
		buffer := getHTLVCRCData(i)
		conn.Send(buffer)
		i++
		time.Sleep(1 * time.Second)
	}
}

// 创建连接的时候执行
func DoClientConnectedBegin(conn ziface.IConnection) {
	zlog.Debug("DoConnecionBegin is Called ... ")
	go business(conn)
}

func main() {
	//创建一个Client句柄，使用Zinx的API
	client := znet.NewClient("127.0.0.1", 8999)
	//添加首次建立链接时的业务
	client.SetOnConnStart(DoClientConnectedBegin)
	//启动客户端client
	client.Start()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("===exit===", sig)

}
