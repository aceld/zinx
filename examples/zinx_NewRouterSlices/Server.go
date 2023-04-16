package main

import (
	"fmt"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

func Test1(request ziface.IRequest) {
	fmt.Println("test1")
}

func Test2(request ziface.IRequest) {
	fmt.Println("Test2")
}
func Test3(request ziface.IRequest) {
	fmt.Println("Test3")
}
func Test4(request ziface.IRequest) {
	fmt.Println("Test4")
}
func Test5(request ziface.IRequest) {
	fmt.Println("Test5")
}
func Test6(request ziface.IRequest) {
	fmt.Println("Test6")
}

func main() {
	server := znet.NewUserConfServer(&zconf.Config{RouterMode: 2, TCPPort: 9512, Host: "127.0.0.1"})
	//分组
	group := server.Group(3, 10, Test1)

	//添加路由 如果不在组范围会直接panic
	//group.AddHandler(11, Test2)

	//在组中 不受Use影响 有 1 2 处理器
	group.AddHandler(3, Test2)

	//既不在组里也在Use之前只会有自己的处理器 3
	server.AddRouterSlices(1, Test3)

	//如果希望group处理器优先，应当在Use之前如下操作
	//可以手动添加 入 group.AddHandler(5, Test4, Test5,Test2, Test3, Test6)
	//或者如下使用Group的Use方法 那么就是 1 4 5 6的顺序 不被use影响
	group.Use(Test2, Test3)
	group.AddHandler(5, Test4, Test5, Test6)

	//公共组件，但是，在使用Use之前的组或者路由不会影响到
	router := server.Use(Test4, Test5)

	//添加路由 不在组中但是收Use影响 有4 5 6处理器
	router.AddHandler(2, Test6)

	//在组里也受到Use影响 有所有处理器 且顺序应该是 4 5 1 2 3 6 因为use中的处理器始终在最前端
	group.AddHandler(4, Test6)

	server.Serve()
}
