# Multipass Backend 快速开始指南

## 5 分钟快速上手

### 前置要求

- Go 1.21+
- PostgreSQL 12+
- Git

### 步骤 1: 克隆项目

```bash
cd /root/.openclaw/workspace
# 项目已存在
cd multipass-backend
```

### 步骤 2: 安装依赖

```bash
go mod tidy
```

### 步骤 3: 创建数据库

```bash
# 连接到 PostgreSQL
psql -U postgres

# 创建数据库
CREATE DATABASE multipass;

# 退出
\q
```

### 步骤 4: 配置环境变量

```bash
# 复制示例配置
cp .env.example .env

# 编辑配置（可选，使用默认值即可）
vim .env
```

默认配置：
- 端口：8080
- 数据库：localhost:5432
- 数据库用户：postgres
- 数据库密码：postgres

### 步骤 5: 运行服务

```bash
# 方式 1: 直接运行
go run cmd/server/main.go

# 方式 2: 使用 Make
make run

# 方式 3: 开发模式（热重载）
make dev
```

### 步骤 6: 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 获取版本
curl http://localhost:8080/api/v1/system/version

# 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "Admin123",
    "confirm_password": "Admin123"
  }'

# 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123"
  }'
```

## 使用 Docker 部署

### 步骤 1: 启动所有服务

```bash
docker-compose up -d
```

### 步骤 2: 查看日志

```bash
docker-compose logs -f backend
```

### 步骤 3: 测试

```bash
curl http://localhost:8080/health
```

### 步骤 4: 停止服务

```bash
docker-compose down
```

## 使用部署脚本（生产环境）

```bash
# 赋予执行权限
chmod +x scripts/deploy.sh

# 生产环境部署
sudo ./scripts/deploy.sh prod

# Docker 部署
sudo ./scripts/deploy.sh docker
```

## 验证安装

### 1. 检查服务状态

```bash
# Systemd 服务
systemctl status multipass-backend

# Docker 容器
docker ps | grep multipass
```

### 2. 检查端口

```bash
netstat -tlnp | grep 8080
# 或
ss -tlnp | grep 8080
```

### 3. 测试 API

```bash
# 健康检查应该返回 healthy
curl http://localhost:8080/health

# 预期响应:
# {"code":0,"message":"healthy","data":{"status":"healthy","database":"connected",...}}
```

## 常见问题

### 问题 1: 端口已被占用

**错误**: `bind: address already in use`

**解决**:
```bash
# 查找占用端口的进程
lsof -i :8080

# 杀死进程
kill -9 <PID>

# 或修改端口
export SERVER_PORT=8081
```

### 问题 2: 数据库连接失败

**错误**: `failed to connect database`

**解决**:
```bash
# 检查 PostgreSQL 是否运行
systemctl status postgresql

# 检查数据库是否存在
psql -U postgres -l | grep multipass

# 检查连接配置
cat .env | grep DB_
```

### 问题 3: 编译失败

**错误**: `package not found`

**解决**:
```bash
# 清理并重新下载依赖
go clean -modcache
go mod tidy

# 检查 Go 版本
go version
# 应该是 1.21 或更高
```

### 问题 4: JWT Token 过期

**错误**: `认证失败或 Token 已过期`

**解决**:
```bash
# Token 默认 24 小时过期
# 重新登录获取新 Token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin123"}'
```

## 下一步

1. **创建第一个虚拟机**: 参考 API 文档的虚拟机管理章节
2. **配置 OpenClaw**: 部署 OpenClaw 到你的虚拟机
3. **设置目录挂载**: 配置工作空间同步
4. **启用远程访问**: 配置 IP 白名单和远程访问

## 获取帮助

- 查看完整 API 文档：`api/API.md`
- 查看 README: `README.md`
- 查看部署脚本：`scripts/deploy.sh`

## 开发模式

### 热重载（推荐）

```bash
# 安装 air
go install github.com/air-verse/air@latest

# 运行
air
```

### 运行测试

```bash
go test ./... -v
```

### 代码格式化

```bash
go fmt ./...
go vet ./...
```

## 生产环境优化

### 1. 修改 JWT 密钥

```bash
# 生成随机密钥
openssl rand -base64 32

# 更新 .env 文件
JWT_SECRET=<生成的密钥>
```

### 2. 配置 HTTPS

使用 Nginx 反向代理：

```bash
# 安装 Nginx
apt-get install nginx

# 配置反向代理
cp nginx.conf /etc/nginx/nginx.conf

# 重启 Nginx
systemctl restart nginx
```

### 3. 配置防火墙

```bash
# 允许 8080 端口
ufw allow 8080/tcp

# 如果使用 Nginx
ufw allow 80/tcp
ufw allow 443/tcp
```

### 4. 设置日志轮转

创建 `/etc/logrotate.d/multipass-backend`:

```
/var/log/multipass-backend/*.log {
    daily
    rotate 7
    compress
    delaycompress
    notifempty
    create 0640 root root
}
```

## 监控和维护

### 查看日志

```bash
# Systemd 日志
journalctl -u multipass-backend -f

# 应用日志
tail -f /var/log/multipass-backend/app.log
```

### 备份数据库

```bash
pg_dump -U postgres multipass > backup_$(date +%Y%m%d).sql
```

### 恢复数据库

```bash
psql -U postgres multipass < backup_20240101.sql
```

### 更新版本

```bash
# 拉取最新代码
git pull

# 重新编译
make build

# 重启服务
sudo systemctl restart multipass-backend
```

---

祝你使用愉快！🎉
