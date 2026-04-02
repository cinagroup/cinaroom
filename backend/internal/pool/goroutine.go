package pool

import (
	"errors"
	"sync"
)

// ErrPoolClosed is returned when submitting a task to a stopped pool.
var ErrPoolClosed = errors.New("pool: goroutine pool is closed")

// GoroutinePool limits the number of concurrent goroutines using a
// buffered channel as a task queue.
type GoroutinePool struct {
	workers   int
	taskQueue chan func()
	wg        sync.WaitGroup
	closeOnce sync.Once
	stopCh    chan struct{}
}

// NewGoroutinePool creates a pool with the given number of workers.
// If maxWorkers <= 0, it defaults to 100.
func NewGoroutinePool(maxWorkers int) *GoroutinePool {
	if maxWorkers <= 0 {
		maxWorkers = 100
	}
	p := &GoroutinePool{
		workers:   maxWorkers,
		taskQueue: make(chan func(), maxWorkers*2),
		stopCh:    make(chan struct{}),
	}
	p.start()
	return p
}

// start launches the worker goroutines.
func (p *GoroutinePool) start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// worker reads tasks from the queue and executes them.
func (p *GoroutinePool) worker() {
	defer p.wg.Done()
	for task := range p.taskQueue {
		if task == nil {
			return
		}
		safeExec(task)
	}
}

// Submit enqueues a task for execution. Returns ErrPoolClosed if the pool
// has been stopped.
func (p *GoroutinePool) Submit(fn func()) error {
	if fn == nil {
		return nil
	}
	select {
	case <-p.stopCh:
		return ErrPoolClosed
	default:
	}
	select {
	case p.taskQueue <- fn:
		return nil
	case <-p.stopCh:
		return ErrPoolClosed
	}
}

// Stop gracefully shuts down the pool. It signals workers to finish and
// waits for all in-flight tasks to complete.
func (p *GoroutinePool) Stop() {
	p.closeOnce.Do(func() {
		close(p.stopCh)
		// Send nil sentinels to wake all workers so they exit the range loop
		// after draining remaining real tasks.
		for i := 0; i < p.workers; i++ {
			p.taskQueue <- nil
		}
		p.wg.Wait()
	})
}

// QueueLen returns the current number of pending tasks in the queue.
func (p *GoroutinePool) QueueLen() int {
	return len(p.taskQueue)
}

// safeExec runs fn and recovers from any panic.
func safeExec(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			// Log but do not crash the worker.
			_ = r
		}
	}()
	fn()
}
