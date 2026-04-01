# CinaSeek 落地执行计划

> 基于《Multipass自主品牌Web远程管理平台 完整落地执行计划》整理  
> 最后更新：2026-04-01

---

## 一、项目核心定位

基于开源 Multipass 深度定制，整合 Go 语言安全中转服务 + Cloudflare Tunnel 商用内网穿透，打造自主品牌 CinaSeek 轻量 Ubuntu 虚拟机 Web 管理平台。

- **零配置远程访问**：无公网 IP、无端口映射、开箱即用
- **安全无密钥泄露**：Cloudflare Token 仅存云端，用户端零密钥
- **OpenClaw 专属适配**：一键部署 AI 开发环境
- **商用合规**：遵循 Multipass GPLv3 + Cloudflare 商用条款
- **面向**：个人开发者、小型团队、AI 研究者

## 二、项目核心原则

| 原则 | 说明 |
|------|------|
| 安全第一 | Cloudflare Token 仅存云端，用户端无任何敏感密钥，杜绝泄露风险 |
| 极简部署 | 用户一键安装，无需配置网络、穿透、端口，开箱即用 |
| 商用合规 | 遵循 Multipass GPLv3 协议，Cloudflare 商用条款，无版权风险 |
| 可规模化 | 免费版/付费版分层，支持海量用户并发 |

## 三、品牌与域名

| 项目 | 值 |
|------|-----|
| 品牌名 | CinaSeek |
| 域名 | **cinaseek.ai**（已注册）|
| API 域名 | **api.cinaseek.ai**（已注册）|
| GitHub | github.com/cinagroup/cinaseek |
| 品牌关系 | CinaGroup → CinaToken（LLM聚合）→ CinaSeek（管理平台）→ CinaClaw（VM引擎）|

---

## 四、前期准备（1-3天）

### 4.1 资质与账号

- [x] 注册域名 `cinaseek.ai`
- [x] 注册域名 `api.cinaseek.ai`
- [ ] Cloudflare 托管 `cinaseek.ai` 域名（NameServer 迁移）
- [ ] Cloudflare 申请 Origin CA 证书（15年）
- [ ] 准备品牌素材：LOGO（爪印图标）、favicon、OG 图片
- [ ] 梳理 Multipass GPLv3 协议合规要点

### 4.2 技术环境搭建

- [ ] 云服务器 A（43.156.66.122）：PostgreSQL 主库 + Go 中转服务
- [ ] 云服务器 B（101.32.108.223）：PostgreSQL 从库 + K8s Master
- [ ] 本地搭建 CinaClaw 测试环境，验证 gRPC 接口
- [ ] 云端安装 Cloudflared，配置基础隧道
- [ ] 编写 Go 项目脚手架，划分模块边界

### 4.3 基础设施复用

| 层级 | 复用来源 | 配置 |
|------|----------|------|
| 认证 | CinaToken OAuth | 9 Provider，SSO 统一登录 |
| 数据库 | CinaSeek PostgreSQL 主从 | Schema: `cinaseek`（独立）|
| 部署 | K8s 集群 | Namespace: `cinaseek` |
| 监控 | Uptime Kuma | 共享监控面板 |
| SSL | Cloudflare Full (Strict) | 15年 Origin CA |

---

## 五、分阶段开发计划（14-21天）

### 阶段一：核心功能开发（7天）

#### 1.1 CinaClaw 引擎编译验证（Day 1-2）

```
目标：确认 fork 的 CinaClaw 引擎可编译运行
```

- [ ] Linux 编译 CinaClaw（CMake + vcpkg + gRPC）
- [ ] macOS 编译 CinaClaw（Homebrew + Xcode）
- [ ] Windows 编译 CinaClaw（vcpkg + MSVC）
- [ ] 验证 `cinaclawd` daemon 启动正常
- [ ] 验证 `cinaclaw` CLI 基本命令（launch/list/info/shell）
- [ ] 确认 gRPC Unix Socket 接口可用
- [ ] 确认环境变量 `CINACLAWS_*` 前缀生效
- [ ] 编写编译文档 `cinaclaw/BUILD.linux.md` 更新

**交付物**：三平台编译通过的 CinaClaw 二进制 + 验证报告

#### 1.2 Go 客户端对接 CinaClaw gRPC（Day 2-4）

```
目标：Go 封装 CinaClaw gRPC，暴露管理接口
```

**模块：`backend/internal/cinaclaw/`**

