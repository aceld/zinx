package main

import (
	"fmt"
	"github.com/aceld/zinx/zpack"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	dp := zpack.NewDataPack()
	for i := 1; i < 4; i++ {
		msg, _ := dp.Pack(zpack.NewMsgPackage(uint32(i), []byte("ZinxPing")))
		fmt.Println(i)
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}
	}

}
