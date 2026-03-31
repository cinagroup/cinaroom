package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// CloudRelay 云端中转服务
type CloudRelay struct {
	mu            sync.RWMutex
	connections   map[string]*RelayConnection
	authTokens    map[string]*RelayAuthToken
	allowedHosts  map[string]bool
	maxConnections int
}

// RelayConnection 中继连接
type RelayConnection struct {
	ID          string
	UserConn    *websocket.Conn
	TerminalConn *websocket.Conn
	CreatedAt   time.Time
	LastActive  time.Time
	Buffer      []Message
	BufferMu    sync.Mutex
}

// RelayAuthToken 中继认证 token
type RelayAuthToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	Revoked   bool
}

// Message 中继消息
type Message struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Payload   string `json:"payload"`
	From      string `json:"from"`
	To        string `json:"to"`
	Timestamp int64  `json:"timestamp"`
}

// NewCloudRelay 创建云端中转服务
func NewCloudRelay() *CloudRelay {
	relay := &CloudRelay{
		connections:    make(map[string]*RelayConnection),
		authTokens:     make(map[string]*RelayAuthToken),
		allowedHosts:   make(map[string]bool),
		maxConnections: 1000,
	}

	// 添加允许的 hosts
	relay.allowedHosts["localhost"] = true
	relay.allowedHosts["127.0.0.1"] = true

	// 启动清理协程
	go relay.cleanup()

	return relay
}

// HandleUserConnect 处理用户端连接
func (cr *CloudRelay) HandleUserConnect(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	authToken, valid := cr.validateToken(token)
	if !valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	connectionID := uuid.New().String()
	relayConn := &RelayConnection{
		ID:         connectionID,
		UserConn:   conn,
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
		Buffer:     make([]Message, 0),
	}

	cr.mu.Lock()
	if len(cr.connections) >= cr.maxConnections {
		cr.mu.Unlock()
		conn.Close()
		http.Error(w, "Too many connections", http.StatusServiceUnavailable)
		return
	}
	cr.connections[connectionID] = relayConn
	cr.mu.Unlock()

	log.Printf("User connected: %s (user: %s)", connectionID, authToken.UserID)

	// 发送连接成功消息
	msg := Message{
		ID:        uuid.New().String(),
		Type:      "connected",
		Payload:   connectionID,
		Timestamp: time.Now().Unix(),
	}
	conn.WriteJSON(msg)

	// 处理用户消息
	go cr.handleUserMessages(relayConn)
	go cr.heartbeatCheck(relayConn)
}

// HandleTerminalConnect 处理终端连接
func (cr *CloudRelay) HandleTerminalConnect(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	connectionID := r.URL.Query().Get("connection_id")

	if token == "" || connectionID == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	_, valid := cr.validateToken(token)
	if !valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	cr.mu.Lock()
	relayConn, exists := cr.connections[connectionID]
	if !exists {
		cr.mu.Unlock()
		conn.Close()
		http.Error(w, "Connection not found", http.StatusNotFound)
		return
	}
	relayConn.TerminalConn = conn
	cr.mu.Unlock()

	log.Printf("Terminal connected to session: %s", connectionID)

	// 发送缓冲消息
	relayConn.BufferMu.Lock()
	for _, msg := range relayConn.Buffer {
		conn.WriteJSON(msg)
	}
	relayConn.Buffer = nil // 清空缓冲
	relayConn.BufferMu.Unlock()

	// 处理终端消息
	go cr.handleTerminalMessages(relayConn)
}

// handleUserMessages 处理用户消息
func (cr *CloudRelay) handleUserMessages(relayConn *RelayConnection) {
	defer func() {
		relayConn.UserConn.Close()
		cr.removeConnection(relayConn.ID)
	}()

	for {
		var msg Message
		if err := relayConn.UserConn.ReadJSON(&msg); err != nil {
			log.Printf("User read error: %v", err)
			return
		}

		relayConn.LastActive = time.Now()
		msg.From = "user"
		msg.ID = uuid.New().String()
		msg.Timestamp = time.Now().Unix()

		// 转发到终端
		if relayConn.TerminalConn != nil {
			relayConn.TerminalConn.WriteJSON(msg)
		} else {
			// 缓冲消息
			relayConn.BufferMu.Lock()
			relayConn.Buffer = append(relayConn.Buffer, msg)
			// 限制缓冲大小
			if len(relayConn.Buffer) > 100 {
				relayConn.Buffer = relayConn.Buffer[len(relayConn.Buffer)-100:]
			}
			relayConn.BufferMu.Unlock()
		}
	}
}

// handleTerminalMessages 处理终端消息
func (cr *CloudRelay) handleTerminalMessages(relayConn *RelayConnection) {
	defer func() {
		relayConn.TerminalConn.Close()
		relayConn.TerminalConn = nil
	}()

	for {
		var msg Message
		if err := relayConn.TerminalConn.ReadJSON(&msg); err != nil {
			log.Printf("Terminal read error: %v", err)
			return
		}

		relayConn.LastActive = time.Now()
		msg.From = "terminal"
		msg.ID = uuid.New().String()
		msg.Timestamp = time.Now().Unix()

		// 转发到用户
		if relayConn.UserConn != nil {
			relayConn.UserConn.WriteJSON(msg)
		}
	}
}

