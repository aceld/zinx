package main

import (
	"github.com/aceld/zinx/examples/zinx_interceptor/interceptors"
	"github.com/aceld/zinx/examples/zinx_interceptor/router"
	"github.com/aceld/zinx/znet"
)

func main() {
	server := znet.NewServer()

	server.AddRouter(1, &router.HelloRouter{})

	// Add Custom Interceptor
	server.AddInterceptor(&interceptors.MyInterceptor{})

	server.Serve()
}
