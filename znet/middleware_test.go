package znet

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/aceld/zinx/ziface"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// mockConnection is a mock implementation of IConnection for testing
type mockConnection struct {
	connID uint64
}

func (m *mockConnection) Start()                           {}
func (m *mockConnection) Stop()                            {}
func (m *mockConnection) Context() context.Context         { return context.Background() }
func (m *mockConnection) GetName() string                  { return "mock" }
func (m *mockConnection) GetConnection() net.Conn          { return nil }
func (m *mockConnection) GetWsConn() *websocket.Conn       { return nil }
func (m *mockConnection) GetTCPConnection() net.Conn       { return nil }
func (m *mockConnection) GetConnID() uint64                { return m.connID }
func (m *mockConnection) GetConnIdStr() string             { return fmt.Sprintf("%d", m.connID) }
func (m *mockConnection) GetMsgHandler() ziface.IMsgHandle { return nil }
func (m *mockConnection) GetWorkerID() uint32              { return 0 }
func (m *mockConnection) RemoteAddr() net.Addr             { return nil }
func (m *mockConnection) LocalAddr() net.Addr              { return nil }
func (m *mockConnection) LocalAddrString() string          { return "" }
func (m *mockConnection) RemoteAddrString() string         { return "" }
func (m *mockConnection) Send(data []byte) error           { return nil }
func (m *mockConnection) SendToQueue(data []byte, opts ...ziface.MsgSendOption) error {
	return nil
}
func (m *mockConnection) SendMsg(msgID uint32, data []byte) error { return nil }
func (m *mockConnection) SendBuffMsg(msgID uint32, data []byte, opts ...ziface.MsgSendOption) error {
	return nil
}
func (m *mockConnection) SetProperty(key string, value interface{})                  {}
func (m *mockConnection) GetProperty(key string) (interface{}, error)                { return nil, nil }
func (m *mockConnection) RemoveProperty(key string)                                  {}
func (m *mockConnection) IsAlive() bool                                              { return true }
func (m *mockConnection) SetHeartBeat(checker ziface.IHeartbeatChecker)              {}
func (m *mockConnection) AddCloseCallback(handler, key interface{}, callback func()) {}
func (m *mockConnection) RemoveCloseCallback(handler, key interface{})               {}
func (m *mockConnection) InvokeCloseCallbacks()                                      {}

