package main

import (
	"errors"
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"time"
)

type TestRouter struct {
	znet.BaseRouter
}

// PreHandle -
func (t *TestRouter) PreHandle(req ziface.IRequest) {
	start := time.Now()

	fmt.Println("--> Call PreHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test1")); err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(start)
	fmt.Println("cost time：", elapsed)
}

// Handle -
func (t *TestRouter) Handle(req ziface.IRequest) {
	fmt.Println("--> Call Handle")

	// Simulated scenario - In the event of an expected error such as incorrect permissions or incorrect information,
	// subsequent function execution will be stopped, but this function will be fully executed.
	// 模拟场景- 出现意料之中的错误 如权限不对或者信息错误 则停止后续函数执行，但是次函数会执行完毕
	if err := Err(); err != nil {
		req.Abort()
		fmt.Println("Insufficient permission")
	}

	// Simulation scenario - In case of a certain situation, repeat the above operation.
	// 模拟场景- 出现某种情况，重复上面的操作
	/*
		if err := Err(); err != nil {
			req.Goto(znet.PRE_HANDLE)
			fmt.Println("repeat")
		}
	*/

	if err := req.GetConnection().SendMsg(0, []byte("test2")); err != nil {
		fmt.Println(err)
	}

	time.Sleep(1 * time.Millisecond)
}

// PostHandle -
func (t *TestRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("--> Call PostHandle")
	if err := req.GetConnection().SendMsg(0, []byte("test3")); err != nil {
		fmt.Println(err)
	}
}

func Err() error {
	//Specific Business Operation (具体业务操作)
	return errors.New("Test")
}

func main() {
	s := znet.NewServer()
	s.AddRouter(1, &TestRouter{})
	s.Serve()
}
