package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

/*
	模拟客户端
 */
 func ClientTest() {

 	fmt.Println("Client Test ... start")
 	//3秒之后发起测试请求，给服务端开启服务的机会
 	time.Sleep(3 * time.Second)

 	conn,err := net.Dial("tcp", "127.0.0.1:7777")
 	if err != nil {
 		fmt.Println("client start err, exit!")
		return
	}

 	for {
 		_, err := conn.Write([]byte("hahaha"))
 		if err !=nil {
 			fmt.Println("write error err ", err)
 			return
		}

 		buf :=make([]byte, 512)
 		cnt, err := conn.Read(buf)
 		if err != nil {
 			fmt.Println("read buf error ")
 			return
		}

 		fmt.Printf(" server call back : %s, cnt = %d\n", buf,  cnt)

 		time.Sleep(1*time.Second)
	}
 }

//Server 模块的测试函数
func TestServer(t *testing.T) {

	/*
		服务端测试
	*/
	//1 创建一个server 句柄 s
	s := NewServer("[zinx V0.1]")

	/*
		客户端测试
	*/
	go ClientTest()

	//2 开启服务
	s.Serve()
}