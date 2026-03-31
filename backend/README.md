# Multipass Backend API

基于 Go + Gin 的 Multipass 远程管理平台后端 API 服务。

## 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Gin v1.9+
- **ORM**: GORM v1.25+
- **数据库**: PostgreSQL 12+
- **认证**: JWT (JSON Web Token)
- **密码加密**: bcrypt

## 项目结构

```
multipass-backend/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── internal/
│   ├── config/              # 配置管理
│   ├── handler/             # HTTP 处理器
│   ├── middleware/          # 中间件
│   ├── model/               # 数据模型
│   ├── repository/          # 数据库操作
│   ├── service/             # 业务逻辑
│   └── utils/               # 工具函数
├── pkg/
│   ├── response/            # 响应封装
│   └── validator/           # 参数验证
├── api/                     # API 文档
├── scripts/                 # 脚本文件
├── tests/                   # 测试文件
├── go.mod                   # Go 模块定义
├── go.sum                   # 依赖校验
└── README.md                # 项目说明
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- PostgreSQL 12+
- Multipass (用于虚拟机管理)

### 2. 安装依赖

```bash
cd multipass-backend
go mod tidy
```

### 3. 配置环境变量

```bash
# 服务器配置
export SERVER_PORT=8080

# 数据库配置
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=multipass
export DB_SSLMODE=disable

# JWT 配置
export JWT_SECRET=your-secret-key-change-in-production
export JWT_EXPIRE=24h
```

### 4. 创建数据库

```bash
psql -U postgres -c "CREATE DATABASE multipass;"
```

### 5. 运行服务

```bash
go run cmd/server/main.go
```

### 6. 构建可执行文件

```bash
go build -o multipass-backend cmd/server/main.go
```

## API 接口文档

### 认证模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/v1/auth/register | 用户注册 | ❌ |
| POST | /api/v1/auth/login | 用户登录 | ❌ |
| POST | /api/v1/auth/reset-pwd | 重置密码 | ❌ |
| POST | /api/v1/auth/logout | 用户登出 | ✅ |
| GET | /api/v1/auth/user-info | 获取用户信息 | ✅ |
| PUT | /api/v1/auth/user-info | 更新用户信息 | ✅ |
| PUT | /api/v1/auth/user-pwd | 修改密码 | ✅ |
| GET | /api/v1/auth/login-logs | 获取登录日志 | ✅ |
| GET | /api/v1/auth/sessions | 获取活跃会话 | ✅ |
| POST | /api/v1/auth/sessions/revoke | 撤销会话 | ✅ |

### 虚拟机管理模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/vm/list | 获取虚拟机列表 | ✅ |
| GET | /api/v1/vm/detail/:id | 获取虚拟机详情 | ✅ |
| POST | /api/v1/vm/create | 创建虚拟机 | ✅ |
| POST | /api/v1/vm/operate/:id | 操作虚拟机 | ✅ |
| PUT | /api/v1/vm/update-config/:id | 更新配置 | ✅ |
| GET | /api/v1/vm/snapshots/:id | 获取快照列表 | ✅ |
| POST | /api/v1/vm/snapshot/:id | 创建快照 | ✅ |
| POST | /api/v1/vm/snapshot/:id/restore | 恢复快照 | ✅ |
| DELETE | /api/v1/vm/snapshot/:id/:snapshot_id | 删除快照 | ✅ |
| GET | /api/v1/vm/logs/:id | 获取操作日志 | ✅ |
| GET | /api/v1/vm/metrics/:id | 获取监控指标 | ✅ |

### 目录挂载模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/mount/list | 获取挂载列表 | ✅ |
| POST | /api/v1/mount/add | 添加挂载 | ✅ |
| POST | /api/v1/mount/operate/:id | 操作挂载 | ✅ |
| GET | /api/v1/mount/openclaw-config | 获取 OpenClaw 配置 | ✅ |
| POST | /api/v1/mount/openclaw-config | 配置 OpenClaw 挂载 | ✅ |

### OpenClaw 管理模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/openclaw/status | 获取状态 | ✅ |
| POST | /api/v1/openclaw/deploy | 部署 OpenClaw | ✅ |
| POST | /api/v1/openclaw/operate/:id | 操作 OpenClaw | ✅ |
| GET | /api/v1/openclaw/log/:id | 获取日志 | ✅ |
| PUT | /api/v1/openclaw/config/:id | 更新配置 | ✅ |
| GET | /api/v1/openclaw/monitor | 获取监控数据 | ✅ |
| GET | /api/v1/openclaw/workspace | 获取工作空间 | ✅ |

### 远程访问模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/remote/status | 获取状态 | ✅ |
| PUT | /api/v1/remote/switch/:id | 切换开关 | ✅ |
| GET | /api/v1/remote/ip-whitelist | 获取白名单 | ✅ |
| POST | /api/v1/remote/ip-whitelist | 添加白名单 | ✅ |
| DELETE | /api/v1/remote/ip-whitelist/:id/:whitelist_id | 删除白名单 | ✅ |
| GET | /api/v1/remote/log/:id | 获取访问日志 | ✅ |

### 系统模块

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/system/setting | 获取设置 | ✅ |
| PUT | /api/v1/system/setting | 更新设置 | ✅ |
| GET | /api/v1/system/version | 获取版本 | ❌ |
| GET | /api/v1/system/dashboard | 获取仪表盘 | ✅ |
| GET | /api/v1/system/statistics | 获取统计 | ✅ |
| GET | /api/v1/system/search | 搜索虚拟机 | ✅ |
| POST | /api/v1/system/batch-vm | 批量操作 | ✅ |

## 请求示例

### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "Admin123",
    "confirm_password": "Admin123"
  }'
```

### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123"
  }'
```

### 创建虚拟机

```bash
curl -X POST http://localhost:8080/api/v1/vm/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "my-vm",
    "image": "ubuntu:22.04",
    "cpu": 2,
    "memory": 4,
    "disk": 50,
    "network_type": "nat"
  }'
```

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "参数错误",
  "data": null
}
```

### 分页响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t multipass-backend .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=multipass \
  --name multipass-backend \
  multipass-backend
```

### Systemd 服务

创建 `/etc/systemd/system/multipass-backend.service`:

```ini
[Unit]
Description=Multipass Backend API
After=network.target postgresql.service

[Service]
Type=simple
User=multipass
WorkingDirectory=/opt/multipass-backend
ExecStart=/opt/multipass-backend/multipass-backend
Restart=on-failure
Environment="SERVER_PORT=8080"
Environment="DB_HOST=localhost"
Environment="DB_USER=multipass"
Environment="DB_PASSWORD=your_password"
Environment="DB_NAME=multipass"

[Install]
WantedBy=multi-user.target
```

启动服务:

```bash
sudo systemctl daemon-reload
sudo systemctl enable multipass-backend
sudo systemctl start multipass-backend
```

## 开发

### 运行测试

```bash
go test ./... -v
```

### 代码格式化

```bash
go fmt ./...
```

### 代码检查

```bash
go vet ./...
```

## 许可证

MIT License

## 联系方式

- 项目地址：https://github.com/your-org/multipass-backend
- 问题反馈：https://github.com/your-org/multipass-backend/issues
