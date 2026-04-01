# CinaSeek - 自主品牌远程管理平台

## 品牌定位

**CinaSeek** 是自主品牌轻量级 Ubuntu 虚拟机远程管理工具，主打：
- 🚀 **零配置远程访问** - Cloudflare Tunnel 加密通道
- 🔧 **OpenClaw 一键部署** - 专属 AI 开发环境
- 🔒 **安全无密钥泄露** - 用户端无 Cloudflare 密钥，云端集中管控
- 💻 **全平台支持** - Windows/macOS/Linux 本地/远程无缝切换

## 品牌关系

```
CinaGroup 技术生态
├── CinaToken         # LLM 聚合平台（核心产品）
│   └── CinaToken    # OAuth 认证（基础设施层）
└── CinaSeek          # 虚拟机远程管理平台
    └── CinaClaw VM 引擎
```

## 核心引擎：Cinaseek

**Cinaseek** 是 CinaSeek 的虚拟机管理引擎，基于 Canonical Multipass (GPLv3) fork 而来：
- 🏭 **自主品牌运营** - 独立品牌标识，非 Multipass 附属
- 🔧 **深度定制** - 针对 CinaSeek 远程管理场景优化
- 🔌 **gRPC 接口** - 通过 Unix Socket 提供 VM 生命周期管理
- 📦 **Snap 分发** - 支持 Linux/macOS/Windows 多平台

## 技术架构

### 复用 CinaSeek 基础设施

| 层级 | 复用内容 | 说明 |
|------|----------|------|
| **认证** | CinaToken OAuth | 9 个 OAuth Provider、39+ LLM 渠道 |
| **数据库** | PostgreSQL 主从 | 服务器 A(43.156.66.122) + 服务器 B(101.32.108.223) |
| **部署** | Kubernetes 集群 | Master 节点：101.32.108.223 |
| **监控** | Uptime Kuma | 共享监控面板、企业微信 + 邮件告警 |
| **SSL** | Cloudflare Full (Strict) | 15 年 Origin CA 证书 |

### 独立运营部分

| 项目 | 配置 |
|------|------|
| **域名** | cinaseek.run（待注册） |
| **API 域名** | api.cinaseek.run（待注册） |
| **GitHub 仓库** | cinagroup/cinaseek |
| **Docker 镜像** | cinagroup/cinaseek:* |
| **数据库 Schema** | cinaseek（独立 schema） |

## 核心功能

1. **虚拟机管理** - 创建/启动/停止/快照/监控
2. **WebShell 终端** - 浏览器内操作虚拟机（xterm.js）
3. **目录挂载** - 宿主机 ↔ 虚拟机双向挂载
4. **OpenClaw 专属部署** - 一键安装、管控、资源监控
5. **远程访问** - Cloudflare Tunnel 加密通道
6. **实时监控** - CPU/内存/磁盘使用率（ECharts）

## 技术栈

- **前端**: Vue 3.4 + Vite 5 + Element Plus 2.5 + Pinia + Vue Router
- **后端**: Go 1.23+ + Gin + PostgreSQL 15 + Redis
- **WebSocket**: xterm.js 5.3 + PTY 实时流
- **部署**: Docker + Kubernetes + Cloudflare Tunnel + Nginx

## 开发里程碑

| 阶段 | 时间 | 目标 |
|------|------|------|
| Phase 1 | 2026-04-01 | 代码整合完成、仓库创建 |
| Phase 2 | 2026-04-07 | 认证模块集成（CinaToken OAuth） |
| Phase 3 | 2026-04-14 | K8s 部署配置完成、测试环境上线 |
| Phase 4 | 2026-04-21 | 生产环境部署、域名配置 |
| Launch | 2026-04-28 | 正式上线 |

## 品牌资产

- **名称**: CinaSeek
- **口号**: "你的云端开发工作室"
- **定位**: 独立品牌，非 CinaSeek 子模块
- **受众**: 个人开发者、小型团队、AI 研究者

---

*最后更新：2026-04-01*
