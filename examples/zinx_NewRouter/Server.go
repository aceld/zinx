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

//PreHandle -
func (t *TestRouter) PreHandle(req ziface.IRequest) {
	//使用场景模拟  完整路由计时
	start := time.Now()
	req.Next()

	fmt.Println("1")
	if err := req.GetConnection().SendMsg(0, []byte("test1")); err != nil {
		fmt.Println(err)
	}
	elapsed := time.Since(start)
	fmt.Println("该路由组执行完成耗时：", elapsed)
}

//Handle -
func (t *TestRouter) Handle(req ziface.IRequest) {
	fmt.Println("2")

	//模拟场景2 出现意料之中的错误 如权限不对或者信息错误 则停止后续函数执行，但是次函数会执行完毕
	if err := Err(); err != nil {
		req.Abort()
		fmt.Println("权限不足")
	}

	if err := req.GetConnection().SendMsg(0, []byte("test2")); err != nil {
		fmt.Println(err)
	}

	time.Sleep(1 * time.Millisecond) //模拟函数计时
}

//PostHandle -
func (t *TestRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("3")
	if err := req.GetConnection().SendMsg(0, []byte("test3")); err != nil {
		fmt.Println(err)
	}
}

func Err() error {
	//具体业务操作

	return errors.New("Test")
}

func main() {
	s := znet.NewServer()
	s.AddRouter(1, &TestRouter{})
	s.Serve()
}
