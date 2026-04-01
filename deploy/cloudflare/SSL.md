# SSL 证书配置指南 — CinaSeek

> Cloudflare 提供多层 SSL 保护。本文档说明如何配置 Full (Strict) 模式 + Origin CA 证书。

---

## SSL 架构

```
用户浏览器 ←──SSL──→ Cloudflare 边缘 ←──Origin CA──→ 你的服务器
   (自动HTTPS)        (免费边缘证书)         (自签 Origin CA)
```

- **用户 → Cloudflare**：Cloudflare 自动提供边缘证书（Universal SSL）
- **Cloudflare → 源站**：使用 Cloudflare Origin CA 证书加密

---

## 方案一：Cloudflare Origin CA（推荐）

### 1. 申请 Origin CA 证书

1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. 选择域名 `cinaseek.ai`
3. 进入 **SSL/TLS → Origin Server**
4. 点击 **Create Certificate**
5. 配置：
   - 私钥类型：**RSA (2048)**
   - 证书有效期：**15 年**（最大）
   - 主机名：`cinaseek.ai, *.cinaseek.ai`
   - 证书签名请求：留空（Cloudflare 生成）
6. 点击 **Create**
7. **立即保存**显示的证书和私钥（私钥只显示一次！）

### 2. 安装证书到服务器

```bash
# 创建证书目录
sudo mkdir -p /etc/ssl/cinaseek

# 保存证书
sudo nano /etc/ssl/cinaseek/origin.pem
# 粘贴 Origin Certificate 内容

# 保存私钥
sudo nano /etc/ssl/cinaseek/origin-key.pem
# 粘贴 Private Key 内容

# 设置权限
sudo chmod 600 /etc/ssl/cinaseek/origin-key.pem
sudo chmod 644 /etc/ssl/cinaseek/origin.pem
sudo chown root:root /etc/ssl/cinaseek/*
```

### 3. 配置 Cloudflare SSL 模式

1. 进入 **SSL/TLS → Overview**
2. 选择 **Full (Strict)**

> ⚠️ 必须选 **Full (Strict)** 而不是 **Full**。Strict 模式会验证源站证书的有效性。

### 4. 配置 Nginx（如果使用 Nginx 反代）

```nginx
# /etc/nginx/sites-available/cinaseek

# 前端
server {
    listen 443 ssl http2;
    server_name cinaseek.ai;

    ssl_certificate     /etc/ssl/cinaseek/origin.pem;
    ssl_certificate_key /etc/ssl/cinaseek/origin-key.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $http_cf_connecting_ip;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# API
server {
    listen 443 ssl http2;
    server_name api.cinaseek.ai;

    ssl_certificate     /etc/ssl/cinaseek/origin.pem;
    ssl_certificate_key /etc/ssl/cinaseek/origin-key.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $http_cf_connecting_ip;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# WebSocket
server {
    listen 443 ssl;
    server_name ws.cinaseek.ai;

    ssl_certificate     /etc/ssl/cinaseek/origin.pem;
    ssl_certificate_key /etc/ssl/cinaseek/origin-key.pem;

    ssl_protocols TLSv1.2 TLSv1.3;

    location / {
        proxy_pass http://127.0.0.1:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $http_cf_connecting_ip;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket 超时
        proxy_read_timeout 86400s;
        proxy_send_timeout 86400s;
    }
}
```

```bash
# 启用配置
sudo ln -sf /etc/nginx/sites-available/cinaseek /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 5. 更新 Tunnel 配置

如果使用 Nginx + Full (Strict)，更新 `config.yml`：

```yaml
ingress:
  - hostname: cinaseek.ai
    service: https://localhost:443
    originRequest:
      noTLSVerify: false    # false = 验证 Origin CA 证书（Strict 模式）
      http2Origin: true

  # ... 其他规则同理改为 https://
```

如果 **不使用 Nginx**，Tunnel 直接连接后端服务，保持 `http://localhost` 即可，Cloudflare 边缘仍提供 HTTPS。

---

## 方案二：Let's Encrypt（不推荐用于 Tunnel）

Cloudflare Tunnel 场景下不推荐 Let's Encrypt，因为：
- Cloudflare 已提供边缘 SSL
- Origin CA 有效期 15 年，无需续期
- Let's Encrypt 需要公网可达的 80 端口（Tunnel 场景下不需要）

---

## 方案三：不使用 Nginx，Tunnel 直接转发

最简方案，无需 Nginx，无需 Origin CA：

```yaml
ingress:
  - hostname: cinaseek.ai
    service: http://localhost:3000    # 直接连后端
```

Cloudflare SSL 模式设为 **Flexible** 或 **Full**。

> 优点：配置最简单
> 缺点：Cloudflare → 源站无加密（仅限本机回环，风险较低）

---

## SSL 验证

```bash
# 检查边缘证书（外部）
curl -vI https://cinaseek.ai 2>&1 | grep -E "SSL|subject|issuer"

# 检查源站证书（内部）
openssl s_client -connect localhost:443 -servername cinaseek.ai </dev/null 2>/dev/null | openssl x509 -noout -subject -issuer -dates

# 检查 Origin CA 链
openssl verify -CAfile <(curl -s https://developers.cloudflare.com/ssl/static/origin_ca_rsa_root.pem) /etc/ssl/cinaseek/origin.pem
```

---

## 常见问题

| 问题 | 解决方案 |
|------|---------|
| 525 SSL Handshake Failed | Origin CA 证书未安装或 Nginx 未监听 443 |
| 526 Invalid SSL Certificate | 证书不匹配或已过期，或 SSL 模式应为 Full 而非 Strict |
| ERR_TOO_MANY_REDIRECTS | SSL 模式设为 Flexible 但 Nginx 又做了 HTTP→HTTPS 重定向 |

---

## Cloudflare Origin CA 根证书

验证时需要下载根证书：
- [Origin CA RSA Root](https://developers.cloudflare.com/ssl/static/origin_ca_rsa_root.pem)
- [Origin CA ECC Root](https://developers.cloudflare.com/ssl/static/origin_ca_ecc_root.pem)

---

## 参考链接

- [Cloudflare Origin CA](https://developers.cloudflare.com/ssl/origin-configuration/origin-ca/)
- [SSL 模式说明](https://developers.cloudflare.com/ssl/origin-configuration/ssl-modes/)
- [Nginx + Cloudflare](https://developers.cloudflare.com/ssl/origin-configuration/nginx/)