- [ ] 从 `cinaclaw.proto` 生成 Go gRPC 客户端代码
- [ ] 实现 gRPC 连接管理（Unix Socket + 认证）
- [ ] 封装虚拟机管理接口：
  - [ ] `CreateVM(name, cpu, mem, disk, image)` → gRPC Launch
  - [ ] `StartVM(name)` / `StopVM(name)` / `RestartVM(name)`
  - [ ] `DeleteVM(name)` / `PurgeVM()`
  - [ ] `ListVMs()` → gRPC List
  - [ ] `GetVMInfo(name)` → gRPC Info
  - [ ] `SnapshotVM(name, snapshotName)` / `RestoreVM(name, snapshot)`
- [ ] 封装目录挂载接口：
  - [ ] `MountPath(hostPath, vmPath, autoMount)`
  - [ ] `UnmountPath(vmPath)`
- [ ] 封装资源监控接口：
  - [ ] `GetMetrics(name)` → CPU/内存/磁盘实时数据
- [ ] 单元测试覆盖率 > 80%

**交付物**：`backend/internal/cinaclaw/` 包 + 测试

#### 1.3 云端中转服务开发（Day 3-5）

```
目标：Go 云端服务，处理用户鉴权 + WebSocket 转发
```

- [ ] **用户鉴权模块**
  - [ ] 集成 CinaToken OAuth（`/auth/cinatoken/callback`）
  - [ ] JWT Token 生成/验证/刷新
  - [ ] 用户信息存储（PostgreSQL `users` 表）
  - [ ] Session 管理（Redis）
- [ ] **WebSocket 连接管理**
  - [ ] 用户端注册连接（`ws://api.cinaseek.ai/ws/register`）
  - [ ] 心跳保活（30s interval）
  - [ ] 断线自动重连（指数退避）
  - [ ] 用户连接池（并发安全 map）
- [ ] **请求转发模块**
  - [ ] 接收 Cloudflare 转发的公网请求
  - [ ] 路由到对应用户端 WebSocket 连接
  - [ ] 回收响应并返回给请求方
  - [ ] 超时处理（30s 默认）
- [ ] **日志与审计**
  - [ ] 结构化日志（slog）
  - [ ] 操作审计日志记录
  - [ ] 异常告警（企业微信 webhook）

**交付物**：`backend/` 完整 API 服务 + WebSocket 中转

#### 1.4 Cloudflare Tunnel 部署（Day 5-6）

```
目标：云端部署 Cloudflared，打通公网访问
```

- [ ] 云服务器安装 `cloudflared`
- [ ] 登录 Cloudflare（`cloudflared tunnel login`）
- [ ] 创建 Named Tunnel：`cinaseek-production`
- [ ] 配置 Ingress 规则：
  ```yaml
  # ~/.cloudflared/config.yml
  tunnel: cinaseek-production
  credentials-file: /root/.cloudflared/<TUNNEL_ID>.json
  
  ingress:
    - hostname: cinaseek.ai
      service: http://localhost:3000    # 前端
    - hostname: api.cinaseek.ai
      service: http://localhost:8080    # 后端 API
    - hostname: ws.cinaseek.ai
      service: ws://localhost:8081      # WebSocket
    - service: http_status:404
  ```
- [ ] DNS CNAME 配置：
  - `cinaseek.ai` → `<TUNNEL_ID>.cfargotunnel.com`
  - `api.cinaseek.ai` → `<TUNNEL_ID>.cfargotunnel.com`
  - `ws.cinaseek.ai` → `<TUNNEL_ID>.cfargotunnel.com`
- [ ] 开启 HTTPS 自动加密（Cloudflare Edge → Origin）
- [ ] 配置速率限制（100 req/min/user）
- [ ] 配置 WAF 规则
- [ ] 验证 Token 隔离：用户端无 Cloudflare 密钥
- [ ] Systemd 守护进程配置（开机自启）

**交付物**：`cinaseek.ai` + `api.cinaseek.ai` 可公网访问

#### 1.5 Web 管理面板功能联调（Day 6-7）

```
目标：前端所有页面与后端 API 打通
```

- [ ] 前端对接后端 API：
  - [ ] `GET /api/v1/vm/list` → VMList.vue
  - [ ] `POST /api/v1/vm/create` → CreateVM.vue
  - [ ] `GET /api/v1/vm/detail` → VMDetail.vue
  - [ ] `POST /api/v1/mount/add` → MountManager.vue
  - [ ] `GET /api/v1/openclaw/status` → Deploy/Monitor.vue
  - [ ] `GET /api/v1/remote/status` → RemoteAccess.vue
- [ ] WebShell 终端联调：
  - [ ] xterm.js ↔ WebSocket ↔ CinaClaw PTY
  - [ ] 多标签页切换
  - [ ] 断线重连
