package znet

import (
	"sync"
	"testing"
	"time"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zpack"
)

type orderRouter struct {
	BaseRouter
	id    int
	delay time.Duration
	mu    *sync.Mutex
	out   *[]int
	done  chan struct{}
}

func (r *orderRouter) Handle(req ziface.IRequest) {
	if r.delay > 0 {
		time.Sleep(r.delay)
	}
	r.mu.Lock()
	*r.out = append(*r.out, r.id)
	if len(*r.out) == 2 {
		close(r.done)
	}
	r.mu.Unlock()
}

func TestClientMsgHandleNoWorkerPoolPreservesOrder(t *testing.T) {
	mh := newCliMsgHandle()

	var mu sync.Mutex
	got := make([]int, 0, 2)
	done := make(chan struct{})

	mh.AddRouter(1, &orderRouter{id: 1, delay: 120 * time.Millisecond, mu: &mu, out: &got, done: done})
	mh.AddRouter(2, &orderRouter{id: 2, delay: 0, mu: &mu, out: &got, done: done})

	req1 := NewRequest(nil, zpack.NewMessageByMsgId(1, 0, nil))
	req2 := NewRequest(nil, zpack.NewMessageByMsgId(2, 0, nil))

	start := time.Now()
	mh.Execute(req1)
	mh.Execute(req2)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for client handlers to complete")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 handled messages, got %d", len(got))
	}
	if got[0] != 1 || got[1] != 2 {
		t.Fatalf("expected client order [1 2], got %v", got)
	}
	if elapsed := time.Since(start); elapsed < 120*time.Millisecond {
		t.Fatalf("expected serialized execution in client mode, elapsed=%v", elapsed)
	}
}

func TestServerMsgHandleNoWorkerPoolRemainsAsync(t *testing.T) {
	mh := newMsgHandle()
	mh.WorkerPoolSize = 0

	var mu sync.Mutex
	got := make([]int, 0, 2)
	done := make(chan struct{})

	mh.AddRouter(1, &orderRouter{id: 1, delay: 120 * time.Millisecond, mu: &mu, out: &got, done: done})
	mh.AddRouter(2, &orderRouter{id: 2, delay: 0, mu: &mu, out: &got, done: done})

	req1 := NewRequest(nil, zpack.NewMessageByMsgId(1, 0, nil))
	req2 := NewRequest(nil, zpack.NewMessageByMsgId(2, 0, nil))

	mh.Execute(req1)
	mh.Execute(req2)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for server handlers to complete")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 handled messages, got %d", len(got))
	}
	if got[0] != 2 || got[1] != 1 {
		t.Fatalf("expected async server order [2 1], got %v", got)
	}
}
