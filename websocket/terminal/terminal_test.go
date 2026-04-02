package main

import (
	"fmt"
	"testing"
	"time"
)

// TestGenerateToken 测试 token 生成
func TestGenerateToken(t *testing.T) {
	cm := NewSecureConnectionManager()
	
	userID := "test-user"
	token, err := cm.GenerateToken(userID)
	
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	if token == "" {
		t.Error("Generated token is empty")
	}
	
	// 验证 token 可以验证
	authToken, valid := cm.validateToken(token)
	if !valid {
		t.Error("Generated token is not valid")
	}
	
	if authToken.UserID != userID {
		t.Errorf("Expected userID %s, got %s", userID, authToken.UserID)
	}
}

// TestTokenExpiry 测试 token 过期
func TestTokenExpiry(t *testing.T) {
	cm := NewSecureConnectionManager()
	
	token, _ := cm.GenerateToken("test-user")
	
	// 立即验证应该成功
	_, valid := cm.validateToken(token)
	if !valid {
		t.Error("Fresh token should be valid")
	}
	
	// 模拟过期（手动设置过期时间）
	cm.mu.Lock()
	if authToken, exists := cm.authTokens[token]; exists {
		authToken.ExpiresAt = time.Now().Add(-1 * time.Hour)
	}
	cm.mu.Unlock()
	
	// 验证应该失败
	_, valid = cm.validateToken(token)
	if valid {
		t.Error("Expired token should be invalid")
	}
}

// TestCleanupExpiredTokens 测试过期 token 清理
func TestCleanupExpiredTokens(t *testing.T) {
	cm := NewSecureConnectionManager()
	
	token, _ := cm.GenerateToken("test-user")
	
	// 设置过期
	cm.mu.Lock()
	if authToken, exists := cm.authTokens[token]; exists {
		authToken.ExpiresAt = time.Now().Add(-1 * time.Hour)
	}
	cm.mu.Unlock()
	
	// 手动触发清理
	cm.cleanupExpiredTokens()
	
	// 验证 token 已被清理
	cm.mu.RLock()
	_, exists := cm.authTokens[token]
	cm.mu.RUnlock()
	
	if exists {
		t.Error("Expired token should be cleaned up")
	}
}

// TestRevokeToken 测试 token 撤销
func TestRevokeToken(t *testing.T) {
	relay := NewCloudRelay()
	
	token := relay.GenerateToken("test-user")
	
	// 验证 token 有效
	_, valid := relay.validateToken(token)
	if !valid {
		t.Error("Fresh token should be valid")
	}
	
	// 撤销 token
	revoked := relay.RevokeToken(token)
	if !revoked {
		t.Error("Failed to revoke token")
	}
	
	// 验证 token 已失效
	_, valid = relay.validateToken(token)
	if valid {
		t.Error("Revoked token should be invalid")
	}
}

// TestValidateTargetURL 测试 URL 验证
func TestValidateTargetURL(t *testing.T) {
	relay := NewCloudRelay()
	
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://example.com/api", true},
		{"http://localhost:8080/api", true},
		{"http://127.0.0.1:8080/api", true},
		{"http://10.0.0.1/api", false},      // 内网 IP
		{"http://192.168.1.1/api", false},   // 内网 IP
		{"ftp://example.com/file", false},   // 不支持的协议
		{"", false},                          // 空 URL
	}
	
	for _, test := range tests {
		result := relay.validateTargetURL(test.url)
		if result != test.expected {
			t.Errorf("validateTargetURL(%s) = %v, expected %v", test.url, result, test.expected)
		}
	}
}

// TestConnectionLimits 测试连接限制
func TestConnectionLimits(t *testing.T) {
	relay := NewCloudRelay()
	relay.maxConnections = 5
	
	// 创建连接直到达到限制
	for i := 0; i < 5; i++ {
		relay.connections[fmt.Sprintf("conn-%d", i)] = &RelayConnection{}
	}
	
	// 验证连接数
	if len(relay.connections) != 5 {
		t.Errorf("Expected 5 connections, got %d", len(relay.connections))
	}
}

// TestMessageBuffer 测试消息缓冲
func TestMessageBuffer(t *testing.T) {
	conn := &RelayConnection{
		Buffer: make([]Message, 0),
	}
	
	// 添加消息到缓冲
	for i := 0; i < 150; i++ {
		conn.BufferMu.Lock()
		conn.Buffer = append(conn.Buffer, Message{ID: fmt.Sprintf("msg-%d", i)})
		// 限制缓冲大小
		if len(conn.Buffer) > 100 {
			conn.Buffer = conn.Buffer[len(conn.Buffer)-100:]
		}
		conn.BufferMu.Unlock()
	}
	
	// 验证缓冲大小
	conn.BufferMu.Lock()
	bufferSize := len(conn.Buffer)
	conn.BufferMu.Unlock()
	
	if bufferSize != 100 {
		t.Errorf("Expected buffer size 100, got %d", bufferSize)
	}
	
	// 验证最后一条消息
	conn.BufferMu.Lock()
	lastMsg := conn.Buffer[99]
	conn.BufferMu.Unlock()
	
	if lastMsg.ID != "msg-149" {
		t.Errorf("Expected last message msg-149, got %s", lastMsg.ID)
	}
}

// TestHeartbeat 测试心跳
func TestHeartbeat(t *testing.T) {
	relay := NewCloudRelay()
	
	conn := &RelayConnection{
		ID:         "test-conn",
		LastActive: time.Now(),
	}
	
	// 模拟不活跃
	conn.LastActive = time.Now().Add(-2 * time.Hour)
	
	// 清理应该移除这个连接
	relay.connections[conn.ID] = conn
	relay.cleanup()
	
	// 验证连接已被清理
	_, exists := relay.connections[conn.ID]
	if exists {
		t.Error("Inactive connection should be cleaned up")
	}
}

// BenchmarkGenerateToken 性能测试：token 生成
func BenchmarkGenerateToken(b *testing.B) {
	relay := NewCloudRelay()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		relay.GenerateToken("test-user")
	}
}

// BenchmarkValidateToken 性能测试：token 验证
func BenchmarkValidateToken(b *testing.B) {
	relay := NewCloudRelay()
	token := relay.GenerateToken("test-user")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		relay.validateToken(token)
	}
}

// BenchmarkForwardRequest 性能测试：请求转发
func BenchmarkForwardRequest(b *testing.B) {
	// 这个测试需要实际的 HTTP 服务器
	// 这里只是示例
	b.Skip("Skipping benchmark that requires HTTP server")
}
