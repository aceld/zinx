package znet

import (
	"fmt"
	"testing"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
)

func A1(request ziface.IRequest) {
	fmt.Println("我要写入一些上下文到 Request 中了")
	request.Set("Hey", "zinx!")
	request.Set("Age", 2)
}
func A2(request ziface.IRequest) {
	name, _ := request.Get("Hey")
	age, _ := request.Get("Age")

	fmt.Printf("我是练习时长%v年半的%v \n", age, name)

	//如果需要开新协程操作应该 copy
	cp := request.Copy()
	go A4(cp)
	request.Abort()
}

func A3(request ziface.IRequest) {
	fmt.Println("No! 不带我玩")
}

func A4(request ziface.IRequest) {
	// 需要新线程同时也需要上下文的情况
	fmt.Println(request)
}

func TestRouterAdd(t *testing.T) {

	server := NewUserConfServer(&zconf.Config{RouterSlicesMode: true, TCPPort: 8999, Host: "127.0.0.1"})
	server.AddRouterSlices(1, A1, A2, A3)
	server.Serve()

}
