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
	msg, _ := dp.Pack(zpack.NewMsgPackage(1, []byte("ZinxPing")))
	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println("write error err ", err)
		return
	}

}
