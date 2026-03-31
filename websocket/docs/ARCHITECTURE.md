# 架构设计文档

## 1. 系统概述

MultiPass WebSocket 系统是一个基于 Go 的实时通信平台，提供 WebShell 终端访问和云端安全中转服务。系统设计支持高并发、低延迟的 WebSocket 连接，并通过 Cloudflare Tunnel 实现安全的公网暴露。

## 2. 架构层次

### 2.1 接入层

**组件**: Cloudflare Tunnel + Nginx

**职责**:
- HTTPS 终止
- DDoS 防护
- 请求路由
- WebSocket 协议升级
- 负载均衡

**配置要点**:
```nginx
# WebSocket 支持
map $http_upgrade $connection_upgrade {
    default upgrade;
    ''      close;
}

location /api/v1/ws/ {
    proxy_pass http://websocket_terminal;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $connection_upgrade;
    proxy_read_timeout 86400s;  # 长连接支持
}
```

### 2.2 应用层

#### 2.2.1 WebSocket Terminal Service (:8080)

**核心功能**:
1. WebSocket 连接管理
2. PTY 会话创建与维护
3. 命令执行与输出流式传输
4. 认证与授权
5. 心跳保活

**数据流**:
```
Client → WebSocket → Message Parser → PTY Manager → Shell
Client ← WebSocket ← Output Stream ← PTY Manager ← Shell
```

**关键结构**:
```go
type PTYTerminal struct {
    ID         string              // 会话 ID
    Ws         *websocket.Conn     // WebSocket 连接
    Pty        *os.File            // PTY 文件描述符
    Cmd        *exec.Cmd           // Shell 进程
    Mu         sync.Mutex          // 并发控制
    LastActive time.Time           // 最后活跃时间
    Timeout    time.Duration       // 超时时间
}
```

#### 2.2.2 Cloud Relay Service (:8081)

**核心功能**:
1. 用户端长连接管理
2. 请求转发与代理
3. 消息缓冲
4. 安全校验

**数据流**:
```
User → WebSocket → Message Buffer → Terminal WebSocket → Terminal
User ← WebSocket ← Message Buffer ← Terminal WebSocket ← Terminal
```

**转发流程**:
```
Client → POST /api/v1/forward → Auth Check → URL Validation 
       → HTTP Client → Target Server → Response → Client
```

### 2.3 存储层

#### 2.3.1 内存存储（默认）

**存储内容**:
- WebSocket 连接映射
- 认证 Token
- 会话状态

**优点**: 低延迟，简单
**缺点**: 重启丢失，单机限制

#### 2.3.2 Redis 存储（集群模式）

**存储内容**:
- 会话状态
- Token 黑名单
- 连接映射

**Key 设计**:
```
session:{session_id} → JSON(PTYTerminal)
token:{token_hash} → JSON(AuthToken)
user:{user_id}:connections → SET[connection_id]
```

### 2.4 基础设施层

**Cloudflare Tunnel**:
- 安全隧道
- 自动 HTTPS
- 全球加速
- DDoS 防护

**Docker**:
- 容器化部署
- 资源隔离
- 自动重启
- 健康检查

## 3. 核心流程

### 3.1 终端连接建立流程

```
1. Client 请求 WS /api/v1/ws/terminal?token=xxx
2. Server 验证 Token
   ├─ 无效 → 返回 401
   └─ 有效 → 继续
3. WebSocket 升级
4. 创建 PTY 会话
   ├─ 启动 bash 进程
   ├─ 分配 PTY
   └─ 记录会话
5. 发送 session_created 消息
6. 启动 Goroutine
   ├─ 读取 PTY 输出 → WebSocket
   ├─ 读取 WebSocket 输入 → PTY
   └─ 心跳监控
```

### 3.2 命令执行流程

```
1. Client 发送 input 消息
2. Server 解析消息
3. 写入 PTY (terminal.Pty.Write)
4. Shell 执行命令
5. PTY 输出捕获
6. 发送 output 消息到 Client
7. Client 渲染输出
```

### 3.3 请求转发流程

```
1. Client 发送 POST /api/v1/forward
2. Server 验证 X-Auth-Token
3. 解析请求体
4. 安全校验
   ├─ URL 白名单
   ├─ 协议检查 (HTTP/HTTPS)
   └─ 内网 IP 检查
5. 创建转发请求
6. 执行 HTTP 请求
7. 复制响应到 Client
8. 记录日志
```

