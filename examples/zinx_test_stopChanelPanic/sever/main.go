package main

import (
	"fmt"

	"github.com/aceld/zinx/examples/zinx_test_stopChanelPanic/router"
	"github.com/aceld/zinx/znet"
)

func main() {
	fmt.Println("[Server] server start")
	s := znet.NewServer()
	s.AddRouter(1, &router.PanicTestRouter{znet.BaseRouter{}})
	s.Serve()
}

// ... existing code ...
