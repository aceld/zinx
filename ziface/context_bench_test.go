package ziface

import (
	"fmt"
	"sync"
	"testing"
)

func BenchmarkContext_SetGet(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("key", "value")
		c.Get("key")
	}
}

func BenchmarkContext_SetOnly(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("key", i)
	}
}

func BenchmarkContext_GetOnly(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()
	c.Set("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get("key")
	}
}

func BenchmarkContext_Concurrent(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%100)
			c.Set(key, i)
			c.Get(key)
			i++
		}
	})
}

func BenchmarkContext_ConcurrentSet(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			c.Set(fmt.Sprintf("key%d", i%100), i)
			i++
		}
	})
}

func BenchmarkContext_ConcurrentGet(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	// Pre-populate keys
	for i := 0; i < 100; i++ {
		c.Set(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			c.Get(fmt.Sprintf("key%d", i%100))
			i++
		}
	})
}

func BenchmarkContext_Pool(b *testing.B) {
	conn := &mockConnection{connID: 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := NewContext(conn, uint32(i), []byte("test"))
		c.Set("key", i)
		c.Release()
	}
}

func BenchmarkContext_PoolParallel(b *testing.B) {
	conn := &mockConnection{connID: 1}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			c := NewContext(conn, uint32(i), []byte("test"))
			c.Set("key", i)
			c.Release()
			i++
		}
	})
}

func BenchmarkContext_Copy(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	// Populate with some keys
	for i := 0; i < 10; i++ {
		c.Set(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cp := c.Copy()
		cp.Release()
	}
}

func BenchmarkContext_Next(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []HandlerFunc{
		func(c *Context) { c.Next() },
		func(c *Context) { c.Next() },
		func(c *Context) { c.Next() },
		func(c *Context) { c.Next() },
		func(c *Context) {},
	}
	c.SetHandlers(handlers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.index = -1
		c.Next()
	}
}

func BenchmarkContext_NextParallel(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []HandlerFunc{
		func(c *Context) { c.Next() },
		func(c *Context) { c.Next() },
		func(c *Context) { c.Next() },
		func(c *Context) { c.Next() },
		func(c *Context) {},
	}
	c.SetHandlers(handlers)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.index = -1
			c.Next()
		}
	})
}

func BenchmarkContext_MultipleKeys(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	numKeys := 100
	keys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			c.Set(key, i)
		}
		for _, key := range keys {
			c.Get(key)
		}
	}
}

func BenchmarkContext_MultipleKeysConcurrent(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	numKeys := 100
	keys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			for _, key := range keys {
				c.Set(key, i)
			}
			for _, key := range keys {
				c.Get(key)
			}
			i++
		}
	})
}

func BenchmarkContext_Abort(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []HandlerFunc{
		func(c *Context) { c.Abort() },
		func(c *Context) {},
		func(c *Context) {},
	}
	c.SetHandlers(handlers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.index = -1
		c.Next()
	}
}

func BenchmarkContext_IsAborted(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.IsAborted()
	}
}

func BenchmarkContext_SetHandlers(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []HandlerFunc{
		func(c *Context) {},
		func(c *Context) {},
		func(c *Context) {},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.SetHandlers(handlers)
	}
}

func BenchmarkContext_GetHandlers(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	handlers := []HandlerFunc{
		func(c *Context) {},
		func(c *Context) {},
		func(c *Context) {},
	}
	c.SetHandlers(handlers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetHandlers()
	}
}

func BenchmarkContext_MustGet(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()
	c.Set("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.MustGet("key")
	}
}

func BenchmarkContext_MustGetPanic(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			defer func() { recover() }()
			c.MustGet("nonexistent")
		}()
	}
}

func BenchmarkContext_Reset(b *testing.B) {
	conn := &mockConnection{connID: 1}
	c := NewContext(conn, 1, []byte("test"))
	defer c.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("key", i)
		c.Reset()
	}
}

func BenchmarkContext_Release(b *testing.B) {
	conn := &mockConnection{connID: 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := NewContext(conn, uint32(i), []byte("test"))
		c.Set("key", i)
		c.Release()
	}
}

func BenchmarkContext_ConcurrentPool(b *testing.B) {
	conn := &mockConnection{connID: 1}

	var wg sync.WaitGroup
	numGoroutines := 10

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < numGoroutines; j++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				c := NewContext(conn, uint32(id), []byte("test"))
				c.Set("key", id)
				c.Release()
			}(j)
		}
		wg.Wait()
	}
}