// HandleForward 处理转发请求
func (cr *CloudRelay) HandleForward(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 验证 token
	token := r.Header.Get("X-Auth-Token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	authToken, valid := cr.validateToken(token)
	if !valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// 解析请求
	var req struct {
		ConnectionID string            `json:"connection_id"`
		TargetURL    string            `json:"target_url"`
		Method       string            `json:"method"`
		Headers      map[string]string `json:"headers"`
		Body         string            `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 验证连接存在
	cr.mu.RLock()
	_, exists := cr.connections[req.ConnectionID]
	cr.mu.RUnlock()

	if !exists {
		http.Error(w, "Connection not found", http.StatusNotFound)
		return
	}

	// 安全校验目标 URL
	if !cr.validateTargetURL(req.TargetURL) {
		http.Error(w, "Invalid target URL", http.StatusForbidden)
		return
	}

	// 创建转发请求
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	forwardReq, err := http.NewRequestWithContext(ctx, req.Method, req.TargetURL, []byte(req.Body))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// 添加 headers
	for key, value := range req.Headers {
		forwardReq.Header.Set(key, value)
	}
	forwardReq.Header.Set("X-Forwarded-By", "cloud-relay")
	forwardReq.Header.Set("X-User-ID", authToken.UserID)

	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(forwardReq)
	if err != nil {
		http.Error(w, "Forward failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 复制响应
	for key, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	log.Printf("Forwarded: %s %s -> %d", req.Method, req.TargetURL, resp.StatusCode)
}

// validateTargetURL 验证目标 URL
func (cr *CloudRelay) validateTargetURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// 检查协议
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}

	// 检查 host
	host := parsed.Hostname()
	if cr.allowedHosts[host] {
		return true
	}

	// 防止内网访问
	if isInternalIP(host) {
		return false
	}

	return true
}

// isInternalIP 检查是否为内网 IP
func isInternalIP(host string) bool {
	// 简化实现，生产环境需要完整的 IP 检查
	internalPrefixes := []string{
		"10.",
		"172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.",
		"172.24.", "172.25.", "172.26.", "172.27.",
		"172.28.", "172.29.", "172.30.", "172.31.",
		"192.168.",
		"127.",
	}

	for _, prefix := range internalPrefixes {
		if len(host) > len(prefix) && host[:len(prefix)] == prefix {
			return true
		}
	}

	return false
}

// validateToken 验证 token
func (cr *CloudRelay) validateToken(token string) (*RelayAuthToken, bool) {
	cr.mu.RLock()
	authToken, exists := cr.authTokens[token]
	cr.mu.RUnlock()

	if !exists || authToken.Revoked {
		return nil, false
	}

	if time.Now().After(authToken.ExpiresAt) {
		return nil, false
	}

	return authToken, true
}

// GenerateToken 生成 token
func (cr *CloudRelay) GenerateToken(userID string) string {
	token := uuid.New().String()
	authToken := &RelayAuthToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	cr.mu.Lock()
	cr.authTokens[token] = authToken
	cr.mu.Unlock()

	return token
}

// RevokeToken 撤销 token
func (cr *CloudRelay) RevokeToken(token string) bool {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if authToken, exists := cr.authTokens[token]; exists {
		authToken.Revoked = true
		return true
	}
	return false
}

// removeConnection 移除连接
func (cr *CloudRelay) removeConnection(id string) {
	cr.mu.Lock()
	delete(cr.connections, id)
	cr.mu.Unlock()
	log.Printf("Connection removed: %s", id)
}

// cleanup 清理不活跃连接
func (cr *CloudRelay) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cr.mu.Lock()
		for id, conn := range cr.connections {
			if time.Since(conn.LastActive) > 1*time.Hour {
				if conn.UserConn != nil {
					conn.UserConn.Close()
				}
				if conn.TerminalConn != nil {
					conn.TerminalConn.Close()
				}
				delete(cr.connections, id)
				log.Printf("Cleaned up inactive connection: %s", id)
			}
		}
		cr.mu.Unlock()
	}
}

// heartbeatCheck 心跳检查
func (cr *CloudRelay) heartbeatCheck(relayConn *RelayConnection) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if time.Since(relayConn.LastActive) > 30*time.Minute {
			log.Printf("Connection timeout: %s", relayConn.ID)
			relayConn.UserConn.Close()
			return
		}

		// 发送心跳
		msg := Message{
			ID:        uuid.New().String(),
			Type:      "heartbeat",
			Timestamp: time.Now().Unix(),
		}
		relayConn.UserConn.WriteJSON(msg)
	}
}

// GetStats 获取服务统计
func (cr *CloudRelay) GetStats() map[string]interface{} {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	return map[string]interface{}{
		"total_connections": len(cr.connections),
		"active_tokens":     len(cr.authTokens),
		"max_connections":   cr.maxConnections,
		"uptime":            time.Since(time.Now()).String(),
	}
}

func main() {
	relay := NewCloudRelay()

	// 路由
	http.HandleFunc("/api/v1/ws/user", relay.HandleUserConnect)
	http.HandleFunc("/api/v1/ws/terminal", relay.HandleTerminalConnect)
	http.HandleFunc("/api/v1/forward", relay.HandleForward)

	// 健康检查
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(relay.GetStats())
	})

	log.Printf("Cloud Relay starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
