package main

import (
	"fmt"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
)

func Auth1(request ziface.IRequest) {

	// Verify business, default to pass. (验证业务 默认固定放行)
	fmt.Println("I am the Auth1, I will always pass.")
	// I am validation handler 1, and I must pass.
	// Note that the next function will start executing and return here after all functions are executed.
	// (注意是进入下一个函数开始执行，全部执行完后 回到此处)
	request.RouterSlicesNext()
}

func Auth2(request ziface.IRequest) {
	// I am the validation handler 2, and by default, I do not allow the request to pass.(验证业务 默认固定不放行)

	// Terminate execution function, no more handlers will be executed after this one.(终结执行函数，再这个处理器结束后不会在执行后面的处理器)
	request.Abort()
	fmt.Println("I am the Auth2, I will definitely not pass.")
	fmt.Println("The business terminates here and the subsequent handlers will not be executed.")
}

func Auth3(request ziface.IRequest) {

	fmt.Println("I am the group validation function.")
}

// I am a business function.
func TestFunc(request ziface.IRequest) {
	fmt.Println("I am a business function.")
}

func main() {

	// New version usage and explanation.(新版本使用方法以及说明)
	server := znet.NewUserConfServer(&zconf.Config{RouterSlicesMode: true, TCPPort: 8999, Host: "127.0.0.1"})

	// Simulation scenario 1: A normal business that only executes a single operation function separately.
	// 模拟场景 1，普通业务单独只执行一个操作函数
	//server.AddRouterSlices(1, TestFunc)

	// Simulated scenario 2: All operations below require verification of request permissions, so a verification function is needed.
	// the verification component has been added to all the routes under the use method, such as 1 and 2.
	// 模拟场景 2, 以下所有操作都需要验证请求权限，所以需要一个验证函数
	// 将验证组件添加到了所有再use方法下的所有路由中了 如1,2 中都会带有
	//routerSlices := server.Use(Auth1)
	//routerSlices.AddHandler(1, TestFunc)
	//routerSlices.AddHandler(2, TestFunc)

	// Equivalent to the following:(等价于下面:)
	//routerSlices.AddHandler(1, Auth1, TestFunc)

	// Simulated scenario 3: Authorization is required, but some route operations require additional verification.
	// 模拟场景3 需要权限，但是某些路由操作需要更多额外验证
	server.Use(Auth1)
	group1 := server.Group(1, 2, Auth3)
	{
		// MsgId=1, there will be Auth3 and Auth1. (1中就会有Auth3和Auth1)
		group1.AddHandler(1, TestFunc)

		// More specific scenario: Some operations within the group require an additional validation process.
		// 更特殊的情况，组内另一些操作还需要另一道校验处理
		group1.Use(Auth2)
		// MsgId=2, Auth3 and Auth1 will be added to all routes under the use method.
		// 2中就会有Auth1和Auth3以及Auth2
		group1.AddHandler(2, TestFunc)

	}
	// MsgId=3, Auth3 will not be included. (3中就不会有Auth3)
	server.AddRouterSlices(3, TestFunc)

	server.Serve()

}
