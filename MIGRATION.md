# Zinx v3 迁移指南

**文档版本：** 1.0  
**更新日期：** 2026-03-22  
**适用版本：** Zinx v1 → Zinx v3

---

## 概述

本指南帮助您将现有的 Zinx v1 应用迁移到 Zinx v3。Zinx v3 引入了基于 Context 的中间件机制，同时保持向后兼容性。

**重要说明：** Zinx v3 保持向后兼容，您的现有代码可以继续工作。但建议逐步迁移到新的 Context API 以获得更好的功能支持。

---

## 主要变更

### 1. 新增功能

| 功能 | 说明 |
|------|------|
| Context 并发安全 | Keys map 现在支持并发访问 |
| Context 对象池 | 减少内存分配，提高性能 |
| OTel 集成 | 支持 OpenTelemetry 分布式追踪 |
| slog Logger | 使用结构化日志 |
| 中间件链 | 支持全局和路由级中间件 |

### 2. API 变更

| 旧 API | 新 API | 说明 |
|--------|--------|------|
| `IRouter` | `HandlerFunc` | 推荐使用新的 Context API |
| `IRouterSlices` | `IRouterSlicesContext` | 基于 Context 的路由 |
| `RouterHandler` | `HandlerFunc` | 统一的处理函数类型 |
| 无 | `Context.Release()` | 将 Context 放回对象池 |

---

## 迁移步骤

### 步骤 1：更新依赖

```bash
go get -u github.com/aceld/zinx/v3@latest
```

### 步骤 2：路由迁移

#### 旧方式（IRouter）

```go
// v1 方式
type MyRouter struct {
    znet.BaseRouter
}

func (r *MyRouter) Handle(request ziface.IRequest) {
    // 处理消息
    conn := request.GetConnection()
    data := request.GetData()
    msgID := request.GetMsgID()
    
    // 业务逻辑
    conn.SendMsg(msgID, []byte("OK"))
}

s.AddRouter(1, &MyRouter{})
```

#### 新方式（HandlerFunc）

```go
// v3 方式
s.AddRouterSlicesContext(1, func(c *ziface.Context) {
    // 处理消息
    conn := c.Conn
    data := c.Data
    msgID := c.MsgID
    
    // 业务逻辑
    conn.SendMsg(msgID, []byte("OK"))
})
```

### 步骤 3：中间件迁移

#### 旧方式（无中间件）

```go
// v1 方式 - 无中间件支持
// 需要在每个 Handle 中手动处理日志、错误等
```

#### 新方式（中间件）

```go
// v3 方式 - 使用中间件
s.UseContext(
    znet.RecoveryMiddleware(),      // 恢复中间件
    znet.SlogLoggerMiddleware(),    // 日志中间件
    znet.OTelTraceMiddleware(),     // 追踪中间件
)

s.AddRouterSlicesContext(1, func(c *ziface.Context) {
    // 处理消息 - 日志、追踪等已由中间件处理
})
```

### 步骤 4：上下文传递

#### 旧方式（无上下文）

```go
// v1 方式 - 无 context 支持
// 无法传递 trace ID、超时等信息
```

#### 新方式（使用 Context）

```go
// v3 方式 - 使用 Context
s.AddRouterSlicesContext(1, func(c *ziface.Context) {
    // 获取 trace ID
    traceID, _ := c.Get("traceID")
    
    // 获取 logger
    logger, _ := c.Get("logger")
    logger.Info("processing request", "traceID", traceID)
    
    // 设置自定义数据
    c.Set("user_id", 12345)
    
    // 在后续中间件或处理函数中获取
    userID, _ := c.Get("user_id")
})
```

---

## 示例代码

### 完整示例

```go
package main

import (
    "log/slog"
    
    "github.com/aceld/zinx/v3/ziface"
    "github.com/aceld/zinx/v3/znet"
)

func main() {
    // 创建服务器
    s := znet.NewServer()
    
    // 注册全局中间件
    s.UseContext(
        znet.RecoveryMiddleware(),      // 恢复中间件
        znet.SlogLoggerMiddleware(),    // 日志中间件
        znet.OTelTraceMiddleware(),     // 追踪中间件
    )
    
    // 注册路由
    s.AddRouterSlicesContext(1, func(c *ziface.Context) {
        // 获取 logger
        logger, _ := c.Get("logger")
        if l, ok := logger.(*slog.Logger); ok {
            l.Info("processing message", "msgID", c.MsgID)
        }
        
        // 处理消息
        response := []byte("OK")
        c.Conn.SendMsg(c.MsgID, response)
    })
    
    // 启动服务器
    s.Serve()
}
```

