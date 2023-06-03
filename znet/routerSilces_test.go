package znet

import (
	"fmt"
	"testing"

	"github.com/aceld/zinx/ziface"
)

func A1(request ziface.IRequest) {
	fmt.Println("test A1")
}
func A2(request ziface.IRequest) {
	fmt.Println("test A2")
}
func A3(request ziface.IRequest) {
	fmt.Println("test A3")
}
func A4(request ziface.IRequest) {
	fmt.Println("test A4")
}
func A5(request ziface.IRequest) {
	fmt.Println("test A5")
}

func TestRouterAdd(t *testing.T) {
	router := NewRouterSlices()
	router.Use(A3)
	router.AddHandler(1, A1, A2)

	testgroup := router.Group(2, 5, A5)
	{
		testgroup.AddHandler(2, A4)

		//正确panic
		//testgroup.AddHandler(6, A4)
	}

	for _, v := range router.Apis[2] {
		v(&Request{
			index: -1,
		})
	}
}
