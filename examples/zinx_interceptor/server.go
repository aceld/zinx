package main

import (
	"github.com/aceld/zinx/examples/zinx_interceptor/interceptors"
	"github.com/aceld/zinx/examples/zinx_interceptor/router"
	"github.com/aceld/zinx/znet"
)

func main() {
	// 创建server 对象
	server := znet.NewServer()
	// 添加路由映射
	server.AddRouter(1, &router.HelloRouter{})
	// 添加自定义拦截器
	server.AddInterceptor(&interceptors.MyInterceptor{})
	// 启动服务
	server.Serve()
}
