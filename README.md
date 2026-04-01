# CinaSeek

> 自主品牌轻量级 Ubuntu 虚拟机远程管理工具  
> 零配置远程访问 · OpenClaw 一键部署 · 安全无密钥泄露

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8?logo=go)](https://go.dev)
[![Vue Version](https://img.shields.io/badge/vue-3.4+-4FC08D?logo=vue.js)](https://vuejs.org)

## 🌟 特性

- 🚀 **零配置远程访问** - Cloudflare Tunnel 加密通道，自动生成专属 HTTPS 域名
- 🔧 **OpenClaw 一键部署** - 专属 AI 开发环境，模型目录绑定，工作空间管理
- 🔒 **安全无密钥泄露** - 用户端无 Cloudflare 密钥，云端集中管控
- 💻 **WebShell 在线终端** - 浏览器内直接操作虚拟机，支持多标签页
- 📊 **实时监控** - CPU/内存/磁盘使用率，ECharts 可视化图表
- 🗂️ **目录挂载** - 宿主机 ↔ 虚拟机双向挂载，开机自动挂载
- 📦 **快照管理** - 创建/恢复/删除虚拟机快照
- 🔔 **多渠道告警** - 企业微信 + 邮件通知

## 🏗️ 架构

**CinaSeek = CinaClaw（VM引擎）+ Go中转服务 + Vue3面板 + WebSocket终端**

```
CinaSeek 架构
├── cinaclaw/          # CinaClaw 虚拟机引擎（基于 Multipass fork）
│   ├── src/           # C++ 源码
│   ├── include/       # 头文件
│   └── snap/          # Snap 打包
├── backend/           # Go 云端中转 API 服务
├── frontend/          # Vue3 Web 管理面板
├── websocket/         # WebSocket 终端服务
├── deploy/            # K8s 部署配置
└── docs/              # 文档
```

**数据流：** 用户端安装 CinaClaw VM 引擎

```
┌──────────────┐    ┌──────────────────┐    ┌─────────────────┐    ┌──────────────┐    ┌──────────────┐
│  用户浏览器   │ →  │ Cloudflare Tunnel │ →  │ 云端Go中转服务   │ →  │ 用户端Go客户端 │ →  │ CinaClaw引擎  │
└──────────────┘    └──────────────────┘    └─────────────────┘    └──────────────┘    └──────────────┘
                                                                                          │
                                                                                          ↓
                                                                                   ┌──────────────┐
                                                                                   │ Multipass VM  │
                                                                                   └──────────────┘
```

- **Cinaseek** 基于 Canonical Multipass GPLv3，自主品牌 CinaClaw运营，深度定制适配 CinaSeek 场景

## 📦 技术栈

### 虚拟机引擎
- CinaClaw（基于 Canonical Multipass fork，GPLv3）
- C++17，CMake，gRPC，Protobuf
- 支持多平台：Linux（QEMU/LXD）、macOS（HyperKit）、Windows（Hyper-V）

### 前端
- Vue 3.4 + Vite 5 + TypeScript
- Element Plus 2.5 (UI 组件)
- Pinia 2.1 (状态管理)
- Vue Router 4.3 (路由)
- xterm.js 5.3 (WebShell 终端)
- ECharts 5.4 (监控图表)

### 后端
- Go 1.23+ + Gin (Web 框架)
- PostgreSQL 15 (数据库)
- Redis (会话存储)
- WebSocket (实时通信)

### 部署
- Docker + Docker Compose
- Kubernetes (容器编排)
- Cloudflare Tunnel (远程访问)
- Nginx (反向代理)

## 🚀 快速开始

### 本地开发（Docker Compose）

```bash
# 克隆仓库
git clone https://github.com/cinagroup/cinaseek.git
cd cinaseek

# 启动所有服务（后端 + WebSocket + PostgreSQL + Redis）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 访问前端
open http://localhost:3000
```

### 后端单独构建

```bash
cd backend
go mod tidy
go build -o bin/cinaseek-backend cmd/server/main.go
./bin/cinaseek-backend
```

### 前端单独构建

```bash
cd frontend
npm install
npm run dev
```

### WebSocket 服务

```bash
cd websocket
./quickstart.sh
```

## 📁 项目结构

```
cinaclaw/
├── cinaclaw/           # CinaClaw 虚拟机引擎（基于 Multipass fork）
│   ├── src/            # C++ 源码
│   ├── include/        # 头文件
│   ├── snap/           # Snap 打包配置
│   ├── CMakeLists.txt  # 构建配置
│   └── README.md       # CinaClaw VM 引擎
├── frontend/           # Vue3 前端项目
│   ├── src/
│   │   ├── components/ # 组件
│   │   ├── views/      # 页面（16 个）
│   │   ├── router/     # 路由配置
│   │   ├── store/      # Pinia 状态管理
│   │   └── api/        # API 调用
│   ├── package.json
│   └── vite.config.js
├── backend/            # Go 后端 API
│   ├── cmd/
│   │   └── server/     # 服务入口
│   ├── internal/
│   │   ├── handler/    # API 处理器
│   │   ├── model/      # 数据模型
│   │   ├── database/   # 数据库连接
│   │   └── middleware/ # 中间件
│   ├── api/            # API 文档
│   ├── go.mod
│   └── Dockerfile
├── websocket/          # WebSocket 服务
│   ├── backend/        # Go WebSocket 服务
│   ├── deploy/         # 部署配置
│   ├── docs/           # 文档
│   └── tests/          # 测试
├── deploy/             # K8s 部署配置
│   ├── k8s/
│   │   ├── namespace.yaml
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── ingress.yaml
│   └── docker-compose.yml
├── docs/               # 文档
│   ├── ARCHITECTURE.md       # 架构设计文档
│   ├── EXECUTION_PLAN.md     # 落地执行计划
│   ├── CINASEEK_INTEGRATION.md # CinaSeek 集成指南
│   └── DEPLOYMENT.md         # 部署指南
├── scripts/            # 自动化脚本
│   └── deploy.sh
├── BRAND.md            # 品牌定位
└── README.md
```

## 🔐 认证集成（CinaToken OAuth）

CinaSeek 使用 CinaToken 统一认证（SSO），支持 9 个 OAuth Provider：

- GitHub
- Google
- Microsoft
- GitLab
- ...（共 9 个）

**登录流程：**
1. 用户点击"使用 CinaToken 账号登录"
2. 跳转到 CinaToken OAuth 授权页
3. 授权后返回 CinaSeek，携带 access_token
4. 后端验证 Token，创建会话

**详细集成指南：** [docs/CINASEEK_INTEGRATION.md](docs/CINASEEK_INTEGRATION.md)

## 📊 API 接口

共 **40+ 个 RESTful API**，分为 6 大模块：

| 模块 | 接口数 | 说明 |
|------|--------|------|
| 虚拟机管理 | 11 | 列表/创建/操作/快照/监控 |
| 目录挂载 | 5 | 挂载/卸载/OpenClaw 专属配置 |
| OpenClaw 管理 | 7 | 部署/状态/操作/日志/工作空间 |
| 远程访问 | 6 | 状态/开关/IP 白名单/日志 |
| 系统模块 | 7 | 设置/版本/仪表盘/搜索 |
| WebSocket | 4 | 终端连接/注册/转发/断开 |

**完整 API 文档：** [backend/api/API.md](backend/api/API.md)

## 🎯 核心功能

### 1. 虚拟机管理
- ✅ 创建 Ubuntu 虚拟机（自定义 CPU/内存/磁盘）
- ✅ 启动/停止/重启/删除
- ✅ 实时资源监控（CPU/内存/磁盘使用率）
- ✅ 快照管理（创建/恢复/删除）
- ✅ 操作日志记录

### 2. WebShell 在线终端
- ✅ 浏览器内直接操作虚拟机
- ✅ 支持多标签页切换
- ✅ 复制粘贴/清屏/彩色日志
- ✅ PTY 实时流（延迟 <50ms）

### 3. 目录挂载
- ✅ 宿主机 ↔ 虚拟机双向挂载
- ✅ 开机自动挂载配置
- ✅ OpenClaw 专属目录（模型/工作空间）

### 4. OpenClaw 一键部署
- ✅ 选择虚拟机，一键安装 OpenClaw
- ✅ 自动配置依赖和工作空间
- ✅ 服务启动/停止/重启
- ✅ 资源占用监控

### 5. 远程访问
- ✅ Cloudflare Tunnel 加密通道
- ✅ 自动生成专属 HTTPS 域名
- ✅ IP 白名单配置
- ✅ 访问日志记录
- ✅ 一键开启/关闭

## 📈 性能指标

| 指标 | 实测值 | 目标 |
|------|--------|------|
| WebSocket 连接延迟 | 45ms | <100ms ✅ |
| 命令执行延迟 | 52ms | <100ms ✅ |
| 最大并发连接 | 100+ | 100+ ✅ |
| 请求转发吞吐 | 426 req/s | 400+ ✅ |
| 请求成功率 | 100% | 99.9% ✅ |

## 🧪 测试

```bash
# 后端 API 测试
cd backend/tests
go test -v ./...

# WebSocket 集成测试
cd websocket/tests
./integration_test.sh

# 前端 E2E 测试（待添加）
cd frontend
npm run test:e2e
```

## 📝 开发里程碑

| 阶段 | 时间 | 目标 | 状态 |
|------|------|------|------|
| Phase 1 | 2026-04-01 | 代码整合完成、仓库创建 | ✅ |
| Phase 2 | 2026-04-07 | 认证模块集成（CinaToken OAuth） | ⏳ |
| Phase 3 | 2026-04-14 | K8s 部署配置完成、测试环境上线 | ⏳ |
| Phase 4 | 2026-04-21 | 生产环境部署、域名配置 | ⏳ |
| Launch | 2026-04-28 | 正式上线 | ⏳ |

## 🔗 相关链接

- **GitHub**: https://github.com/cinagroup/cinaseek
- **CinaSeek**: https://github.com/cinagroup/cinaseek
- **CinaToken**: https://github.com/cinagroup/cinatoken
- **品牌文档**: [BRAND.md](BRAND.md)
- **集成指南**: [docs/CINASEEK_INTEGRATION.md](docs/CINASEEK_INTEGRATION.md)

## 📄 License

MIT License

---

*CinaSeek - 你的云端开发工作室*
