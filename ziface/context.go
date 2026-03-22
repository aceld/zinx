// @Title context.go
// @Description Provides a Gin-like context for middleware chain support
// @Author Aceld - Upgrade for context/OTEL support
package ziface

import (
	"context"
	"math"
	"sync"
)

// HandlerFunc defines the handler function used by the middleware chain
// (定义中间件链使用的处理函数)
type HandlerFunc func(*Context)

// contextPool is a pool for Context objects to reduce memory allocations
// (Context对象池，用于减少内存分配)
var contextPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			Keys: make(map[string]interface{}),
		}
	},
}

// Context is a request context similar to Gin's Context, used for middleware chain
// (类似于Gin的Context，用于中间件链)
type Context struct {
	// Ctx is the standard Go context for request-scoped data like traceID
	// (标准Go context，用于请求范围的数据如traceID)
	Ctx context.Context

	// Conn is the connection associated with this request
	// (与该请求关联的连接)
	Conn IConnection

	// MsgID is the message ID of the request
	// (请求的消息ID)
	MsgID uint32

	// Data is the raw request data
	// (原始请求数据)
	Data []byte

	// handlers is the middleware chain
	// (中间件链)
	handlers []HandlerFunc

	// index is the current position in the middleware chain
	// (中间件链中的当前位置)
	index int8

	// Keys is a map for storing key-value pairs during request processing
	// (用于在请求处理过程中存储键值对的map)
	Keys map[string]interface{}

	// keysLock protects the Keys map for concurrent access
	// (保护Keys map的并发访问)
	keysLock sync.RWMutex
}

// NewContext creates a new Context from the pool
// (从对象池创建一个新的Context)
func NewContext(conn IConnection, msgID uint32, data []byte) *Context {
	c := contextPool.Get().(*Context)
	c.Ctx = context.Background()
	c.Conn = conn
	c.MsgID = msgID
	c.Data = data
	c.index = -1
	// Keys already initialized in pool
	return c
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// (应该只在中间件内部使用。它在调用处理程序中执行链中待处理的处理程序)
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort prevents pending handlers from being called.
// Note that this will not stop the current handler.
// (阻止待处理的处理程序被调用。注意，这不会停止当前处理程序)
func (c *Context) Abort() {
	c.index = math.MaxInt8 / 2
}

// IsAborted returns true if the current context was aborted.
// (如果当前上下文被中止，则返回true)
func (c *Context) IsAborted() bool {
	return c.index >= math.MaxInt8/2
}

// Set is used to store a new key/value pair exclusively for this context.
// Thread-safe: uses write lock to protect concurrent access.
// (用于为此上下文专门存储一个新的键/值对，线程安全：使用写锁保护并发访问)
func (c *Context) Set(key string, value interface{}) {
	c.keysLock.Lock()
	defer c.keysLock.Unlock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// Thread-safe: uses read lock to protect concurrent access.
// (返回给定键的值，线程安全：使用读锁保护并发访问)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.keysLock.RLock()
	defer c.keysLock.RUnlock()
	if c.Keys != nil {
		value, exists = c.Keys[key]
	}
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
// (如果给定的键存在则返回其值，否则panic)
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// Copy returns a copy of the current Context that can be safely used outside the request's scope.
// (返回当前Context的副本，可以在请求范围之外安全使用)
func (c *Context) Copy() *Context {
	cp := contextPool.Get().(*Context)
	cp.Ctx = c.Ctx
	cp.Conn = c.Conn
	cp.MsgID = c.MsgID
	cp.Data = c.Data
	cp.handlers = c.handlers
	cp.index = c.index

	// Copy keys with lock protection
	c.keysLock.RLock()
	cp.keysLock.Lock()
	if c.Keys != nil {
		for k, v := range c.Keys {
			cp.Keys[k] = v
		}
	}
	cp.keysLock.Unlock()
	c.keysLock.RUnlock()

	return cp
}

// SetHandlers sets the handlers chain
// (设置处理程序链)
func (c *Context) SetHandlers(handlers []HandlerFunc) {
	c.handlers = handlers
}

// GetHandlers returns the handlers chain
// (返回处理程序链)
func (c *Context) GetHandlers() []HandlerFunc {
	return c.handlers
}

// Reset resets the context for reuse in the pool
// (重置上下文以便在对象池中重用)
func (c *Context) Reset() {
	c.Ctx = nil
	c.Conn = nil
	c.MsgID = 0
	c.Data = nil
	c.handlers = nil
	c.index = -1
	// Clear Keys map while keeping the underlying memory
	c.keysLock.Lock()
	for k := range c.Keys {
		delete(c.Keys, k)
	}
	c.keysLock.Unlock()
}

// Release returns the context to the pool for reuse
// (将Context放回对象池以便重用)
func (c *Context) Release() {
	c.Reset()
	contextPool.Put(c)
}
