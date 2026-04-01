.PHONY: help build run test clean docker-build docker-run migrate dev prod

# 生成 protobuf Go 代码
proto:
	@echo "Generating protobuf Go code..."
	@mkdir -p proto/cinaclaw
	protoc \
		--go_out=proto/cinaclaw --go_opt=paths=source_relative \
		--go-grpc_out=proto/cinaclaw --go-grpc_opt=paths=source_relative \
		-I ../cinaclaw/src/rpc/ \
		-I /usr/include/ \
		../cinaclaw/src/rpc/cinaclaw.proto
	@echo "Proto generation complete"

# 默认目标
help:
	@echo "Multipass Backend - Makefile Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make build        - 编译应用"
	@echo "  make run          - 运行应用（开发模式）"
	@echo "  make test         - 运行测试"
	@echo "  make clean        - 清理构建文件"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo "  make docker-run   - 运行 Docker Compose"
	@echo "  make migrate      - 运行数据库迁移"
	@echo "  make dev          - 开发环境（热重载）"
	@echo "  make prod         - 生产环境编译"
	@echo ""

# 编译应用
build:
	@echo "Building application..."
	go build -o bin/cinaseek-backend cmd/server/main.go
	@echo "Build complete: bin/cinaseek-backend"

# 运行应用
run:
	@echo "Starting application..."
	go run cmd/server/main.go

# 运行测试
test:
	@echo "Running tests..."
	go test -v ./... -cover

# 清理构建文件
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f cinaseek-backend
	@echo "Clean complete"

# 构建 Docker 镜像
docker-build:
	@echo "Building Docker image..."
	docker build -t cinaseek-backend:latest .
	@echo "Docker image built"

# 运行 Docker Compose
docker-run:
	@echo "Starting Docker Compose..."
	docker-compose up -d
	@echo "Services started"

# 停止 Docker Compose
docker-stop:
	@echo "Stopping Docker Compose..."
	docker-compose down
	@echo "Services stopped"

# 查看 Docker 日志
docker-logs:
	docker-compose logs -f backend

# 重启 Docker 服务
docker-restart:
	docker-compose restart backend

# 开发环境（需要 air 支持热重载）
dev:
	@echo "Starting development mode with hot reload..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air -c .air.toml

# 生产环境编译
prod:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o bin/cinaseek-backend cmd/server/main.go
	@echo "Production build complete"

# 代码格式化
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete"

# 代码检查
vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "Vet complete"

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

# 生成 API 文档（需要 swag）
docs:
	@echo "Generating API documentation..."
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/server/main.go -o ./api/docs
	@echo "Documentation generated"

# 数据库迁移（如果有迁移工具）
migrate:
	@echo "Running database migrations..."
	@echo "Migrations are auto-run on startup"

# 创建数据库
db-create:
	@echo "Creating database..."
	psql -U postgres -c "CREATE DATABASE cinaseek;" || true
	@echo "Database created"

# 备份数据库
db-backup:
	@echo "Backing up database..."
	pg_dump -U postgres cinaseek > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Backup complete"

# 恢复数据库
db-restore:
	@echo "Restoring database from backup..."
	@read -p "Enter backup file: " file; psql -U postgres cinaseek < $$file
	@echo "Restore complete"

# 压力测试（需要 wrk）
stress-test:
	@echo "Running stress test..."
	@which wrk > /dev/null || (echo "Please install wrk: brew install wrk" && exit 1)
	wrk -t12 -c400 -d30s http://localhost:8080/api/v1/health
	@echo "Stress test complete"

# 性能分析
profile:
	@echo "Starting profiling..."
	go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
	@echo "Profiling complete"

# 生成覆盖率报告
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# 安全扫描（需要 gosec）
security:
	@echo "Running security scan..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securego/gosec/v2/cmd/gosec@latest)
	gosec ./...
	@echo "Security scan complete"

# 初始化项目
init: deps db-create build
	@echo "Project initialized"

# 完整 CI 流程
ci: fmt vet test security build
	@echo "CI pipeline complete"
