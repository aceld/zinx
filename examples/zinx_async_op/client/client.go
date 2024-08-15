package main

import (
	"fmt"
	"io"
	"net"

	"github.com/aceld/zinx/zpack"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	dp := zpack.NewDataPack()
	msg, _ := dp.Pack(zpack.NewMsgPackage(1, []byte("async_op_router test=========>")))
	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println("write error err ", err)
		return
	}

	for {
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("client read head err: ", err)
			return
		}

		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client unpack head err: ", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*zpack.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("client unpack data err")
				return
			}

			fmt.Printf("==> Client receive Msg: ID = %d, len = %d , data = %s\n", msg.ID, msg.DataLen, msg.Data)
		}
	}

}
