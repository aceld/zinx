package main

import (
	"fmt"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

func Auth1(request ziface.IRequest) {

	//验证业务 默认固定放行
	fmt.Println("我是验证处理器1，我一定通过")
	//注意是进入下一个函数开始执行，全部执行完后 回到此处
	request.RouterSlicesNext()
}

func Auth2(request ziface.IRequest) {
	//验证业务 默认固定不放行
	//RouterAbort 终结执行函数，再这个处理器结束后不会在执行后面的处理器
	request.RouterAbort()
	fmt.Println("我是验证处理器2，我一定不通过")
	fmt.Println("业务到此终止不会执行后面的")
}

func Auth3(request ziface.IRequest) {

	fmt.Println("我是组验证函数")
}

// 实际业务
func TestFunc(request ziface.IRequest) {
	fmt.Println("我是业务函数")
}

func main() {

	//新版本使用方法以及说明
	server := znet.NewUserConfServer(&zconf.Config{RouterSlicesMode: true, TCPPort: 8999, Host: "127.0.0.1"})

	//模拟场景 1，普通业务单独只执行一个操作函数
	//server.AddRouterSlices(1, TestFunc)

	//模拟场景 2, 以下所有操作都需要验证请求权限，所以需要一个验证函数
	//将验证组件添加到了所有再use方法下的所有路由中了 如1,2 中都会带有
	//routerSlices := server.Use(Auth1)
	//routerSlices.AddHandler(1, TestFunc)
	//routerSlices.AddHandler(2, TestFunc)

	//等价于下面
	//routerSlices.AddHandler(1, Auth1, TestFunc)

	//模拟场景3 需要权限，但是某些路由操作需要更多额外验证
	server.Use(Auth1)
	group1 := server.Group(1, 2, Auth3)
	{
		//1中就会有Auth3和Auth1
		group1.AddHandler(1, TestFunc)

		//更特殊的情况，组内另一些操作还需要另一道校验处理
		group1.Use(Auth2)
		//2中就会有Auth1和Auth3以及Auth2
		group1.AddHandler(2, TestFunc)

	}
	//3中就不会有Auth3
	server.AddRouterSlices(3, TestFunc)

	server.Serve()

}
