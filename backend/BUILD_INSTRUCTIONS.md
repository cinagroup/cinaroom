# 构建说明

## 环境要求

### 必需
- **Go**: 1.23+ (由于依赖包要求)
- **PostgreSQL**: 12+
- **Git**

### 可选
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **Make**: GNU Make 4.0+
- **Nginx**: 1.20+ (生产环境)

## Go 1.23 安装

如果系统 Go 版本低于 1.23，需要升级：

### 方式 1: 官方二进制安装

```bash
# 下载 Go 1.23
cd /tmp
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz

# 删除旧版本（如果有）
sudo rm -rf /usr/local/go

# 解压
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz

# 添加到 PATH
export PATH=$PATH:/usr/local/go/bin

# 验证
go version
# 应该显示：go version go1.23.0 linux/amd64
```

### 方式 2: 使用 GVM (Go Version Manager)

```bash
# 安装 gvm
bash < <(curl -s https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)

# 加载 gvm
source ~/.gvm/scripts/gvm

# 安装 Go 1.23
gvm install go1.23

# 使用 Go 1.23
gvm use go1.23

# 设为默认
gvm use go1.23 --default
```

### 方式 3: 使用 Docker（推荐用于开发）

```bash
# 使用官方 Go 镜像
docker run --rm -v $(pwd):/app -w /app golang:1.23 go build -o bin/multipass-backend cmd/server/main.go
```

## 构建步骤

### 1. 下载依赖

```bash
cd /root/.openclaw/workspace/multipass-backend
go mod download
```

### 2. 编译

```bash
# 开发版本
go build -o bin/multipass-backend cmd/server/main.go

# 生产版本（优化体积）
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o bin/multipass-backend cmd/server/main.go
```

### 3. 验证

```bash
# 检查二进制文件
ls -lh bin/multipass-backend

# 运行帮助（如果实现了）
./bin/multipass-backend --help
```

## 使用 Make

```bash
# 编译
make build

# 生产编译
make prod

# 运行
make run

# 开发模式（热重载）
make dev

# 测试
make test

# Docker 构建
make docker-build

# Docker 运行
make docker-run
```

## Docker 构建

### 方式 1: 直接构建

```bash
docker build -t multipass-backend:latest .
```

### 方式 2: 使用 Docker Compose

```bash
# 构建并启动所有服务
docker-compose up -d --build

# 只构建 backend
docker-compose build backend

# 查看日志
docker-compose logs -f backend
```

## 常见问题

### 问题 1: Go 版本过低

**错误**: `toolchain upgrade needed to resolve github.com/rogpeppe/go-internal/fmtsort`

**解决**: 升级到 Go 1.23+（见上方安装说明）

### 问题 2: 依赖下载失败

**错误**: `package not found` 或 `timeout`

**解决**:
```bash
# 清理缓存
go clean -modcache

# 使用国内镜像
export GOPROXY=https://goproxy.cn,direct

# 重新下载
go mod download
```

### 问题 3: 编译时内存不足

**错误**: `out of memory`

**解决**:
```bash
# 限制编译使用的内存
GOMEMLIMIT=512MiB go build -o bin/multipass-backend cmd/server/main.go

# 或添加 swap 空间
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### 问题 4: PostgreSQL 驱动编译失败

**错误**: `gcc failed` 或 `C compiler required`

**解决**:
```bash
# 安装编译工具
apt-get install build-essential

# 或使用 CGO_ENABLED=0
CGO_ENABLED=0 go build -o bin/multipass-backend cmd/server/main.go
```

## 交叉编译

### 编译 Linux AMD64

```bash
GOOS=linux GOARCH=amd64 go build -o bin/multipass-backend-linux-amd64 cmd/server/main.go
```

### 编译 macOS AMD64

```bash
GOOS=darwin GOARCH=amd64 go build -o bin/multipass-backend-darwin-amd64 cmd/server/main.go
```

### 编译 Windows AMD64

```bash
GOOS=windows GOARCH=amd64 go build -o bin/multipass-backend-windows-amd64.exe cmd/server/main.go
```

## 性能优化

### 1. 使用 UPX 压缩

```bash
# 安装 upx
apt-get install upx

# 压缩二进制文件
upx --best bin/multipass-backend

# 可减小 50-70% 体积
```

### 2. 去除调试信息

```bash
go build -ldflags="-s -w" -o bin/multipass-backend cmd/server/main.go
```

### 3. 静态编译

```bash
CGO_ENABLED=0 go build -a -ldflags="-s -w" -o bin/multipass-backend cmd/server/main.go
```

## 验证构建

### 1. 检查文件

```bash
# 查看文件大小
ls -lh bin/multipass-backend

# 查看文件类型
file bin/multipass-backend

# 应该显示：ELF 64-bit LSB executable, x86-64
```

### 2. 运行测试

```bash
# 运行单元测试
go test ./... -v

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 3. 健康检查

```bash
# 启动服务
./bin/multipass-backend &

# 等待几秒后检查
curl http://localhost:8080/health

# 应该返回：{"code":0,"message":"healthy",...}
```

## 部署检查清单

- [ ] Go 版本 >= 1.23
- [ ] 依赖已下载（go.sum 存在）
- [ ] 编译成功（bin/multipass-backend 存在）
- [ ] 二进制文件可执行
- [ ] 单元测试通过
- [ ] 健康检查通过
- [ ] 环境变量已配置
- [ ] 数据库已创建
- [ ] 防火墙规则已配置
- [ ] 日志目录已创建
- [ ] 备份策略已配置

## 下一步

构建成功后，参考以下文档：

1. **QUICKSTART.md** - 快速开始指南
2. **README.md** - 项目说明
3. **api/API.md** - API 文档
4. **scripts/deploy.sh** - 部署脚本

---

**最后更新**: 2026-04-01
**Go 版本要求**: 1.23+
