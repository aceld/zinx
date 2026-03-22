package main

import (
	"github.com/aceld/zinx/v3/examples/zinx_interceptor/interceptors"
	"github.com/aceld/zinx/v3/examples/zinx_interceptor/router"
	"github.com/aceld/zinx/v3/znet"
)

func main() {
	server := znet.NewServer()

	server.AddRouter(1, &router.HelloRouter{})

	// Add Custom Interceptor
	server.AddInterceptor(&interceptors.MyInterceptor{})

	server.Serve()
}
