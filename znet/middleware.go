// @Title middleware.go
// @Description Provides common middleware implementations for zinx
// @Author Aceld - Upgrade for context/OTEL support
package znet

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
)

// RecoveryMiddleware returns a middleware that recovers from panics and logs the error
// (返回一个从panic中恢复并记录错误的中间件)
func RecoveryMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				// (记录panic和堆栈信息)
				zlog.Ins().ErrorF("Panic recovered: %v\n%s", err, debug.Stack())

				// Abort the middleware chain
				// (中止中间件链)
				c.Abort()
			}
		}()

		// Continue to the next middleware
		// (继续执行下一个中间件)
		c.Next()
	}
}

// LoggingMiddleware returns a middleware that logs request information
// (返回一个记录请求信息的中间件)
func LoggingMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		start := time.Now()

		// Continue to the next middleware
		// (继续执行下一个中间件)
		c.Next()

		// Log request information after processing
		// (处理后记录请求信息)
		latency := time.Since(start)
		zlog.Ins().InfoF("MessageID: %d, ConnectionID: %d, Latency: %v",
			c.MsgID, c.Conn.GetConnID(), latency)
	}
}

// AuthMiddleware returns a middleware that performs authentication
// (返回一个执行认证的中间件)
func AuthMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		// Example: Check for authentication token in context
		// (示例：检查context中的认证token)
		token, exists := c.Get("auth_token")
		if !exists {
			zlog.Ins().ErrorF("Authentication failed: no token found")
			c.Abort()
			return
		}

		// Validate token (simplified example)
		// (验证token（简化示例）)
		if token == "" {
			zlog.Ins().ErrorF("Authentication failed: invalid token")
			c.Abort()
			return
		}

		// Continue to the next middleware
		// (继续执行下一个中间件)
		c.Next()
	}
}

// RateLimitMiddleware returns a middleware that performs rate limiting
// (返回一个执行限流的中间件)
func RateLimitMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		// Example: Simple rate limiting based on connection ID
		// (示例：基于连接ID的简单限流)
		connID := c.Conn.GetConnID()

		// In a real implementation, you would use a more sophisticated rate limiting algorithm
		// (在实际实现中，您会使用更复杂的限流算法)
		log.Printf("Rate limit check for connection %d", connID)

		// Continue to the next middleware
		// (继续执行下一个中间件)
		c.Next()
	}
}

// TraceMiddleware returns a middleware that adds tracing information
// Note: This is a simplified example. For real OTEL integration, you would use the OpenTelemetry SDK
// (返回一个添加追踪信息的中间件)
// 注意：这是一个简化示例。对于真正的OTEL集成，您需要使用OpenTelemetry SDK
func TraceMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		// In a real implementation, you would:
		// 1. Start a new span
		// 2. Add trace ID to context
		// 3. Log trace information
		// (在实际实现中，您需要：)
		// 1. 启动一个新的span
		// 2. 将trace ID添加到context
		// 3. 记录追踪信息)

		// Example: Add a trace ID to the context
		// (示例：向context添加trace ID)
		traceID := fmt.Sprintf("trace-%d-%d", c.Conn.GetConnID(), c.MsgID)
		c.Set("trace_id", traceID)

		zlog.Ins().DebugF("Trace started: %s", traceID)

		// Continue to the next middleware
		// (继续执行下一个中间件)
		c.Next()

		// After processing, you could end the span here
		// (处理后，您可以在这里结束span)
		zlog.Ins().DebugF("Trace ended: %s", traceID)
	}
}

// CORSMiddleware returns a middleware that handles CORS headers
// Note: This is more relevant for HTTP servers, but included for completeness
// (返回一个处理CORS头的中间件)
// 注意：这更适用于HTTP服务器，但为了完整性而包含)
func CORSMiddleware() ziface.HandlerFunc {
	return func(c *ziface.Context) {
		// Set CORS headers (simplified example)
		// (设置CORS头（简化示例）)
		// In a real implementation, you would set appropriate headers
		// (在实际实现中，您需要设置适当的头)

		// Continue to the next middleware
		// (继续执行下一个中间件)
		c.Next()
	}
}

// TimeoutMiddleware returns a middleware that adds timeout handling
// (返回一个添加超时处理的中间件)
func TimeoutMiddleware(timeout time.Duration) ziface.HandlerFunc {
	return func(c *ziface.Context) {
		// Create a channel to signal completion
		// (创建一个信号完成的channel)
		done := make(chan struct{})

		go func() {
			// Execute the middleware chain
			// (执行中间件链)
			c.Next()
			close(done)
		}()

		// Wait for completion or timeout
		// (等待完成或超时)
		select {
		case <-done:
			// Completed successfully
			// (成功完成)
		case <-time.After(timeout):
			// Timeout occurred
			// (发生超时)
			zlog.Ins().ErrorF("Request timeout after %v", timeout)
			c.Abort()
		}
	}
}