### 3.4 Token 认证流程

```
1. 生成 Token
   ├─ UUID 生成
   ├─ Bcrypt 哈希
   └─ 存储 Token + 过期时间
2. 验证 Token
   ├─ 查找 Token
   ├─ 检查过期
   ├─ 检查撤销状态
   └─ 返回验证结果
3. 清理过期 Token
   └─ 定时任务每小时执行
```

## 4. 并发模型

### 4.1 Goroutine 使用

```
Main Thread
├─ HTTP Server (Goroutine Pool)
│  ├─ HandleTerminal (per connection)
│  ├─ HandleRegister (per connection)
│  └─ HandleForward (per request)
│
├─ Cleanup Goroutines
│  ├─ cleanupExpiredTokens (ticker)
│  └─ cleanupInactiveSessions (ticker)
│
└─ Session Goroutines
   ├─ handlePTYSession (per session)
   │  ├─ PTY Read → WS Write
   │  └─ WS Read → PTY Write
   └─ heartbeatMonitor (per session)
```

### 4.2 并发控制

**Mutex 使用**:
```go
// 保护共享数据
cm.mu.Lock()
cm.sessions[id] = session
cm.mu.Unlock()

// 只读访问
cm.mu.RLock()
session := cm.sessions[id]
cm.mu.RUnlock()
```

**Channel 使用**:
```go
// 命令队列
Commands: make(chan string, 100)

// 结果队列
Results: make(chan string, 100)
```

## 5. 错误处理

### 5.1 连接错误

```go
// WebSocket 读取错误
if err := conn.ReadJSON(&msg); err != nil {
    log.Printf("Read error: %v", err)
    conn.Close()
    return  // 退出 Goroutine
}
```

### 5.2 PTY 错误

```go
// PTY 启动失败
ptyFile, err := pty.Start(cmd)
if err != nil {
    log.Printf("PTY start error: %v", err)
    conn.Close()
    return
}
```

### 5.3 超时处理

```go
// 会话超时
select {
case <-activity:
    // 重置超时
case <-time.After(timeout):
    // 清理会话
    conn.Close()
    return
}
```

## 6. 安全设计

### 6.1 认证机制

- UUID Token 生成
- Bcrypt 密码哈希
- Token 过期时间
- Token 撤销机制

### 6.2 授权控制

- Token 验证中间件
- URL 白名单
- 内网访问限制
- SSRF 防护

### 6.3 数据保护

- HTTPS 加密传输
- Token 哈希存储
- 敏感信息不日志
- 连接隔离

### 6.4 防护措施

```go
// URL 验证
func validateTargetURL(rawURL string) bool {
    parsed, _ := url.Parse(rawURL)
    
    // 协议检查
    if parsed.Scheme != "http" && parsed.Scheme != "https" {
        return false
    }
    
    // 内网 IP 检查
    if isInternalIP(parsed.Hostname()) {
        return false
    }
    
    return true
}
```

## 7. 性能优化

### 7.1 内存管理

- Buffer Pool 复用
- 限制消息缓冲大小
- 及时关闭连接
- 定期清理过期数据

### 7.2 网络优化

- WebSocket 长连接
- Nginx Keepalive
- Cloudflare CDN
- 压缩传输

### 7.3 数据库优化（Redis）

- Pipeline 批量操作
- 合理设置 TTL
- 使用 Hash 结构
- 连接池复用

## 8. 监控与日志

### 8.1 日志级别

```go
log.Printf("[INFO] Session created: %s", id)
log.Printf("[WARN] Session timeout: %s", id)
log.Printf("[ERROR] WebSocket error: %v", err)
```

### 8.2 指标收集

- 连接数
- 消息数
- 延迟
- 错误率

### 8.3 健康检查

```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    stats := manager.GetStats()
    json.NewEncoder(w).Encode(stats)
})
```

## 9. 扩展性设计

### 9.1 水平扩展

- 无状态设计
- Redis 共享会话
- Nginx 负载均衡
- Docker 容器化

### 9.2 功能扩展

- 插件系统
- Webhook 支持
- 审计日志
- 录制回放

## 10. 容灾设计

### 10.1 故障恢复

- 自动重启
- 健康检查
- 连接重连
- 数据持久化

### 10.2 降级策略

- 限流保护
- 熔断机制
- 备用服务
- 手动开关
