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

type router struct {
	znet.BaseRouter
}

func (r *router) PreHandle(req ziface.IRequest) {
	fmt.Println(" hello router1")
}
func (r *router) Handle(req ziface.IRequest) {
	req.Abort()
	fmt.Println(" hello router2")
}
func (r *router) PostHandle(req ziface.IRequest) {
	fmt.Println(" hello router3")
}

func main() {

	// Old version router method (旧版本路由方法)
	//{
	//	server := znet.NewUserConfServer(&zconf.Config{TCPPort: 8999, Host: "127.0.0.1"})
	//
	//	// Even without manually calling the router mode, the default is 1 (old version) 即使不手动调路由模式也可以,默认是1（旧版本）
	//	//server := znet.NewServer()
	//
	//	// Old version runs normally(旧版正常执行)
	//	r := &router{}
	//	server.AddRouter(1, r)
	//	server.Serve()
	//}
	//{

	// New version usage and explanation(新版本使用方法以及说明)
	{
		server := znet.NewUserConfServer(&zconf.Config{RouterSlicesMode: true, TCPPort: 8999, Host: "127.0.0.1"})
		// Grouping(分组)
		group := server.Group(3, 10, Test1)

		// Add router. Will panic if not within the group range.(添加路由 如果不在组范围会直接panic)
		//group.AddHandler(11, Test2)

		// Within the group, not affected by Use, has processors 1 and 2.(在组中 不受Use影响 有 1 2 处理器)
		group.AddHandler(3, Test2)

		// Not within the group and before Use, only has its own processor 3.(既不在组里也在Use之前只会有自己的处理器 3)
		server.AddRouterSlices(1, Test3)

		// If you want the group processor to have priority, you should do the following before Use.
		// You can manually add it via group.AddHandler(5, Test4, Test5,Test2, Test3, Test6)
		// or use the Group's Use method as shown below, which would have the order of 1 4 5 6 and not be affected by Use.
		// 如果希望group处理器优先，应当在Use之前如下操作
		// 可以手动添加 入 group.AddHandler(5, Test4, Test5,Test2, Test3, Test6)
		// 或者如下使用Group的Use方法 那么就是 1 4 5 6的顺序 不被use影响
		group.Use(Test2, Test3)
		group.AddHandler(5, Test4, Test5, Test6)

		// Common components, but not affected by the groups or routers before Use. (公共组件，但是，在使用Use之前的组或者路由不会影响到)
		router := server.Use(Test4, Test5)
		// Add router. Not within the group but is affected by Use, has processors 4, 5, and 6.
		// (添加路由 不在组中但是收Use影响 有4 5 6处理器)
		router.AddHandler(2, Test6)

		// Within the group and affected by Use, has all processors in the order of 4, 5, 1, 2, 3, 6 because the processors in Use are always at the forefront.
		// (在组里也受到Use影响 有所有处理器 且顺序应该是 4 5 1 2 3 6 因为use中的处理器始终在最前端)
		group.AddHandler(4, Test6)

		server.Serve()
	}

}