func TestRecoveryMiddleware(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	panicked := false
	handlers := []ziface.HandlerFunc{
		RecoveryMiddleware(),
		func(c *ziface.Context) {
			panicked = true
			panic("test panic")
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	if !panicked {
		t.Error("Expected handler to panic")
	}

	if !c.IsAborted() {
		t.Error("Expected context to be aborted after panic")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	executed := false
	handlers := []ziface.HandlerFunc{
		LoggingMiddleware(),
		func(c *ziface.Context) {
			executed = true
			time.Sleep(10 * time.Millisecond)
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	if !executed {
		t.Error("Expected handler to execute")
	}
}

func TestSlogLoggerMiddleware(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	var receivedLogger *slog.Logger
	handlers := []ziface.HandlerFunc{
		SlogLoggerMiddleware(),
		func(c *ziface.Context) {
			if logger, exists := c.Get("logger"); exists {
				if l, ok := logger.(*slog.Logger); ok {
					receivedLogger = l
				}
			}
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	if receivedLogger == nil {
		t.Error("Expected logger to be set in context")
	}
}

func TestSlogLoggerMiddlewareWithTraceID(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	// Set trace ID before middleware
	c.Set("traceID", "test-trace-id")

	var receivedLogger *slog.Logger
	handlers := []ziface.HandlerFunc{
		SlogLoggerMiddleware(),
		func(c *ziface.Context) {
			if logger, exists := c.Get("logger"); exists {
				if l, ok := logger.(*slog.Logger); ok {
					receivedLogger = l
				}
			}
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	if receivedLogger == nil {
		t.Error("Expected logger to be set in context")
	}
}

func TestOTelTraceMiddleware(t *testing.T) {
	// Set up test exporter
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())

	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []ziface.HandlerFunc{
		OTelTraceMiddleware(),
		func(c *ziface.Context) {
			// Handler execution
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	// Verify trace ID was set
	traceID, exists := c.Get("traceID")
	if !exists {
		t.Error("Expected traceID to be set")
	}

	if traceID == "" {
		t.Error("Expected traceID to be non-empty")
	}

	// Verify span ID was set
	spanID, exists := c.Get("spanID")
	if !exists {
		t.Error("Expected spanID to be set")
	}

	if spanID == "" {
		t.Error("Expected spanID to be non-empty")
	}
}

func TestOTelTraceMiddlewareWithError(t *testing.T) {
	// Set up test exporter
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())

	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []ziface.HandlerFunc{
		OTelTraceMiddleware(),
		func(c *ziface.Context) {
			// Set error in context
			c.Set("error", errors.New("test error"))
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	// Verify error was recorded in span
	spans := exporter.GetSpans()
	if len(spans) == 0 {
		t.Error("Expected at least one span")
	}
}

func TestOTelTraceMiddlewareWithName(t *testing.T) {
	// Set up test exporter
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())

	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []ziface.HandlerFunc{
		OTelTraceMiddlewareWithName("custom-tracer"),
		func(c *ziface.Context) {},
	}

	c.SetHandlers(handlers)
	c.Next()

	traceID, exists := c.Get("traceID")
	if !exists || traceID == "" {
		t.Error("Expected traceID to be set")
	}
}

func TestMiddlewareChain(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	executionOrder := []string{}
	handlers := []ziface.HandlerFunc{
		func(c *ziface.Context) {
			executionOrder = append(executionOrder, "middleware1-before")
			c.Next()
			executionOrder = append(executionOrder, "middleware1-after")
		},
		func(c *ziface.Context) {
			executionOrder = append(executionOrder, "middleware2-before")
			c.Next()
			executionOrder = append(executionOrder, "middleware2-after")
		},
		func(c *ziface.Context) {
			executionOrder = append(executionOrder, "handler")
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	expected := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executionOrder) != len(expected) {
		t.Errorf("Expected %d executions, got %d", len(expected), len(executionOrder))
	}

	for i, v := range executionOrder {
		if v != expected[i] {
			t.Errorf("Expected %s at position %d, got %s", expected[i], i, v)
		}
	}
}

func TestMiddlewareAbort(t *testing.T) {
	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	executionOrder := []string{}
	handlers := []ziface.HandlerFunc{
		func(c *ziface.Context) {
			executionOrder = append(executionOrder, "middleware1")
			c.Abort()
		},
		func(c *ziface.Context) {
			executionOrder = append(executionOrder, "middleware2")
		},
		func(c *ziface.Context) {
			executionOrder = append(executionOrder, "handler")
		},
	}

	c.SetHandlers(handlers)
	c.Next()

	if len(executionOrder) != 1 {
		t.Errorf("Expected 1 execution, got %d", len(executionOrder))
	}

	if executionOrder[0] != "middleware1" {
		t.Errorf("Expected middleware1 to execute, got %s", executionOrder[0])
	}
}

func TestMiddlewareConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			conn := &mockConnection{connID: uint64(id)}
			c := ziface.NewContext(conn, uint32(id), []byte("test"))
			defer c.Release()

			handlers := []ziface.HandlerFunc{
				RecoveryMiddleware(),
				SlogLoggerMiddleware(),
				func(c *ziface.Context) {
					c.Set("result", "ok")
				},
			}

			c.SetHandlers(handlers)
			c.Next()

			result, exists := c.Get("result")
			if !exists || result != "ok" {
				t.Errorf("Goroutine %d: Expected result ok, got %v", id, result)
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkRecoveryMiddleware(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []ziface.HandlerFunc{
		RecoveryMiddleware(),
		func(c *ziface.Context) {},
	}

	c.SetHandlers(handlers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Next()
	}
}

func BenchmarkLoggingMiddleware(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []ziface.HandlerFunc{
		LoggingMiddleware(),
		func(c *ziface.Context) {},
	}

	c.SetHandlers(handlers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Next()
	}
}

func BenchmarkSlogLoggerMiddleware(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []ziface.HandlerFunc{
		SlogLoggerMiddleware(),
		func(c *ziface.Context) {},
	}

	c.SetHandlers(handlers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Next()
	}
}

func BenchmarkOTelTraceMiddleware(b *testing.B) {
	// Set up noop tracer
	otel.SetTracerProvider(otel.GetTracerProvider())

	conn := &mockConnection{connID: 1}
	c := ziface.NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []ziface.HandlerFunc{
		OTelTraceMiddleware(),
		func(c *ziface.Context) {},
	}

	c.SetHandlers(handlers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Next()
	}
}
