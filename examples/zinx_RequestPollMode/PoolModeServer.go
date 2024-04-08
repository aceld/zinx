package main

import (
	"fmt"
	"time"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

func Poll1(request ziface.IRequest) {
	// 如果需要连接信息
	request.Set("conn", request.GetConnection())
	request.Set("num", 1)
	fmt.Printf("request 1 addr:%p,conn:%p \n", &request, request.GetConnection())

	// 需要新线程同时也需要上下文的情况,则需要调用 copy 方法拷贝一份
	cp := request.Copy()
	go Poll2(cp)

	// 如果不使用 copy 方法拷贝对象则会出现同一个对象但是信息可能不一致的问题,不启动 poll2 会更直接
	go Poll3(request)
}

func Poll2(request ziface.IRequest) {
	defer func() {
		if err := recover(); err != nil {
			// 接收一个panic
			fmt.Println(err)
		}

	}()
	get_conn, ok := request.Get("conn")
	if ok {
		// 如果直接取用则会导致空指针
		request.GetConnection().GetConnID()
		//  打印出的 Request 对象的地址是不一致的
		conn := get_conn.(ziface.IConnection)
		fmt.Printf("request copy addr:%p,conn:%p \n", &request, conn)
		// conn.sendMsg()
	}
}

// 如果请求的次数多,则开启对象池且直接传递不copy Request 就可能导致值不一致
func Poll3(request ziface.IRequest) {
	time.Sleep(time.Second * 3)
	get, _ := request.Get("num")
	// 池化对象如果直接传递被影响可能随机打印被修改的值 3
	fmt.Printf("num:%v \n", get)

}

func Poll4(request ziface.IRequest) {
	// 影响原本的 request 对象
	request.Set("num", 3)
}

func main() {

	// 开启 Request 对象池模式
	server := znet.NewUserConfServer(&zconf.Config{RouterSlicesMode: true, TCPPort: 8999, Host: "127.0.0.1", RequestPoolMode: true})
	server.AddRouterSlices(1, Poll1)
	server.AddRouterSlices(2, Poll4)
	server.Serve()
}
