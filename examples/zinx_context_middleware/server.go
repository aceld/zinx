// Example of using zinx with context-based middleware
// (使用zinx和基于Context的中间件的示例)
package main

import (
	"fmt"
	"log"

	"github.com/aceld/zinx/v3/ziface"
	"github.com/aceld/zinx/v3/znet"
)

// CustomMiddleware is an example middleware that demonstrates context usage
// (自定义中间件示例，演示context的使用)
func CustomMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		log.Printf("CustomMiddleware: Processing message %d from connection %d", c.MsgID, c.Conn.GetConnID())

		// Set some data in the context
		// (在context中设置一些数据)
		c.Set("middleware_data", "custom_value")

		// Continue to the next middleware
		// (继续执行下一个中间件)
		c.Next()

		log.Printf("CustomMiddleware: Finished processing message %d", c.MsgID)
	}
}

// MessageHandler is the final handler for the message
// (消息的最终处理程序)
func MessageHandler(c *ziface.Context) {
	// Get data from context
	// (从context获取数据)
	data, exists := c.Get("middleware_data")
	if exists {
		log.Printf("MessageHandler: Received middleware data: %v", data)
	}

	// Process the message
	// (处理消息)
	log.Printf("MessageHandler: Processing message ID %d, data: %s", c.MsgID, string(c.Data))

	// Send a response back to the client
	// (向客户端发送响应)
	response := fmt.Sprintf("Server received: %s", string(c.Data))
	err := c.Conn.SendMsg(c.MsgID, []byte(response))
	if err != nil {
		log.Printf("MessageHandler: Error sending response: %v", err)
	}
}

func main() {
	// Create a new server
	// (创建一个新服务器)
	s := znet.NewServer()

	// Register global middleware
	// (注册全局中间件)
	s.UseContext(
		znet.RecoveryMiddleware(), // Recovery middleware (恢复中间件)
		znet.LoggingMiddleware(),  // Logging middleware (日志中间件)
		znet.TraceMiddleware(),    // Trace middleware (追踪中间件)
	)

	// Register route-specific middleware and handler
	// (注册路由特定的中间件和处理程序)
	s.AddRouterSlicesContext(1, // Message ID 1
		CustomMiddleware(), // Custom middleware (自定义中间件)
		MessageHandler,     // Final handler (最终处理程序)
	)

	// Start the server
	// (启动服务器)
	log.Println("Starting server with context-based middleware...")
	s.Serve()
}
