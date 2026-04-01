package cinaclaw

import (
	"context"
	"fmt"
	"sync"
)

// ClientManager manages gRPC connections for multiple users/tenants.
// Each user gets their own CinaClawClient connection (possibly to different sockets).
type ClientManager struct {
	clients           sync.Map // map[string]*CinaClawClient
	defaultSocketPath string
	mu                sync.Mutex
}

// NewClientManager creates a new client manager with the default socket path.
func NewClientManager(socketPath string) *ClientManager {
	return &ClientManager{
		defaultSocketPath: socketPath,
	}
}

// GetClient returns (or creates) a gRPC client for the given user.
// It uses the default socket path unless a user-specific path is provided.
func (m *ClientManager) GetClient(userID string) (*CinaClawClient, error) {
	if val, ok := m.clients.Load(userID); ok {
		client := val.(*CinaClawClient)
		// Verify the connection is still alive
		if err := client.Ping(context.Background()); err == nil {
			return client, nil
		}
		// Connection dead, remove and recreate
		m.clients.Delete(userID)
		client.Close()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring lock
	if val, ok := m.clients.Load(userID); ok {
		client := val.(*CinaClawClient)
		if err := client.Ping(context.Background()); err == nil {
			return client, nil
		}
		m.clients.Delete(userID)
		client.Close()
	}

	client, err := NewClient(m.defaultSocketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for user %s: %w", userID, err)
	}

	m.clients.Store(userID, client)
	return client, nil
}

// GetClientForSocket creates or returns a client connected to a specific socket path.
func (m *ClientManager) GetClientForSocket(userID, socketPath string) (*CinaClawClient, error) {
	key := fmt.Sprintf("%s@%s", userID, socketPath)

	if val, ok := m.clients.Load(key); ok {
		client := val.(*CinaClawClient)
		if err := client.Ping(context.Background()); err == nil {
			return client, nil
		}
		m.clients.Delete(key)
		client.Close()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	client, err := NewClient(socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for user %s at %s: %w", userID, socketPath, err)
	}

	m.clients.Store(key, client)
	return client, nil
}

// RemoveClient closes and removes the client for the given user.
func (m *ClientManager) RemoveClient(userID string) error {
	if val, ok := m.clients.LoadAndDelete(userID); ok {
		client := val.(*CinaClawClient)
		return client.Close()
	}
	return nil
}

// CloseAll closes all managed client connections.
func (m *ClientManager) CloseAll() error {
	var lastErr error
	m.clients.Range(func(key, value interface{}) bool {
		client := value.(*CinaClawClient)
		if err := client.Close(); err != nil {
			lastErr = err
		}
		m.clients.Delete(key)
		return true
	})
	return lastErr
}

// ActiveClients returns the number of active client connections.
func (m *ClientManager) ActiveClients() int {
	count := 0
	m.clients.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}