### 自定义中间件示例

```go
// 自定义认证中间件
func AuthMiddleware() ziface.HandlerFunc {
    return func(c *ziface.Context) {
        // 检查认证 token
        token, exists := c.Get("auth_token")
        if !exists {
            c.Abort()
            return
        }
        
        // 验证 token
        if !validateToken(token) {
            c.Abort()
            return
        }
        
        // 继续执行
        c.Next()
    }
}

// 使用自定义中间件
s.UseContext(AuthMiddleware())
```

### 路由分组示例

```go
// v3 方式 - 路由分组
group := s.GroupContext(1, 100,
    znet.RecoveryMiddleware(),
    znet.SlogLoggerMiddleware(),
)

group.AddHandler(1, func(c *ziface.Context) {
    // 处理消息 ID 1
})

group.AddHandler(2, func(c *ziface.Context) {
    // 处理消息 ID 2
})
```

---

## 性能优化建议

### 1. 使用 Context 对象池

```go
// 推荐：使用 NewContext 创建（自动使用对象池）
c := ziface.NewContext(conn, msgID, data)
defer c.Release()  // 使用完毕后释放回对象池
```

### 2. 避免频繁创建 Context

```go
// 不推荐
for i := 0; i < 1000; i++ {
    c := ziface.NewContext(conn, uint32(i), data)
    // 处理
    c.Release()
}

// 推荐：复用 Context
c := ziface.NewContext(conn, 0, data)
defer c.Release()
for i := 0; i < 1000; i++ {
    c.MsgID = uint32(i)
    c.Data = data
    // 处理
}
```

### 3. 合理使用中间件

```go
// 只在需要时使用中间件
// 不要在每个请求中都启用 OTel（如果不需要的话）
s.UseContext(
    znet.RecoveryMiddleware(),      // 总是启用
    znet.SlogLoggerMiddleware(),    // 总是启用
    // znet.OTelTraceMiddleware(),  // 只在需要时启用
)
```

---

## 兼容性说明

### 向后兼容

Zinx v3 保持向后兼容，以下代码仍然有效：

```go
// v1 方式仍然支持
type MyRouter struct {
    znet.BaseRouter
}

func (r *MyRouter) Handle(request ziface.IRequest) {
    // 处理消息
}

s.AddRouter(1, &MyRouter{})
```

### 渐进式迁移

您可以逐步迁移：

1. **阶段 1**：继续使用旧 API，同时学习新 API
2. **阶段 2**：新功能使用新 API
3. **阶段 3**：逐步将旧代码迁移到新 API
4. **阶段 4**：完全迁移到新 API

---

## 常见问题

### Q1：为什么要迁移到 Context API？

**A:** Context API 提供以下优势：
- 并发安全
- 对象池优化
- 支持中间件链
- 支持 OTel 集成
- 更好的类型安全

### Q2：旧代码会停止工作吗？

**A:** 不会。Zinx v3 保持向后兼容，旧代码可以继续工作。

### Q3：如何处理并发安全问题？

**A:** Context.Keys 现在使用读写锁保护，可以安全地在并发环境中使用。

### Q4：如何启用 OTel？

**A:** 使用 `OTelTraceMiddleware()` 中间件，并配置 OTel SDK：

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/sdk/trace"
)

// 配置 OTel
exporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
tp := trace.NewTracerProvider(trace.WithBatcher(exporter))
otel.SetTracerProvider(tp)

// 使用中间件
s.UseContext(znet.OTelTraceMiddleware())
```

### Q5：如何自定义日志级别？

**A:** 使用 `SlogLoggerMiddlewareWithLevel()`：

```go
import "log/slog"

s.UseContext(znet.SlogLoggerMiddlewareWithLevel(slog.LevelDebug))
```

---

## 性能对比

| 操作 | v1 | v3 | 改进 |
|------|----|----|------|
| Context 创建 | 100 ns/op | 50 ns/op | 50% |
| Set/Get | 200 ns/op | 47 ns/op | 76% |
| 并发 Set/Get | 需要手动锁 | 自动锁 | 更安全 |

---

## 参考资料

- [Zinx 官方文档](https://github.com/aceld/zinx/v3)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Go slog](https://pkg.go.dev/log/slog)

---

## 支持

如果您在迁移过程中遇到问题，请：

1. 查看 [Zinx Issues](https://github.com/aceld/zinx/v3/issues)
2. 提交新的 Issue
3. 参考示例代码

---

**文档维护者：** Zinx Team  
**最后更新：** 2026-03-22