- [ ] OpenClaw 一键部署联调：
  - [ ] 选择虚拟机 → 一键安装 → 进度展示
  - [ ] 工作空间目录挂载
- [ ] CinaToken OAuth 登录联调：
  - [ ] Login.vue → CinaToken 授权 → 回调 → Token 存储
- [ ] ECharts 监控图表实时数据

**交付物**：全功能 Web 面板可操作

---

### 阶段二：安全优化（5天）

#### 2.1 安全加固（Day 8-9）

- [ ] **Token 隔离验证**
  - [ ] Wireshark 抓包确认：用户端流量无 Cloudflare Token
  - [ ] 代码审计：用户端代码无 Cloudflare 配置/密钥
  - [ ] 配置文件检查：`.env` 中无 `CLOUDFLARE_*` 变量
- [ ] **传输加密**
  - [ ] Cloudflare Edge → Origin: TLS 1.3（Full Strict）
  - [ ] WebSocket 连接: `wss://` 强制加密
  - [ ] gRPC 通信: Unix Socket（本地，无网络暴露）
- [ ] **访问控制**
  - [ ] IP 黑白名单（CIDR 支持）
  - [ ] 单用户并发连接限制（免费版 1，专业版 5）
  - [ ] 请求频率限制（令牌桶，100 req/min）
  - [ ] JWT Token 过期自动刷新
- [ ] **鉴权测试**
  - [ ] 伪造 Token 请求 → 401
  - [ ] 过期 Token 请求 → 401 + 自动刷新
  - [ ] 无 Token 请求 → 401
  - [ ] 越权访问他人虚拟机 → 403

#### 2.2 OpenClaw 专属适配（Day 9-10）

- [ ] 虚拟机配置优化
  - [ ] OpenClaw 模板：4 CPU / 8G 内存 / 50G 磁盘
  - [ ] Ubuntu 22.04 LTS 预装依赖
- [ ] 一键部署脚本
  - [ ] `install-openclaw.sh`：安装 Node.js + OpenClaw + 依赖
  - [ ] 自动配置工作空间 `/root/.openclaw/workspace`
  - [ ] 自动挂载模型目录
- [ ] 数据持久化
  - [ ] 工作空间目录挂载（宿主机 ↔ 虚拟机）
  - [ ] 配置文件自动备份
- [ ] OpenClaw 服务管理
  - [ ] `systemctl enable openclaw` 自启
  - [ ] 资源占用监控
  - [ ] 日志查看

#### 2.3 并发性能优化（Day 10-12）

- [ ] Go 服务优化
  - [ ] Goroutine 池（ants 库，默认 100 worker）
  - [ ] WebSocket 连接池（分桶管理，按用户 ID）
  - [ ] gRPC 连接复用（连接池，避免频繁建连）
  - [ ] PostgreSQL 连接池（MaxOpen=50, MaxIdle=10）
  - [ ] Redis Pipeline 批量操作
- [ ] 内存优化
  - [ ] pprof 内存分析
  - [ ] 减少不必要的内存拷贝
  - [ ] 流式响应（避免大包缓冲）
- [ ] 低配设备适配
  - [ ] 用户端 Go 客户端内存 < 50MB
  - [ ] WebSocket 断线重连轻量化
  - [ ] 前端懒加载

**性能目标**：

| 指标 | 目标 | 说明 |
|------|------|------|
| API 响应延迟 | < 100ms | P99 |
| WebSocket 连接延迟 | < 50ms | 建连时间 |
| 命令执行延迟 | < 100ms | WebShell 输入到输出 |
| 并发连接支持 | 1000+ | 免费版 |
| 请求转发吞吐 | 500+ req/s | 单实例 |

---

### 阶段三：打包测试（4-9天）

#### 3.1 全平台打包（Day 13-16）

**Windows 安装包**
- [ ] CinaClaw Windows 版（Hyper-V 后端）
- [ ] Go 客户端编译为 `.exe`
- [ ] Inno Setup/WiX 打包为安装包
- [ ] 集成：CinaClaw + Go 客户端 + 开机自启服务
- [ ] 安装向导：品牌 LOGO、协议页、路径选择
- [ ] Windows Defender 白名单测试

**macOS 安装包**
- [ ] CinaClaw macOS 版（Apple Silicon + Intel 双架构）
- [ ] Go 客户端编译（Universal Binary）
- [ ] `.dmg` 镜像打包
- [ ] 签名 + 公证（Apple Developer 账号）
- [ ] LaunchAgent 自启配置

