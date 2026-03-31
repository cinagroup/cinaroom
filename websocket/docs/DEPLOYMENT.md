# MultiPass WebSocket 部署文档

## 目录

1. [快速开始](#快速开始)
2. [环境准备](#环境准备)
3. [配置 Cloudflare Tunnel](#配置-cloudflare-tunnel)
4. [Docker 部署](#docker-部署)
5. [手动部署](#手动部署)
6. [验证与测试](#验证与测试)
7. [故障排查](#故障排查)

## 快速开始

```bash
# 1. 克隆项目
cd /root/.openclaw/workspace/multipass-websocket

# 2. 配置环境变量
cp deploy/.env.example .env
# 编辑 .env 文件，填入 Cloudflare Tunnel 配置

# 3. 启动服务
cd deploy
docker-compose up -d

# 4. 查看日志
docker-compose logs -f
```

## 环境准备

### 系统要求

- Linux (Ubuntu 20.04+ 推荐)
- Docker 20.10+
- Docker Compose 2.0+
- 2GB+ RAM
- 2+ CPU 核心

### 安装 Docker

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# 验证安装
docker --version
docker-compose --version
```

## 配置 Cloudflare Tunnel

### 1. 创建 Tunnel

```bash
# 安装 cloudflared
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64
chmod +x cloudflared-linux-amd64
sudo mv cloudflared-linux-amd64 /usr/local/bin/cloudflared

# 登录 Cloudflare
cloudflared tunnel login

# 创建 tunnel
cloudflared tunnel create multipass
```

### 2. 获取配置信息

```bash
# 获取 Tunnel ID
cloudflared tunnel list

# 导出凭证
cloudflared tunnel creds --tunnel-id <TUNNEL_ID> > deploy/cloudflared-creds.json
```

### 3. 配置 DNS

在 Cloudflare Dashboard 中配置 DNS CNAME 记录：

```
terminal.example.com -> <TUNNEL_ID>.cfargotunnel.com
relay.example.com -> <TUNNEL_ID>.cfargotunnel.com
```

### 4. 更新环境变量

编辑 `.env` 文件：

```bash
CLOUDFLARE_TUNNEL_ID=<你的 Tunnel ID>
CLOUDFLARE_HOSTNAME=example.com
CLOUDFLARE_SECRET=<你的 Secret>
```

## Docker 部署

### 1. 启动服务

```bash
cd deploy

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f websocket-terminal
docker-compose logs -f cloud-relay
docker-compose logs -f cloudflared
```

### 2. 服务端口

| 服务 | 容器端口 | 主机端口 | 说明 |
|------|---------|---------|------|
| websocket-terminal | 8080 | 8080 | WebSocket 终端服务 |
| cloud-relay | 8081 | 8081 | 云端中转服务 |
| cloudflared | - | - | Cloudflare Tunnel（无暴露端口） |
| nginx | 80/443 | 80/443 | 反向代理（可选） |
| redis | 6379 | 6379 | Redis 存储（可选） |

### 3. 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

### 4. 更新服务

```bash
# 重新构建并启动
docker-compose up -d --build

# 只更新特定服务
docker-compose up -d --build websocket-terminal
```

## 手动部署

### 1. 编译后端服务

```bash
cd backend

# 下载依赖
go mod download

# 编译终端服务
go build -o terminal .

# 编译中转服务
go build -o cloud-relay .
```

### 2. 运行服务

```bash
# 设置环境变量
export CLOUDFLARE_TUNNEL_ID=your-id
export CLOUDFLARE_HOSTNAME=example.com

# 运行终端服务
./terminal

# 运行中转服务（新终端）
./cloud-relay
```

### 3. 使用 systemd 管理

创建 `/etc/systemd/system/websocket-terminal.service`:

```ini
[Unit]
Description=WebSocket Terminal Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/multipass-websocket/backend
EnvironmentFile=/opt/multipass-websocket/.env
ExecStart=/opt/multipass-websocket/backend/terminal
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable websocket-terminal
sudo systemctl start websocket-terminal
sudo systemctl status websocket-terminal
```

## 验证与测试

### 1. 健康检查

```bash
# 终端服务
curl http://localhost:8080/health

# 中转服务
curl http://localhost:8081/health
```

### 2. 运行集成测试

```bash
cd tests

# 运行集成测试
./integration_test.sh

# 运行 Python WebSocket 测试
python3 websocket_test.py --url ws://localhost:8080 --token your-token

# 压力测试
RUN_STRESS_TEST=true ./integration_test.sh
```

### 3. WebSocket 连接测试

```bash
# 使用 websocat
websocat "ws://localhost:8080/api/v1/ws/terminal?token=your-token"

# 使用 wscat
wscat -c "ws://localhost:8080/api/v1/ws/terminal?token=your-token"
```

## 故障排查

### 常见问题

#### 1. Cloudflare Tunnel 连接失败

```bash
# 检查 tunnel 状态
cloudflared tunnel info <TUNNEL_ID>

# 查看 tunnel 日志
docker-compose logs cloudflared

# 验证凭证
cat deploy/cloudflared-creds.json
```

#### 2. WebSocket 连接断开

- 检查 Nginx WebSocket 配置
- 验证防火墙规则
- 检查超时设置

```bash
# 查看连接数
docker-compose exec websocket-terminal netstat -an | grep ESTABLISHED | wc -l
```

#### 3. 内存溢出

```bash
# 查看容器资源使用
docker stats

# 调整资源限制
# 编辑 docker-compose.yml 中的 deploy.resources
```

#### 4. 日志查看

```bash
# 实时日志
docker-compose logs -f

# 最近 100 行
docker-compose logs --tail=100 websocket-terminal

# 导出日志
docker-compose logs > deployment.log
```

### 性能调优

#### 1. 增加连接数限制

编辑 `docker-compose.yml`:

```yaml
environment:
  - MAX_CONNECTIONS=500  # 增加最大连接数
```

#### 2. 调整超时时间

```yaml
environment:
  - SESSION_TIMEOUT_MINUTES=60
  - TOKEN_EXPIRY_HOURS=48
```

#### 3. 启用 Redis 集群

取消 `docker-compose.yml` 中 Redis 服务的注释，并更新后端配置使用 Redis 存储会话。

## 安全建议

1. **使用 HTTPS**: 生产环境必须启用 HTTPS
2. **限制来源**: 配置 `ALLOWED_ORIGINS` 环境变量
3. **定期更新 Token**: 设置合理的 `TOKEN_EXPIRY_HOURS`
4. **启用速率限制**: 配置 `RATE_LIMIT_REQUESTS_PER_MINUTE`
5. **监控日志**: 定期检查访问日志和错误日志
6. **防火墙配置**: 只开放必要的端口

## 监控与告警

### Prometheus 指标

服务暴露以下指标（需要启用 metrics 端点）:

- `websocket_connections_total`: 总连接数
- `websocket_active_connections`: 活跃连接数
- `websocket_messages_total`: 消息总数
- `forward_requests_total`: 转发请求数

### Grafana 仪表板

导入提供的 Grafana 仪表板 JSON 文件（`deploy/grafana-dashboard.json`）进行可视化监控。

## 备份与恢复

### 备份 Redis 数据

```bash
docker-compose exec redis redis-cli SAVE
cp deploy/redis-data/dump.rdb backup-$(date +%Y%m%d).rdb
```

### 恢复 Redis 数据

```bash
cp backup-20240101.rdb deploy/redis-data/dump.rdb
docker-compose restart redis
```

## 升级指南

### 从 v1.x 升级到 v2.x

1. 备份当前配置和数据
2. 查看 CHANGELOG.md 了解变更
3. 更新 docker-compose.yml
4. 重新构建并启动服务
5. 验证功能正常

```bash
# 备份
docker-compose down
cp -r deploy deploy.backup

# 升级
git pull
docker-compose up -d --build

# 验证
./tests/integration_test.sh
```
