package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zpack"
	"io"
	"net"
	"time"
)

// 模拟客户端
func main() {

	fmt.Println("Client Test ... start")
	// Send a test request after 3 seconds to give the server a chance to start the service. (3秒之后发起测试请求，给服务端开启服务的机会)
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	dp := zpack.Factory().NewPack(ziface.ZinxDataPack)
	msg, _ := dp.Pack(zpack.NewMsgPackage(1, []byte("client test message")))
	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println("client write err: ", err)
		return
	}

	for {
		// Read the "head" section from the stream first. (先读出流中的head部分)
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("client read head err: ", err)
			return
		}

		// Unpack the headData byte stream into msg. (将headData字节流 拆包到msg中)
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client unpack head err: ", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			// Read the "data" section from the stream. (再读出流中的data部分)
			msg := msgHead.(*zpack.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// read from io.Reader into msg.Data (根据dataLen从io中读取字节流)
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("client unpack data err")
				return
			}

			fmt.Printf("==> Client receive Msg: ID = %d, len = %d , data = %s\n", msg.ID, msg.DataLen, msg.Data)
		}
	}
}
