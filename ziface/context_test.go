package ziface

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

// mockConnection is a mock implementation of IConnection for testing
type mockConnection struct {
	connID uint64
}

func (m *mockConnection) Start()                     {}
func (m *mockConnection) Stop()                      {}
func (m *mockConnection) Context() context.Context   { return context.Background() }
func (m *mockConnection) GetName() string            { return "mock" }
func (m *mockConnection) GetConnection() net.Conn    { return nil }
func (m *mockConnection) GetWsConn() *websocket.Conn { return nil }
func (m *mockConnection) GetTCPConnection() net.Conn { return nil }
func (m *mockConnection) GetConnID() uint64          { return m.connID }
func (m *mockConnection) GetConnIdStr() string       { return fmt.Sprintf("%d", m.connID) }
func (m *mockConnection) GetMsgHandler() IMsgHandle  { return nil }
func (m *mockConnection) GetWorkerID() uint32        { return 0 }
func (m *mockConnection) RemoteAddr() net.Addr       { return nil }
func (m *mockConnection) LocalAddr() net.Addr        { return nil }
func (m *mockConnection) LocalAddrString() string    { return "" }
func (m *mockConnection) RemoteAddrString() string   { return "" }
func (m *mockConnection) Send(data []byte) error     { return nil }
func (m *mockConnection) SendToQueue(data []byte, opts ...MsgSendOption) error {
	return nil
}
func (m *mockConnection) SendMsg(msgID uint32, data []byte) error { return nil }
func (m *mockConnection) SendBuffMsg(msgID uint32, data []byte, opts ...MsgSendOption) error {
	return nil
}
func (m *mockConnection) SetProperty(key string, value interface{})                  {}
func (m *mockConnection) GetProperty(key string) (interface{}, error)                { return nil, nil }
func (m *mockConnection) RemoveProperty(key string)                                  {}
func (m *mockConnection) IsAlive() bool                                              { return true }
func (m *mockConnection) SetHeartBeat(checker IHeartbeatChecker)                     {}
func (m *mockConnection) AddCloseCallback(handler, key interface{}, callback func()) {}
func (m *mockConnection) RemoveCloseCallback(handler, key interface{})               {}
func (m *mockConnection) InvokeCloseCallbacks()                                      {}

func TestContext_SetGet(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	// Test Set and Get
	c.Set("key1", "value1")
	c.Set("key2", 42)
	c.Set("key3", true)

	val1, ok1 := c.Get("key1")
	if !ok1 || val1 != "value1" {
		t.Errorf("Expected value1, got %v", val1)
	}

	val2, ok2 := c.Get("key2")
	if !ok2 || val2 != 42 {
		t.Errorf("Expected 42, got %v", val2)
	}

	val3, ok3 := c.Get("key3")
	if !ok3 || val3 != true {
		t.Errorf("Expected true, got %v", val3)
	}

	// Test non-existent key
	val4, ok4 := c.Get("nonexistent")
	if ok4 {
		t.Errorf("Expected false for non-existent key, got true with value %v", val4)
	}
}

func TestContext_MustGet(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	c.Set("key1", "value1")

	// Test MustGet with existing key
	val := c.MustGet("key1")
	if val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	// Test MustGet with non-existent key (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for non-existent key, but didn't panic")
		}
	}()
	c.MustGet("nonexistent")
}

func TestContext_Next(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	executionOrder := []int{}
	handlers := []HandlerFunc{
		func(c *Context) {
			executionOrder = append(executionOrder, 1)
			c.Next()
		},
		func(c *Context) {
			executionOrder = append(executionOrder, 2)
			c.Next()
		},
		func(c *Context) {
			executionOrder = append(executionOrder, 3)
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	if len(executionOrder) != 3 {
		t.Errorf("Expected 3 handlers to execute, got %d", len(executionOrder))
	}

	for i, v := range executionOrder {
		if v != i+1 {
			t.Errorf("Expected execution order %d, got %d", i+1, v)
		}
	}
}

func TestContext_Abort(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	executionOrder := []int{}
	handlers := []HandlerFunc{
		func(c *Context) {
			executionOrder = append(executionOrder, 1)
			c.Abort()
		},
		func(c *Context) {
			executionOrder = append(executionOrder, 2)
		},
		func(c *Context) {
			executionOrder = append(executionOrder, 3)
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	if len(executionOrder) != 1 {
		t.Errorf("Expected 1 handler to execute, got %d", len(executionOrder))
	}

	if executionOrder[0] != 1 {
		t.Errorf("Expected handler 1 to execute, got %d", executionOrder[0])
	}

	if !c.IsAborted() {
		t.Error("Expected context to be aborted")
	}
}

func TestContext_Concurrent(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				c.Set(key, id*numOperations+j)
			}
		}(i)
	}

	wg.Wait()

	// Verify all values
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numOperations; j++ {
			key := fmt.Sprintf("key_%d_%d", i, j)
			expected := i*numOperations + j
			val, ok := c.Get(key)
			if !ok {
				t.Errorf("Key %s not found", key)
			} else if val != expected {
				t.Errorf("Key %s: expected %d, got %v", key, expected, val)
			}
		}
	}
}

func TestContext_ConcurrentReadWrite(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	// Initialize some keys
	for i := 0; i < 10; i++ {
		c.Set(fmt.Sprintf("key%d", i), i)
	}

	var wg sync.WaitGroup
	numReaders := 50
	numWriters := 50

	// Concurrent readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Get(fmt.Sprintf("key%d", j%10))
			}
		}()
	}

	// Concurrent writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Set(fmt.Sprintf("key%d", j%10), id*100+j)
			}
		}(i)
	}

	wg.Wait()
}

func TestContext_Copy(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	c.Set("key1", "value1")
	c.Set("key2", 42)

	cp := c.Copy()

	// Verify copy has same values
	val1, ok1 := cp.Get("key1")
	if !ok1 || val1 != "value1" {
		t.Errorf("Copy: Expected value1, got %v", val1)
	}

	val2, ok2 := cp.Get("key2")
	if !ok2 || val2 != 42 {
		t.Errorf("Copy: Expected 42, got %v", val2)
	}

	// Modify original, copy should not change
	c.Set("key1", "modified")
	val1Orig, _ := c.Get("key1")
	val1Copy, _ := cp.Get("key1")

	if val1Orig != "modified" {
		t.Errorf("Original: Expected modified, got %v", val1Orig)
	}

	if val1Copy != "value1" {
		t.Errorf("Copy: Expected value1, got %v", val1Copy)
	}

	cp.Release()
}

func TestContext_Pool(t *testing.T) {
	// Test that contexts are reused from pool
	conn := &mockConnection{connID: 1}

	// Create and release multiple contexts
	for i := 0; i < 100; i++ {
		c := NewContext(conn, uint32(i), []byte("test"))
		c.Set("key", i)
		c.Release()
	}

	// Create a new context and verify it's reset
	c := NewContext(conn, 999, []byte("test"))
	defer c.Release()

	// Keys should be cleared
	val, ok := c.Get("key")
	if ok {
		t.Errorf("Expected no value, got %v", val)
	}

	// New values should work
	c.Set("newkey", "newvalue")
	val, ok = c.Get("newkey")
	if !ok || val != "newvalue" {
		t.Errorf("Expected newvalue, got %v", val)
	}
}

func TestContext_Handlers(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []HandlerFunc{
		func(c *Context) {},
		func(c *Context) {},
	}

	c.SetHandlers(handlers)
	retrieved := c.GetHandlers()

	if len(retrieved) != len(handlers) {
		t.Errorf("Expected %d handlers, got %d", len(handlers), len(retrieved))
	}
}
