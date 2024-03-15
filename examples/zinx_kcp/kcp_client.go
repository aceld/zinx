package main

import (
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zpack"
	"github.com/xtaci/kcp-go"
	"io"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("Client Test ... start")
	// Replace net.Dial with kcp.DialWithOptions
	conn, err := kcp.Dial("127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	dp := zpack.Factory().NewPack(ziface.ZinxDataPack)
	sendMsg, _ := dp.Pack(zpack.NewMsgPackage(1, []byte("client test message")))
	_, err = conn.Write(sendMsg)
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

			time.Sleep(1 * time.Second)
			_, err = conn.Write(sendMsg)
			if err != nil {
				fmt.Println("client write err: ", err)
				return
			}
		}
	}
}
