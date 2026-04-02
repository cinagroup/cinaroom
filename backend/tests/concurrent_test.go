package tests

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/pool"
)

// ─── Concurrent VM operations ───────────────────────────────────────────────

func TestConcurrentVMOperations(t *testing.T) {
	const goroutines = 50

	var (
		wg       sync.WaitGroup
		errors   atomic.Int64
		created  atomic.Int64
	)

	// Simulate concurrent VM "creations" via the GoroutinePool
	p := pool.NewGoroutinePool(goroutines)
	defer p.Stop()

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		vmName := fmt.Sprintf("test-vm-%d", i)
		err := p.Submit(func() {
			defer wg.Done()
			// Simulate a VM creation step
			time.Sleep(time.Millisecond * 10)
			if vmName == "" {
				errors.Add(1)
				return
			}
			created.Add(1)
		})
		if err != nil {
			t.Errorf("Submit failed: %v", err)
			wg.Done()
		}
	}

	wg.Wait()

	if errors.Load() > 0 {
		t.Errorf("got %d errors during concurrent VM creation", errors.Load())
	}
	if created.Load() != goroutines {
		t.Errorf("expected %d created, got %d", goroutines, created.Load())
	}
}

// ─── Concurrent WebSocket pool register/unregister ─────────────────────────

// mockWSConn is a simple mock WSConn for testing.
type mockWSConn struct {
	closed int32
}

func (m *mockWSConn) WriteMessage(messageType int, data []byte) error {
	if atomic.LoadInt32(&m.closed) == 1 {
		return fmt.Errorf("connection closed")
	}
	return nil
}

func (m *mockWSConn) Close() error {
	atomic.StoreInt32(&m.closed, 1)
	return nil
}

func TestConcurrentWebSocketConnections(t *testing.T) {
	wp := pool.NewWSPool()
	const numConns = 100

	var wg sync.WaitGroup

	// Concurrently register connections
	for i := 0; i < numConns; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			userID := fmt.Sprintf("user-%d", idx)
			conn := &mockWSConn{}
			wp.Register(userID, conn)
		}(i)
	}
	wg.Wait()

	count := wp.GetUserCount()
	if count != numConns {
		t.Errorf("expected %d connections, got %d", numConns, count)
	}

	// Concurrently unregister half
	for i := 0; i < numConns/2; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			userID := fmt.Sprintf("user-%d", idx)
			wp.Unregister(userID)
		}(i)
	}
	wg.Wait()

	count = wp.GetUserCount()
	expected := numConns - numConns/2
	if int(count) != expected {
		t.Errorf("expected %d connections after unregister, got %d", expected, count)
	}

	// Verify remaining users can still send
	for i := numConns / 2; i < numConns; i++ {
		userID := fmt.Sprintf("user-%d", i)
		err := wp.Send(userID, []byte("hello"))
		if err != nil {
			t.Errorf("Send to %s failed: %v", userID, err)
		}
	}

	wp.CloseAll()
}

func TestConcurrentWSPool_RegisterReplace(t *testing.T) {
	wp := pool.NewWSPool()
	defer wp.CloseAll()

	const iterations = 50
	userID := "same-user"
	var wg sync.WaitGroup

	// All goroutines register the same user — only the last one should win
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wp.Register(userID, &mockWSConn{})
		}()
	}
	wg.Wait()

	count := wp.GetUserCount()
	if count != 1 {
		t.Errorf("expected 1 connection for same user, got %d", count)
	}
}

func TestConcurrentWSPool_SendOffline(t *testing.T) {
	wp := pool.NewWSPool()
	defer wp.CloseAll()

	err := wp.Send("nonexistent", []byte("hello"))
	if err != pool.ErrUserOffline {
		t.Errorf("expected ErrUserOffline, got %v", err)
	}
}

func TestConcurrentWSPool_Broadcast(t *testing.T) {
	wp := pool.NewWSPool()

	// Register several connections
	for i := 0; i < 10; i++ {
		userID := fmt.Sprintf("broadcast-user-%d", i)
		wp.Register(userID, &mockWSConn{})
	}

	wp.Broadcast([]byte("broadcast test"))

	count := wp.GetUserCount()
	if count != 10 {
		t.Errorf("expected 10 connections after broadcast, got %d", count)
	}

	wp.CloseAll()

	if wp.GetUserCount() != 0 {
		t.Errorf("expected 0 connections after CloseAll, got %d", wp.GetUserCount())
	}
}

// ─── Connection pool stress tests ───────────────────────────────────────────

func TestGoroutinePoolStress(t *testing.T) {
	const workers = 100
	const tasks = 1000

	p := pool.NewGoroutinePool(workers)
	defer p.Stop()

	var executed atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < tasks; i++ {
		wg.Add(1)
		err := p.Submit(func() {
			defer wg.Done()
			time.Sleep(time.Microsecond * 100) // simulate work
			executed.Add(1)
		})
		if err != nil {
			wg.Done()
			t.Errorf("Submit failed at task %d: %v", i, err)
		}
	}

	wg.Wait()

	if executed.Load() != tasks {
		t.Errorf("expected %d executions, got %d", tasks, executed.Load())
	}
}

func TestGoroutinePool_StopRejectsNew(t *testing.T) {
	p := pool.NewGoroutinePool(5)

	// Submit a few tasks first
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		_ = p.Submit(func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
		})
	}
	wg.Wait()

	p.Stop()

	err := p.Submit(func() {})
	if err != pool.ErrPoolClosed {
		t.Errorf("expected ErrPoolClosed after Stop, got %v", err)
	}
}

func TestGoroutinePool_PanicRecovery(t *testing.T) {
	p := pool.NewGoroutinePool(5)
	defer p.Stop()

	var wg sync.WaitGroup

	// Submit a task that panics
	wg.Add(1)
	err := p.Submit(func() {
		defer wg.Done()
		panic("test panic")
	})
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()

	// Pool should still work after panic
	wg.Add(1)
	err = p.Submit(func() {
		defer wg.Done()
		// normal task
	})
	if err != nil {
		t.Errorf("pool should accept tasks after panic recovery: %v", err)
	}
	wg.Wait()
}

func TestGoroutinePool_NilTask(t *testing.T) {
	p := pool.NewGoroutinePool(5)
	defer p.Stop()

	err := p.Submit(nil)
	if err != nil {
		t.Errorf("nil task should be accepted silently, got %v", err)
	}
}

// ─── GRPCPool tests (without real connections) ──────────────────────────────

func TestGRPCPool_Creation(t *testing.T) {
	p := pool.NewGRPCPool("/tmp/test.sock", 5)
	defer p.Close()

	// Pool should be created without error
	if p == nil {
		t.Fatal("expected non-nil pool")
	}
}

func TestGRPCPool_ClosedRejects(t *testing.T) {
	p := pool.NewGRPCPool("/tmp/test.sock", 5)
	p.Close()

	_, err := p.Get()
	if err == nil {
		t.Error("expected error when getting from closed pool")
	}
}

func TestGRPCPool_PutNil(t *testing.T) {
	p := pool.NewGRPCPool("/tmp/test.sock", 5)
	defer p.Close()

	// Putting nil should not panic
	p.Put(nil)
}
