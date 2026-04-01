# CinaSeek 架构设计

## 整体架构

```
[用户浏览器] → [Cloudflare Tunnel] → [云端Go中转服务] → [用户端Go客户端] → [CinaClaw引擎] → [Multipass虚拟机]
```

## 模块说明

### Cinaseek（cinaclaw/）
- **职责**：虚拟机生命周期管理（创建/启动/停止/快照/挂载）
- **技术**：C++17，CMake，gRPC，Protobuf
- **接口**：gRPC Unix Socket（本地通信）
- **品牌**：基于 Canonical Multipass fork，GPLv3 协议
- **定制**：
  - 客户端/守护进程名：`cinaseek` / `cinaseekd`
  - 环境变量前缀：`CINASEEK_*`
  - Protobuf 命名空间：`cinaseek`
  - 深度集成 CinaSeek 管理平台 + OpenClaw 专属优化

### Go 云端中转服务（backend/）
- **职责**：用户鉴权、WebSocket 连接管理、请求转发
- **技术**：Go 1.23+，Gin，PostgreSQL 15，Redis
- **安全**：Cloudflare Token 仅云端存储，用户端零密钥
- **API**：46 个 RESTful 接口

### WebSocket 终端服务（websocket/）
- **职责**：WebShell 远程终端、实时日志
- **技术**：Go，WebSocket，PTY，xterm.js
- **性能**：连接延迟 45ms，命令执行 52ms

### Web 管理面板（frontend/）
- **职责**：可视化操作界面
- **技术**：Vue3，TypeScript，Element Plus，xterm.js，ECharts
- **页面**：16 个功能页面

## 数据流

```
1. 用户访问 https://user.cinaseek.ai
2. Cloudflare Tunnel 转发至云端 Go 服务
3. Go 服务通过 WebSocket 查找用户端连接
4. 请求转发至用户端 Go 客户端
5. Go 客户端调用 Cinaseek gRPC 接口
6. Cinaseek 执行虚拟机操作
7. 结果原路返回
```

## 安全设计

| 层级 | 机制 |
|------|------|
| 传输加密 | 全程 TLS（Cloudflare → 云端 → WebSocket） |
| Cloudflare Token | 仅云端存储，用户端绝对不持有 |
| 用户鉴权 | CinaToken OAuth SSO，平台 Token 鉴权 |
| 数据隔离 | 用户间完全隔离，杜绝数据泄露 |
| 输入校验 | API 层参数验证，防止注入攻击 |
| 速率限制 | 令牌桶限流，防止恶意请求 |

## 品牌关系

```
CinaGroup 技术生态
├── CinaSeek           # LLM 聚合平台（核心产品）
│   └── CinaToken      # OAuth 认证（基础设施层）
└── CinaSeek           # 虚拟机远程管理平台
    └── Cinaseek       # VM 引擎（基于 Multipass fork，GPLv3）
```

## 商用分层

| 版本 | 实例数 | 并发 | 适用场景 |
|------|--------|------|----------|
| 免费版 | 1 | 低 | 个人开发者 |
| 专业版 | 多 | 中 | 小型团队 |
| 企业版 | 多 | 高 | 企业商用 |

---

*最后更新：2026-04-01*
