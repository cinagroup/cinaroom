package pool

import (
	"errors"
	"io"
	"log/slog"
	"sync"
	"sync/atomic"
)

// WSConn is a minimal interface that satisfies both gorilla/websocket and
// nhooyr.io/websocket Conn types for pool purposes.
type WSConn interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
}

// ErrUserOffline is returned when trying to send to a user with no connection.
var ErrUserOffline = errors.New("pool: user is offline")

const (
	// TextMessage denotes a text data message.
	TextMessage = 1
	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2
)

// wsEntry wraps a single WebSocket connection.
type wsEntry struct {
	conn WSConn
	mu   sync.Mutex // protect concurrent writes to this connection
}

// WSPool manages WebSocket connections keyed by user ID.
// It is safe for concurrent use via sync.Map.
type WSPool struct {
	connections sync.Map // map[string]*wsEntry
	count       atomic.Int64
}

// NewWSPool creates a new WebSocket connection pool.
func NewWSPool() *WSPool {
	return &WSPool{}
}

// Register adds a WebSocket connection for a user. If the user already has a
// connection it is closed and replaced.
func (p *WSPool) Register(userID string, conn WSConn) {
	entry := &wsEntry{conn: conn}
	if old, loaded := p.connections.Swap(userID, entry); loaded {
		if oldEntry, ok := old.(*wsEntry); ok {
			_ = oldEntry.conn.Close()
		}
	} else {
		p.count.Add(1)
	}
}

// Unregister removes and closes the WebSocket connection for a user.
func (p *WSPool) Unregister(userID string) {
	if old, loaded := p.connections.LoadAndDelete(userID); loaded {
		if oldEntry, ok := old.(*wsEntry); ok {
			_ = oldEntry.conn.Close()
		}
		p.count.Add(-1)
	}
}

// Send transmits a message to a specific user.
func (p *WSPool) Send(userID string, msg []byte) error {
	val, ok := p.connections.Load(userID)
	if !ok {
		return ErrUserOffline
	}
	entry := val.(*wsEntry)
	entry.mu.Lock()
	defer entry.mu.Unlock()
	if err := entry.conn.WriteMessage(TextMessage, msg); err != nil {
		if isClosedErr(err) {
			p.Unregister(userID)
		}
		return err
	}
	return nil
}

// Broadcast sends a message to all connected users.
func (p *WSPool) Broadcast(msg []byte) {
	p.connections.Range(func(key, val interface{}) bool {
		entry := val.(*wsEntry)
		entry.mu.Lock()
		if err := entry.conn.WriteMessage(TextMessage, msg); err != nil {
			slog.Warn("ws broadcast write failed", "user", key, "error", err)
			entry.mu.Unlock()
			p.Unregister(key.(string))
			return true
		}
		entry.mu.Unlock()
		return true
	})
}

// GetUserCount returns the number of currently connected users.
func (p *WSPool) GetUserCount() int {
	return int(p.count.Load())
}

// CloseAll closes every connection in the pool.
func (p *WSPool) CloseAll() {
	p.connections.Range(func(key, val interface{}) bool {
		entry := val.(*wsEntry)
		_ = entry.conn.Close()
		p.connections.Delete(key)
		return true
	})
	p.count.Store(0)
}

// isClosedErr checks if the error indicates a closed connection.
func isClosedErr(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrClosedPipe)
}
