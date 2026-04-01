# Cloudflare Tunnel 部署指南 — CinaSeek

> 域名：cinaseek.ai | 服务器：Ubuntu 24.04, 2核3.6G

---

## 目录

1. [架构概览](#架构概览)
2. [前置条件](#前置条件)
3. [安装 Cloudflared](#安装-cloudflared)
4. [登录 Cloudflare](#步骤-1登录-cloudflare)
5. [创建 Named Tunnel](#步骤-2创建-named-tunnel)
6. [配置 Ingress 规则](#步骤-3配置-ingress-规则)
7. [DNS CNAME 配置](#步骤-4dns-cname-配置)
8. [Systemd 服务配置](#步骤-5systemd-服务配置)
9. [SSL 配置](#步骤-6ssl-配置)
10. [安全加固](#步骤-7安全加固waf--速率限制)
11. [验证与排障](#验证与排障)
12. [一键部署](#一键部署脚本)

---

## 架构概览

```
用户 → Cloudflare CDN (Anycast)
         ├── cinaseek.ai     → Tunnel → localhost:3000  (前端)
         ├── api.cinaseek.ai  → Tunnel → localhost:8080  (API)
         └── ws.cinaseek.ai   → Tunnel → localhost:8081  (WebSocket)
                    ↓
            cloudflared (本机)
                    ↓
            Nginx / 直接转发到后端服务
```

**优势：**
- 无需开放公网端口（入站流量全部走 Tunnel）
- Cloudflare 自动提供 SSL 终止
- 自带 DDoS 防护和 WAF
- WebSocket 原生支持

---

## 前置条件

- [x] Cloudflare 账户，域名 `cinaseek.ai` 已添加到 Cloudflare（NS 已切换）
- [x] 服务器已安装 cloudflared（见下方）
- [x] 本地服务已运行：
  - 前端 `localhost:3000`
  - API `localhost:8080`
  - WebSocket `localhost:8081`

### 安装 Cloudflared

```bash
# 下载最新版
curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb \
  -o /tmp/cloudflared.deb

# 安装
sudo dpkg -i /tmp/cloudflared.deb

# 验证
cloudflared --version
# 应输出: cloudflared version 2026.x.x
```

---

## 步骤 1：登录 Cloudflare

```bash
cloudflared tunnel login
```

执行后会输出一个 URL，在浏览器中打开并授权 `cinaseek.ai` 域名。授权完成后会在 `~/.cloudflared/` 目录下生成 `cert.pem` 证书文件。

> ⚠️ `cert.pem` 包含账户级权限，**切勿提交到 Git**。已添加到 `.gitignore`。

---

## 步骤 2：创建 Named Tunnel

```bash
cloudflared tunnel create cinaseek-production
```

创建成功后会：
1. 生成 Tunnel UUID（例如 `a1b2c3d4-e5f6-7890-abcd-ef1234567890`）
2. 在 `~/.cloudflared/` 下生成凭证文件 `<TUNNEL_UUID>.json`

**记下 Tunnel UUID**，后续配置需要用到。

> 💡 也可以使用 `cloudflared tunnel list` 查看所有隧道及其 UUID。

---

## 步骤 3：配置 Ingress 规则

将本仓库的 `deploy/cloudflare/config.yml` 复制到 `/etc/cloudflared/`：

```bash
# 创建配置目录
sudo mkdir -p /etc/cloudflared

# 复制配置文件
sudo cp deploy/cloudflare/config.yml /etc/cloudflared/config.yml

# 复制凭证文件（从 ~/.cloudflared/ 复制）
sudo cp ~/.cloudflared/<TUNNEL_UUID>.json /etc/cloudflared/credentials.json
sudo chmod 600 /etc/cloudflared/credentials.json
```

### 配置文件说明 (`config.yml`)

```yaml
tunnel: cinaseek-production
credentials-file: /etc/cloudflared/credentials.json

ingress:
  # 前端 — Next.js / Nuxt
  - hostname: cinaseek.ai
    service: http://localhost:3000
    originRequest:
      noTLSVerify: false

  # API 后端
  - hostname: api.cinaseek.ai
    service: http://localhost:8080
    originRequest:
      noTLSVerify: false

  # WebSocket 服务
  - hostname: ws.cinaseek.ai
    service: ws://localhost:8081
    originRequest:
      noTLSVerify: false

  # 兜底规则（必须）
  - service: http_status:404
```

**验证配置文件语法：**

```bash
cloudflared tunnel ingress validate --config /etc/cloudflared/config.yml
```

> 💡 `noTLSVerify: false` 表示 Tunnel 到源站使用标准 TLS 验证。如果使用 Cloudflare Origin CA，保持 `false`。

---

## 步骤 4：DNS CNAME 配置

为每个子域名创建 CNAME 记录，指向 Tunnel：

```bash
# 前端
cloudflared tunnel route dns cinaseek-production cinaseek.ai

# API
cloudflared tunnel route dns cinaseek-production api.cinaseek.ai

# WebSocket
cloudflared tunnel route dns cinaseek-production ws.cinaseek.ai
```

以上命令会自动在 Cloudflare DNS 中创建 CNAME 记录：

| 子域名 | 类型 | 值 | 代理 |
|--------|------|------|------|
| `cinaseek.ai` | CNAME | `<TUNNEL_UUID>.cfargotunnel.com` | ☁️ 已代理 |
| `api.cinaseek.ai` | CNAME | `<TUNNEL_UUID>.cfargotunnel.com` | ☁️ 已代理 |
| `ws.cinaseek.ai` | CNAME | `<TUNNEL_UUID>.cfargotunnel.com` | ☁️ 已代理 |

> 也可以手动在 Cloudflare Dashboard → DNS 中添加 CNAME 记录。

---

## 步骤 5：Systemd 服务配置

```bash
# 复制 service 文件
sudo cp deploy/cloudflare/cloudflared.service /etc/systemd/system/cloudflared.service

# 加载并启动
sudo systemctl daemon-reload
sudo systemctl enable cloudflared
sudo systemctl start cloudflared

# 查看状态
sudo systemctl status cloudflared

# 查看日志
sudo journalctl -u cloudflared -f
```

### 服务文件说明

```ini
[Unit]
Description=Cloudflared Tunnel for CinaSeek
After=network.target

[Service]
Type=notify
ExecStart=/usr/local/bin/cloudflared --config /etc/cloudflared/config.yml tunnel run
Restart=on-failure
RestartSec=5
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

- `Type=notify`：cloudflared 会在就绪后通知 systemd
- `Restart=on-failure`：异常退出时自动重启
- `RestartSec=5`：5 秒后重启
- `TimeoutStartSec=0`：不限制启动超时（Tunnel 连接可能较慢）

---

## 步骤 6：SSL 配置

> Cloudflare Tunnel 默认提供 Cloudflare → 用户 的 SSL 加密（边缘证书）。
> 如果需要 Tunnel → 源站 也加密（Full Strict 模式），参见 [SSL.md](./SSL.md)。

### 推荐模式

| 模式 | 说明 | 推荐场景 |
|------|------|---------|
| **Flexible** | Cloudflare → 源站 HTTP | 开发/测试环境 |
| **Full** | Cloudflare → 源站 HTTPS（不验证证书） | 快速部署 |
| **Full (Strict)** | Cloudflare → 源站 HTTPS（验证 Origin CA） | ✅ **生产环境推荐** |

**生产环境建议使用 Full (Strict) + Cloudflare Origin CA。** 详见 [SSL.md](./SSL.md)。

---

## 步骤 7：安全加固（WAF + 速率限制）

### 7.1 WAF 规则（Cloudflare Dashboard）

路径：**Security → WAF → Custom rules**

#### 规则 1：阻止已知恶意请求

```
Expression: (http.request.uri.path contains "/wp-admin") or 
            (http.request.uri.path contains "/xmlrpc.php") or
            (http.request.uri.path contains "/.env")
Action: Block
```

#### 规则 2：API 速率限制

路径：**Security → WAF → Rate limiting rules**

```
URL: api.cinaseek.ai/*
Requests: 100 / 10 seconds
Action: Block for 60 seconds
```

#### 规则 3：WebSocket 连接保护

```
URL: ws.cinaseek.ai/*
Requests: 30 / 1 minute
Action: Challenge
```

### 7.2 安全设置建议

在 **Security → Settings** 中：

| 设置 | 推荐值 |
|------|--------|
| Security Level | Medium |
| Challenge Passage | 30 minutes |
| Browser Integrity Check | ✅ On |
| Privacy Support | ✅ On |

### 7.3 Bot Management

在 **Security → Bots** 中：

- 启用 **Bot Fight Mode**
- 考虑使用 **Super Bot Fight Mode**（付费计划）

### 7.4 防火墙关闭公网入站端口

使用 Cloudflare Tunnel 后，无需开放 80/443 端口。建议：

```bash
# 只保留 SSH（可改为非标准端口）
sudo ufw default deny incoming
sudo ufw allow 22/tcp       # SSH（建议改为非标准端口）
sudo ufw allow out 53/tcp   # DNS
sudo ufw allow out 53/udp   # DNS
sudo ufw allow out 443/tcp  # HTTPS outbound（Tunnel 需要）
sudo ufw allow out 80/tcp   # HTTP outbound
sudo ufw enable
```

> ⚠️ 确保 SSH 端口已允许后再启用，否则会锁死自己。

---

## 验证与排障

### 验证连通性

```bash
# 检查 Tunnel 状态
cloudflared tunnel info cinaseek-production

# 测试各域名
curl -I https://cinaseek.ai
curl -I https://api.cinaseek.ai
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" -H "Sec-WebSocket-Key: test" \
  https://ws.cinaseek.ai
```

### 常见问题

| 问题 | 排查方法 |
|------|---------|
| Tunnel 无法启动 | `journalctl -u cloudflared -n 50` 查看日志 |
| DNS 未解析 | `dig cinaseek.ai` 检查 CNAME 记录 |
| 502 Bad Gateway | 检查本地服务是否运行：`curl localhost:3000` |
| WebSocket 断连 | 确认 ws.cinaseek.ai 的 service 是 `ws://` 协议 |
| 证书错误 | 确认 Cloudflare SSL 模式设置 |

### 日志查看

```bash
# 实时日志
sudo journalctl -u cloudflared -f

# 最近 100 行
sudo journalctl -u cloudflared -n 100

# 按时间过滤
sudo journalctl -u cloudflared --since "1 hour ago"
```

---

## 一键部署脚本

使用交互式脚本自动完成所有步骤：

```bash
chmod +x deploy/cloudflare/setup-tunnel.sh
sudo ./deploy/cloudflare/setup-tunnel.sh
```

详见 [setup-tunnel.sh](./setup-tunnel.sh)。

---

## 文件清单

```
deploy/cloudflare/
├── SETUP.md              # 本文档
├── SSL.md                # SSL 证书配置指南
├── config.yml            # Tunnel 配置模板
├── cloudflared.service   # Systemd 服务文件
└── setup-tunnel.sh       # 一键部署脚本
```

---

## 参考链接

- [Cloudflare Tunnel 官方文档](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/)
- [cloudflared GitHub](https://github.com/cloudflare/cloudflared)
- [Cloudflare WAF 配置](https://developers.cloudflare.com/waf/)
