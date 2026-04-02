package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TerminalSession 表示一个 WebShell 终端会话
type TerminalSession struct {
	ID        string
	Ws        *websocket.Conn
	Pty       interface{} // pty.PtyMaster
	Commands  chan string
	Results   chan string
	CreatedAt time.Time
	LastActive time.Time
	AuthToken string
}

// ConnectionManager 管理所有 WebSocket 连接
type ConnectionManager struct {
	sessions     map[string]*TerminalSession
	registeredUsers map[string]*websocket.Conn
	mu           sync.RWMutex
	authTokens   map[string]bool
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		sessions:        make(map[string]*TerminalSession),
		registeredUsers: make(map[string]*websocket.Conn),
		authTokens:      make(map[string]bool),
	}
}

// WebSocket 配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应该限制来源
	},
}

// TerminalMessage 终端消息结构
type TerminalMessage struct {
	Type    string `json:"type"`    // "command", "result", "heartbeat", "auth"
	Payload string `json:"payload"`
	Token   string `json:"token,omitempty"`
}

// HandleTerminal 处理 WebShell 终端连接
func (cm *ConnectionManager) HandleTerminal(w http.ResponseWriter, r *http.Request) {
	// 验证 token
	token := r.URL.Query().Get("token")
	if token == "" || !cm.validateToken(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 升级 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// 创建新会话
	sessionID := uuid.New().String()
	session := &TerminalSession{
		ID:         sessionID,
		Ws:         conn,
		Commands:   make(chan string, 100),
		Results:    make(chan string, 100),
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
		AuthToken:  token,
	}

	cm.mu.Lock()
	cm.sessions[sessionID] = session
	cm.mu.Unlock()

	log.Printf("New terminal session: %s", sessionID)

	// 启动会话处理
	go cm.handleSession(session)

	// 启动心跳检测
	go cm.heartbeatCheck(session)
}

// HandleRegister 处理用户端注册连接
func (cm *ConnectionManager) HandleRegister(w http.ResponseWriter, r *http.Request) {
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
	response := TerminalMessage{
		Type:    "registered",
		Payload: userID,
	}
	conn.WriteJSON(response)

	// 处理用户消息
	go cm.handleUserMessages(userID, conn)
}

// handleSession 处理终端会话的命令和结果
func (cm *ConnectionManager) handleSession(session *TerminalSession) {
	defer func() {
		session.Ws.Close()
		cm.mu.Lock()
		delete(cm.sessions, session.ID)
		cm.mu.Unlock()
		log.Printf("Session closed: %s", session.ID)
	}()

	for {
		// 读取命令
		select {
		case cmd := <-session.Commands:
			// 执行命令并返回结果
			result := executeCommand(cmd)
			session.Results <- result
		case result := <-session.Results:
			msg := TerminalMessage{
				Type:    "result",
				Payload: result,
			}
			if err := session.Ws.WriteJSON(msg); err != nil {
				log.Printf("Write error: %v", err)
				return
			}
		case <-time.After(30 * time.Minute): // 30 分钟超时
			log.Printf("Session timeout: %s", session.ID)
			return
		}
	}
}

// handleUserMessages 处理注册用户消息
func (cm *ConnectionManager) handleUserMessages(userID string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		cm.mu.Lock()
		delete(cm.registeredUsers, userID)
		cm.mu.Unlock()
	}()

	for {
		var msg TerminalMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		switch msg.Type {
		case "command":
			// 转发命令到终端
			cm.forwardCommand(msg.Payload)
		case "heartbeat":
			// 心跳响应
			conn.WriteJSON(TerminalMessage{Type: "heartbeat_ack"})
		}
	}
}

// forwardCommand 转发命令到所有终端
func (cm *ConnectionManager) forwardCommand(cmd string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	for _, session := range cm.sessions {
		session.Commands <- cmd
	}
}

// executeCommand 执行命令（简化版本）
func executeCommand(cmd string) string {
	// 实际实现应该使用 pty 执行命令
	return fmt.Sprintf("Executed: %s\n", cmd)
}

// validateToken 验证 token
func (cm *ConnectionManager) validateToken(token string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.authTokens[token]
}

// GenerateToken 生成新的认证 token
func (cm *ConnectionManager) GenerateToken() string {
	token := uuid.New().String()
	cm.mu.Lock()
	cm.authTokens[token] = true
	cm.mu.Unlock()
	return token
}

// heartbeatCheck 心跳检测
func (cm *ConnectionManager) heartbeatCheck(session *TerminalSession) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if time.Since(session.LastActive) > 5*time.Minute {
			log.Printf("Session inactive: %s", session.ID)
			session.Ws.Close()
			return
		}
	}
}

// HandleForward 云端中转请求转发接口
func (cm *ConnectionManager) HandleForward(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 验证请求
	token := r.Header.Get("X-Auth-Token")
	if token == "" || !cm.validateToken(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		TargetURL string `json:"target_url"`
		Method    string `json:"method"`
		Body      string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 转发请求（简化版本）
	log.Printf("Forwarding request to: %s", req.TargetURL)
	
	response := map[string]interface{}{
		"status": "forwarded",
		"url":    req.TargetURL,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	manager := NewConnectionManager()

	// 生成测试 token
	testToken := manager.GenerateToken()
	log.Printf("Test auth token: %s", testToken)

	http.HandleFunc("/api/v1/ws/terminal", manager.HandleTerminal)
	http.HandleFunc("/api/v1/ws/register", manager.HandleRegister)
	http.HandleFunc("/api/v1/forward", manager.HandleForward)
	
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"service":   "cinaseek-websocket",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	addr := os.Getenv("WS_PORT")
	if addr == "" {
		addr = "8081"
	}
	if addr[0] != ':' {
		addr = ":" + addr
	}

	log.Printf("WebSocket server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
