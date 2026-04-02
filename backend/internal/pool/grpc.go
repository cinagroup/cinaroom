package pool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCPool manages a pool of gRPC client connections over a Unix socket.
type GRPCPool struct {
	socketPath string
	maxConns   int
	conns      chan *grpc.ClientConn
	mu         sync.Mutex
	closed     atomic.Bool
	dialOpts   []grpc.DialOption
}

// NewGRPCPool creates a connection pool for the given Unix socket path.
// maxConns controls the pool size (default 10 if <= 0).
func NewGRPCPool(socketPath string, maxConns int) *GRPCPool {
	if maxConns <= 0 {
		maxConns = 10
	}
	p := &GRPCPool{
		socketPath: socketPath,
		maxConns:   maxConns,
		conns:      make(chan *grpc.ClientConn, maxConns),
		dialOpts: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		},
	}
	return p
}

// Get retrieves a connection from the pool or dials a new one if needed.
func (p *GRPCPool) Get() (*grpc.ClientConn, error) {
	if p.closed.Load() {
		return nil, fmt.Errorf("pool: grpc pool is closed")
	}

	// Try to reuse an existing connection.
	select {
	case conn := <-p.conns:
		if conn != nil && conn.GetState().String() != "SHUTDOWN" {
			return conn, nil
		}
		// Connection is dead; dial a new one.
		_ = conn.Close()
	default:
	}

	// Dial a new connection.
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "unix://"+p.socketPath, p.dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("pool: grpc dial failed: %w", err)
	}
	return conn, nil
}

// Put returns a connection to the pool for reuse.
// If the pool is closed or full, the connection is closed immediately.
func (p *GRPCPool) Put(conn *grpc.ClientConn) {
	if conn == nil {
		return
	}
	if p.closed.Load() {
		_ = conn.Close()
		return
	}
	select {
	case p.conns <- conn:
		// Returned to pool.
	default:
		// Pool is full; close excess connection.
		_ = conn.Close()
	}
}

// Close shuts down all pooled connections.
func (p *GRPCPool) Close() {
	if !p.closed.CompareAndSwap(false, true) {
		return // already closed
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	close(p.conns)
	for conn := range p.conns {
		_ = conn.Close()
	}
}

// ActiveCount returns the number of connections currently in the pool.
func (p *GRPCPool) ActiveCount() int {
	return len(p.conns)
}