**Linux 安装包**
- [ ] CinaClaw Linux 版（QEMU/KVM 后端）
- [ ] `.deb` 包（Ubuntu 20.04/22.04/24.04）
- [ ] `.rpm` 包（CentOS/RHEL 8/9）
- [ ] Systemd 服务托管
- [ ] Snap 包（`sudo snap install cinaseek`）

**品牌素材**
- [ ] 爪印 LOGO 替换所有 Multipass/CinaClaw 原始图标
- [ ] 安装包图标
- [ ] 系统托盘图标
- [ ] 关于页面品牌信息

#### 3.2 全场景测试（Day 16-19）

**功能测试**
| 功能 | 测试场景 | 预期 |
|------|----------|------|
| 本地管理 | 创建/启动/停止/删除虚拟机 | 正常 |
| 远程管理 | 通过 cinaseek.ai 远程操作 | 正常 |
| WebShell | 浏览器终端执行命令 | 延迟 < 100ms |
| 目录挂载 | 宿主机 ↔ VM 双向挂载 | 数据一致 |
| OpenClaw | 一键部署 + 服务管理 | 正常 |
| 快照 | 创建/恢复/删除 | 正常 |
| 监控 | CPU/内存/磁盘实时图表 | 数据准确 |

**安全测试**
- [ ] Wireshark 抓包：确认无 Cloudflare Token 泄露
- [ ] 鉴权失效：伪造/过期/空 Token 全部 401
- [ ] 恶意访问：SQL 注入、XSS、CSRF 全部拦截
- [ ] 越权测试：用户 A 无法操作用户 B 的虚拟机

**并发测试**
- [ ] 免费版：1000 用户同时在线，响应 P99 < 500ms
- [ ] 专业版：5000 用户同时在线，响应 P99 < 300ms
- [ ] WebSocket 长连接稳定性：24h 无断连

**网络测试**
- [ ] 家庭宽带（电信/联通/移动）
- [ ] 公司内网（NAT/防火墙）
- [ ] 4G/5G 移动网络
- [ ] 无公网 IP 环境（纯内网）

---

### 阶段四：部署上线（3天）

#### 4.1 云端生产环境（Day 20）

- [ ] **Cloudflare 配置**
  - [ ] 域名 `cinaseek.ai` NameServer 托管
  - [ ] 生产隧道 `cinaseek-production` 配置
  - [ ] DNS 解析：
    - `cinaseek.ai` → 前端（Vite build → Nginx）
    - `api.cinaseek.ai` → 后端 API（Go Gin :8080）
    - `ws.cinaseek.ai` → WebSocket 服务（Go :8081）
  - [ ] SSL：Full (Strict) + 15年 Origin CA
  - [ ] 缓存策略：静态资源 7天，API 不缓存
  - [ ] Brotli 压缩开启
  - [ ] WAF 规则配置

- [ ] **Go 中转服务部署**
  - [ ] Docker 镜像构建（`cinagroup/cinaseek:latest`）
  - [ ] K8s Deployment（2 副本 + HPA）
  - [ ] ConfigMap / Secret 配置
  - [ ] Nginx Ingress（WebSocket 支持）
  - [ ] 健康检查（liveness + readiness）

- [ ] **数据库部署**
  - [ ] PostgreSQL Schema `cinaseek` 初始化
  - [ ] 13 张数据表迁移
  - [ ] 主从复制验证
  - [ ] 自动备份（每日全量 + WAL 归档）

- [ ] **进程守护**
  - [ ] Cloudflared: Systemd 管理
  - [ ] Go 服务: K8s Deployment 自动重启
  - [ ] 日志切割: logrotate

#### 4.2 用户端分发（Day 21）

- [ ] **下载页**：`cinaseek.ai/download`
  - [ ] Windows (.exe)
  - [ ] macOS (.dmg)
  - [ ] Linux (.deb / .rpm / snap)
- [ ] **文档**：
  - [ ] 快速开始指南（5分钟上手）
  - [ ] OpenClaw 部署教程
  - [ ] 常见问题 FAQ
  - [ ] API 文档（Swagger UI）
- [ ] **版本更新**
  - [ ] 用户端自动检测 GitHub Release
  - [ ] 增量更新机制
  - [ ] 回滚支持
- [ ] **反馈渠道**
  - [ ] GitHub Issues
  - [ ] 企业微信群

#### 4.3 商用方案落地

| 版本 | 价格 | 实例数 | 并发连接 | 域名 | 技术支持 |
|------|------|--------|----------|------|----------|
| **免费版** | 免费 | 1 | 1 | `*.cinaseek.ai` 共享 | GitHub Issues |
| **专业版** | ¥49/月 | 5 | 5 | `{user}.cinaseek.ai` | 邮件 + 工单 |
| **企业版** | ¥299/月 | 无限 | 无限 | 自定义域名 | 专属技术支持 |

