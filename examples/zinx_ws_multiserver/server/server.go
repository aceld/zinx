package main

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_server/s_router"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Set global mode to websocket-only.
	zconf.GlobalObject.Mode = zconf.ServerModeWebsocket

	// Start two WebSocket servers on different ports with the same path "/".
	// Before the fix (http.HandleFunc on DefaultServeMux), this would panic:
	//   http: multiple registrations for /
	// After the fix (per-server http.ServeMux), each server handles its own path independently.
	servers := []struct {
		name string
		port int
	}{
		{"WsServer-1", 9001},
		{"WsServer-2", 9002},
	}

	for _, s := range servers {
		server := znet.NewUserConfServer(&zconf.Config{
			Name:   s.name,
			Host:   "0.0.0.0",
			WsPort: s.port,
			WsPath: "/",
		})
		server.AddRouter(100, &s_router.PingRouter{})
		server.AddRouter(1, &s_router.HelloZinxRouter{})
		go server.Serve()
	}

	fmt.Println("Both WebSocket servers started. No route pollution!")
	fmt.Println("Server 1: ws://localhost:9001/")
	fmt.Println("Server 2: ws://localhost:9002/")

	// Wait for signal to stop
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("===exit===")
}
