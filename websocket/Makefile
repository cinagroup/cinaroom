# Makefile for MultiPass WebSocket

.PHONY: all build test clean deploy help

# 变量
BINARY_NAME=terminal
RELAY_NAME=cloud-relay
GO=go
DOCKER=docker
DOCKER_COMPOSE=docker-compose

# 默认目标
all: build

# 帮助信息
help:
	@echo "MultiPass WebSocket Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build       - 编译后端服务"
	@echo "  test        - 运行单元测试"
	@echo "  bench       - 运行基准测试"
	@echo "  clean       - 清理编译文件"
	@echo "  run         - 运行本地服务"
	@echo "  deploy      - Docker 部署"
	@echo "  logs        - 查看 Docker 日志"
	@echo "  stop        - 停止 Docker 服务"
	@echo "  test-integration - 运行集成测试"
	@echo "  help        - 显示帮助信息"

# 编译后端服务
build:
	@echo "Building backend services..."
	cd backend && $(GO) mod download
	cd backend && CGO_ENABLED=0 GOOS=linux $(GO) build -o terminal .
	cd backend && CGO_ENABLED=0 GOOS=linux $(GO) build -o cloud-relay .
	@echo "Build completed!"

# 运行单元测试
test:
	@echo "Running unit tests..."
	cd backend && $(GO) test -v ./...

# 运行基准测试
bench:
	@echo "Running benchmark tests..."
	cd backend && $(GO) test -bench=. -benchmem ./...

# 清理编译文件
clean:
	@echo "Cleaning..."
	rm -f backend/terminal backend/cloud-relay
	rm -rf deploy/logs/*
	@echo "Clean completed!"

# 运行本地服务
run:
	@echo "Starting local services..."
	cd backend && $(GO) run . &
	sleep 2
	cd backend && $(GO) run cloud_relay.go
	@echo "Services started!"

# Docker 部署
deploy:
	@echo "Deploying with Docker..."
	cd deploy && $(DOCKER_COMPOSE) up -d --build
	@echo "Deployment completed!"

# 查看 Docker 日志
logs:
	cd deploy && $(DOCKER_COMPOSE) logs -f

# 停止 Docker 服务
stop:
	@echo "Stopping services..."
	cd deploy && $(DOCKER_COMPOSE) down
	@echo "Services stopped!"

# 运行集成测试
test-integration:
	@echo "Running integration tests..."
	cd tests && bash integration_test.sh
	@echo "Integration tests completed!"

# Python WebSocket 测试
test-ws:
	@echo "Running WebSocket tests..."
	cd tests && python3 websocket_test.py --url ws://localhost:8080

# 生成测试覆盖率报告
coverage:
	@echo "Generating coverage report..."
	cd backend && $(GO) test -coverprofile=coverage.out ./...
	cd backend && $(GO) tool cover -html=coverage.out -o ../docs/coverage.html
	@echo "Coverage report generated: docs/coverage.html"

# 压力测试
stress-test:
	@echo "Running stress test..."
	RUN_STRESS_TEST=true STRESS_CONCURRENT=50 STRESS_REQUESTS=1000 \
		cd tests && bash integration_test.sh

# 检查代码格式
fmt:
	@echo "Formatting code..."
	cd backend && $(GO) fmt ./...
	cd backend && goimports -w .

# 代码检查
lint:
	@echo "Running linter..."
	cd backend && golangci-lint run

# 安全扫描
security:
	@echo "Running security scan..."
	cd backend && govulncheck ./...

# 构建 Docker 镜像
docker-build:
	@echo "Building Docker images..."
	cd deploy && $(DOCKER) build -f Dockerfile.terminal -t multipass/terminal:latest .
	cd deploy && $(DOCKER) build -f Dockerfile.relay -t multipass/relay:latest .
	@echo "Docker images built!"

# 推送 Docker 镜像
docker-push:
	@echo "Pushing Docker images..."
	$(DOCKER) push multipass/terminal:latest
	$(DOCKER) push multipass/relay:latest
	@echo "Docker images pushed!"

# 查看服务状态
status:
	cd deploy && $(DOCKER_COMPOSE) ps

# 重启服务
restart:
	@echo "Restarting services..."
	cd deploy && $(DOCKER_COMPOSE) restart
	@echo "Services restarted!"

# 查看资源使用
stats:
	cd deploy && $(DOCKER_COMPOSE) stats

# 进入容器
shell-terminal:
	cd deploy && $(DOCKER_COMPOSE) exec websocket-terminal /bin/sh

shell-relay:
	cd deploy && $(DOCKER_COMPOSE) exec cloud-relay /bin/sh

# 初始化部署
init:
	@echo "Initializing deployment..."
	cp deploy/.env.example deploy/.env
	@echo "Please edit deploy/.env with your configuration"
	@echo "Then run: make deploy"