**免费版限制**：
- 1 个虚拟机实例
- 1 个并发 WebSocket 连接
- 社区支持
- Cloudflare 免费隧道

**专业版增值**：
- 多实例管理
- 专属子域名
- 快照管理
- 高并发支持
- 优先响应

**企业版增值**：
- 无限实例
- 独立 Cloudflare 隧道
- 自定义域名绑定
- SLA 保障（99.9%）
- 7×24 技术支持
- 定制化功能

---

## 六、运维与售后方案（长期）

### 6.1 日常运维

- [ ] **监控告警**
  - [ ] Uptime Kuma：API/WebSocket/Tunnel 状态
  - [ ] 企业微信告警：服务异常、连接断开、资源告警
  - [ ] 邮件告警：数据库备份状态
- [ ] **版本更新**
  - [ ] Cloudflared 版本更新
  - [ ] Go 服务版本更新（滚动部署）
  - [ ] CinaClaw 引擎版本更新（跟踪上游）
- [ ] **资源优化**
  - [ ] 清理无效用户连接（24h 无活动）
  - [ ] PostgreSQL VACUUM 定期维护
  - [ ] Redis 内存清理

### 6.2 售后支持

- [ ] 使用教程（图文 + 视频）
- [ ] FAQ 知识库
- [ ] 远程排查工具（WebShell 日志导出）
- [ ] OpenClaw 专属技术支持
- [ ] 用户反馈收集 → 产品迭代

### 6.3 合规与风险控制

- [ ] 严格遵循 Cloudflare 商用条款，杜绝流量滥用
- [ ] GPLv3 协议合规：公开 CinaClaw 修改源码，标注原版权
- [ ] 用户数据隔离：Schema 级别隔离，杜绝跨用户数据泄露
- [ ] 隐私政策 + 用户协议
- [ ] 数据保留策略（用户注销后 30 天删除）

---

## 七、项目里程碑

| 阶段 | 时间 | 核心交付物 | 状态 |
|------|------|------------|------|
| 前期准备 | Day 0-3 | 域名、环境、品牌素材、合规梳理 | 🔄 进行中 |
| 核心开发 | Day 4-10 | CinaClaw 验证 + Go 客户端 + 云端服务 + Tunnel + 前端联调 | ⏳ 待开始 |
| 安全优化 | Day 11-15 | 安全加固 + OpenClaw 适配 + 并发优化 | ⏳ 待开始 |
| 打包测试 | Day 16-24 | 全平台安装包 + 全场景测试报告 | ⏳ 待开始 |
| 部署上线 | Day 25-27 | 云端生产环境 + 用户端分发 + 商用方案 | ⏳ 待开始 |
| **正式上线** | **Day 28** | **cinaseek.ai 正式运营** | ⏳ 待开始 |

---

## 八、核心风险与规避

| 风险 | 等级 | 规避方案 |
|------|------|----------|
| Cloudflare Token 泄露 | 🔴 高 | Token 仅云端加密存储，用户端零密钥，全程反向 WebSocket 连接 |
| Multipass GPLv3 合规 | 🟡 中 | 标注原版权，公开 CinaClaw 修改源码，遵循 GPLv3 商用规则 |
| 海量用户并发瓶颈 | 🟡 中 | 多账号隔离，K8s HPA 弹性扩容，付费版独享，免费版限流 |
| 用户端网络适配 | 🟢 低 | 纯出站 WebSocket 连接，无需端口映射，自动重连，适配所有内网 |
| Cloudflare 账号封禁 | 🟡 中 | 遵守商用条款，不滥用免费额度，准备备用域名 |
| 数据丢失 | 🔴 高 | PostgreSQL 主从复制 + 每日备份 + WAL 归档 + 快照保护 |

---

## 九、核心优势总结

| 优势 | 说明 |
|------|------|
| 🛡️ 绝对安全 | 用户端无敏感密钥，Cloudflare 账号完全隔离 |
| ⚡ 零配置远程 | 无需懂内网穿透、端口映射，开箱即用 |
| 💰 商用可行 | 分层定价，支持海量用户，GPLv3 合规 |
| 🤖 OpenClaw 专属 | 完美适配 AI 开发环境，一键部署 |
| 🌐 全球加速 | Cloudflare 全球 CDN 节点，低延迟访问 |
| 🔓 开源透明 | CinaClaw 引擎 GPLv3 开源，代码可审计 |

---

*CinaSeek - 你的云端开发工作室* 🐾
