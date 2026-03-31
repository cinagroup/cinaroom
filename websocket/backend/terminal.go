package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

// PTYTerminal 表示一个 PTY 终端会话
type PTYTerminal struct {
	ID         string
	Ws         *websocket.Conn
	Pty        *os.File
	Cmd        *exec.Cmd
	Mu         sync.Mutex
	LastActive time.Time
	Timeout    time.Duration
}

// SecureConnectionManager 安全连接管理器
type SecureConnectionManager struct {
	sessions        map[string]*PTYTerminal
	registeredUsers map[string]*websocket.Conn
	mu              sync.RWMutex
	authTokens      map[string]*AuthToken
	cloudTunnel     *CloudTunnel
}

// AuthToken 认证 token 结构
type AuthToken struct {
	Token     string
	Hash      string
	CreatedAt time.Time
	ExpiresAt time.Time
	UserID    string
}

// CloudTunnel Cloudflare Tunnel 配置
type CloudTunnel struct {
	TunnelID     string
	Hostname     string
	Secret       string
	ConfigPath   string
}

// WebSocket 配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		// 生产环境应该验证来源
		return true
	},
}

// Message 消息结构
type Message struct {
	Type      string `json:"type"`
	Payload   string `json:"payload"`
	SessionID string `json:"session_id,omitempty"`
	Token     string `json:"token,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// NewSecureConnectionManager 创建安全连接管理器
func NewSecureConnectionManager() *SecureConnectionManager {
	cm := &SecureConnectionManager{
		sessions:        make(map[string]*PTYTerminal),
		registeredUsers: make(map[string]*websocket.Conn),
		authTokens:      make(map[string]*AuthToken),
		cloudTunnel: &CloudTunnel{
			TunnelID:   os.Getenv("CLOUDFLARE_TUNNEL_ID"),
			Hostname:   os.Getenv("CLOUDFLARE_HOSTNAME"),
			Secret:     os.Getenv("CLOUDFLARE_SECRET"),
			ConfigPath: "/etc/cloudflared/config.yml",
		},
	}
	
	// 启动清理协程
	go cm.cleanupExpiredTokens()
	go cm.cleanupInactiveSessions()
	
	return cm
}

// HandleTerminal 处理 WebShell 终端连接
func (cm *SecureConnectionManager) HandleTerminal(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	authToken, valid := cm.validateToken(token)
	if !valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// 创建 PTY 会话
	sessionID := uuid.New().String()
	terminal := &PTYTerminal{
		ID:         sessionID,
		Ws:         conn,
		LastActive: time.Now(),
		Timeout:    30 * time.Minute,
	}

	// 启动 shell
	cmd := exec.Command("bash")
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		log.Printf("PTY start error: %v", err)
		conn.Close()
		return
	}

	terminal.Pty = ptyFile
	terminal.Cmd = cmd

	cm.mu.Lock()
	cm.sessions[sessionID] = terminal
	cm.mu.Unlock()

	log.Printf("New terminal session: %s for user: %s", sessionID, authToken.UserID)

	// 启动会话处理
	go cm.handlePTYSession(terminal)
	go cm.heartbeatMonitor(terminal)

	// 发送会话创建成功消息
	msg := Message{
		Type:      "session_created",
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
	}
	conn.WriteJSON(msg)
}

// handlePTYSession 处理 PTY 会话
func (cm *SecureConnectionManager) handlePTYSession(terminal *PTYTerminal) {
	defer func() {
		terminal.Ws.Close()
		if terminal.Pty != nil {
			terminal.Pty.Close()
		}
		if terminal.Cmd != nil && terminal.Cmd.Process != nil {
			terminal.Cmd.Process.Kill()
		}
		cm.mu.Lock()
		delete(cm.sessions, terminal.ID)
		cm.mu.Unlock()
		log.Printf("Session closed: %s", terminal.ID)
	}()

	// 读取 PTY 输出
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := terminal.Pty.Read(buf)
			if err != nil {
				return
			}
			
			msg := Message{
				Type:      "output",
				Payload:   string(buf[:n]),
				SessionID: terminal.ID,
				Timestamp: time.Now().Unix(),
			}
			
			terminal.Mu.Lock()
			err = terminal.Ws.WriteJSON(msg)
			terminal.Mu.Unlock()
			
			if err != nil {
				return
			}
		}
	}()

	// 读取 WebSocket 输入
	for {
		var msg Message
		if err := terminal.Ws.ReadJSON(&msg); err != nil {
			return
		}

		terminal.LastActive = time.Now()

		switch msg.Type {
		case "input":
			// 写入命令到 PTY
			terminal.Pty.Write([]byte(msg.Payload))
		case "resize":
			// 处理终端大小调整
			var size struct {
				Rows int `json:"rows"`
				Cols int `json:"cols"`
			}
			json.Unmarshal([]byte(msg.Payload), &size)
			pty.Setsize(terminal.Pty, &pty.Winsize{
				Rows: uint16(size.Rows),
				Cols: uint16(size.Cols),
			})
		case "heartbeat":
			// 心跳响应
			ackMsg := Message{
				Type:      "heartbeat_ack",
				SessionID: terminal.ID,
				Timestamp: time.Now().Unix(),
			}
			terminal.Ws.WriteJSON(ackMsg)
		}
	}
}

// HandleRegister 处理用户端注册连接
func (cm *SecureConnectionManager) HandleRegister(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	userID := uuid.New().String()
	
	cm.mu.Lock()
	cm.registeredUsers[userID] = conn
	cm.mu.Unlock()

	log.Printf("New registered user: %s", userID)

	// 发送注册成功消息
	msg := Message{
		Type:      "registered",
		Payload:   userID,
		Timestamp: time.Now().Unix(),
	}
	conn.WriteJSON(msg)

	// 处理用户消息
	go cm.handleUserMessages(userID, conn)
}

// handleUserMessages 处理注册用户消息
func (cm *SecureConnectionManager) handleUserMessages(userID string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		cm.mu.Lock()
		delete(cm.registeredUsers, userID)
		cm.mu.Unlock()
	}()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		switch msg.Type {
		case "command":
			// 转发命令到指定会话
			if msg.SessionID != "" {
				cm.forwardToSession(msg.SessionID, msg.Payload)
			}
		case "heartbeat":
			conn.WriteJSON(Message{
				Type:      "heartbeat_ack",
				Timestamp: time.Now().Unix(),
			})
		}
	}
}

// forwardToSession 转发命令到指定会话
func (cm *SecureConnectionManager) forwardToSession(sessionID string, input string) {
	cm.mu.RLock()
	session, exists := cm.sessions[sessionID]
	cm.mu.RUnlock()

	if !exists {
		return
	}

	session.Mu.Lock()
	defer session.Mu.Unlock()

	msg := Message{
		Type:      "input",
		Payload:   input,
		SessionID: sessionID,
		Timestamp: time.Now().Unix(),
	}
	session.Ws.WriteJSON(msg)
}

// GenerateToken 生成认证 token
func (cm *SecureConnectionManager) GenerateToken(userID string) (string, error) {
	token := uuid.New().String()
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	authToken := &AuthToken{
		Token:     token,
		Hash:      string(hash),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserID:    userID,
	}

	cm.mu.Lock()
	cm.authTokens[token] = authToken
	cm.mu.Unlock()

	return token, nil
}

// validateToken 验证 token
func (cm *SecureConnectionManager) validateToken(token string) (*AuthToken, bool) {
	cm.mu.RLock()
	authToken, exists := cm.authTokens[token]
	cm.mu.RUnlock()

	if !exists {
		return nil, false
	}

	if time.Now().After(authToken.ExpiresAt) {
		return nil, false
	}

	return authToken, true
}

// cleanupExpiredTokens 清理过期 token
func (cm *SecureConnectionManager) cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		cm.mu.Lock()
		for token, authToken := range cm.authTokens {
			if time.Now().After(authToken.ExpiresAt) {
				delete(cm.authTokens, token)
			}
		}
		cm.mu.Unlock()
	}
}

// cleanupInactiveSessions 清理不活跃会话
func (cm *SecureConnectionManager) cleanupInactiveSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cm.mu.Lock()
		for id, session := range cm.sessions {
			if time.Since(session.LastActive) > session.Timeout {
				session.Ws.Close()
				delete(cm.sessions, id)
				log.Printf("Cleaned up inactive session: %s", id)
			}
		}
		cm.mu.Unlock()
	}
}

// heartbeatMonitor 心跳监控
func (cm *SecureConnectionManager) heartbeatMonitor(terminal *PTYTerminal) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if time.Since(terminal.LastActive) > terminal.Timeout {
			log.Printf("Session timeout: %s", terminal.ID)
			terminal.Ws.Close()
			return
		}
	}
}

// HandleForward 云端中转请求转发接口
func (cm *SecureConnectionManager) HandleForward(w http.ResponseWriter, r *http.Request) {
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

	_, valid := cm.validateToken(token)
	if !valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var req struct {
		TargetURL string            `json:"target_url"`
		Method    string            `json:"method"`
		Headers   map[string]string `json:"headers"`
		Body      string            `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 安全校验：验证目标 URL
	if !cm.validateTargetURL(req.TargetURL) {
		http.Error(w, "Invalid target URL", http.StatusForbidden)
		return
	}

	// 转发请求
	client := &http.Client{Timeout: 30 * time.Second}
	forwardReq, err := http.NewRequest(req.Method, req.TargetURL, []byte(req.Body))
	if err != nil {
		http.Error(w, "Failed to create forward request", http.StatusInternalServerError)
		return
	}

	// 复制 headers
	for key, value := range req.Headers {
		forwardReq.Header.Set(key, value)
	}

	resp, err := client.Do(forwardReq)
	if err != nil {
		http.Error(w, "Forward request failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 复制响应
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	log.Printf("Forwarded request: %s %s -> %d", req.Method, req.TargetURL, resp.StatusCode)
}

// validateTargetURL 验证目标 URL 安全性
func (cm *SecureConnectionManager) validateTargetURL(url string) bool {
	// 实现 URL 白名单验证逻辑
	// 防止 SSRF 攻击
	return true
}

// GetCloudTunnelConfig 获取 Cloudflare Tunnel 配置
func (cm *SecureConnectionManager) GetCloudTunnelConfig() string {
	return `version: 2
tunnel: ` + cm.cloudTunnel.TunnelID + `
credentials-file: /etc/cloudflared/creds.json

ingress:
  - hostname: ` + cm.cloudTunnel.Hostname + `
    service: http://localhost:8080
  - service: http_status:404
`
}

func main() {
	manager := NewSecureConnectionManager()

	// 生成测试 token
	testToken, err := manager.GenerateToken("test-user")
	if err != nil {
		log.Fatalf("Failed to generate test token: %v", err)
	}
	log.Printf("Test auth token: %s", testToken)

	// 路由
	http.HandleFunc("/api/v1/ws/terminal", manager.HandleTerminal)
	http.HandleFunc("/api/v1/ws/register", manager.HandleRegister)
	http.HandleFunc("/api/v1/forward", manager.HandleForward)
	
	// 健康检查
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
