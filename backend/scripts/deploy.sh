#!/bin/bash

# Multipass Backend 部署脚本
# 使用方式：./scripts/deploy.sh [dev|prod]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
APP_NAME="multipass-backend"
APP_DIR="/opt/multipass-backend"
SYSTEMD_SERVICE="/etc/systemd/system/${APP_NAME}.service"
BACKUP_DIR="/opt/backups/multipass-backend"

# 检查是否以 root 运行
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}请使用 sudo 运行此脚本${NC}"
  exit 1
fi

# 检测部署环境
ENV="${1:-prod}"
echo -e "${GREEN}部署环境：${ENV}${NC}"

# 函数：打印信息
log_info() {
  echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
  echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
  echo -e "${RED}[ERROR]${NC} $1"
}

# 函数：检查依赖
check_dependencies() {
  log_info "检查依赖..."
  
  local deps=("go" "git")
  for dep in "${deps[@]}"; do
    if ! command -v "$dep" &> /dev/null; then
      log_error "$dep 未安装"
      exit 1
    fi
  done
  
  # 生产环境需要 Docker
  if [ "$ENV" = "prod" ]; then
    if ! command -v docker &> /dev/null; then
      log_warn "Docker 未安装，将使用二进制部署"
    fi
  fi
  
  log_info "依赖检查完成"
}

# 函数：创建目录
create_directories() {
  log_info "创建目录..."
  
  mkdir -p "$APP_DIR"
  mkdir -p "$BACKUP_DIR"
  
  log_info "目录创建完成"
}

# 函数：备份现有版本
backup_existing() {
  if [ -d "$APP_DIR" ] && [ "$(ls -A $APP_DIR)" ]; then
    log_info "备份现有版本..."
    
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_path="${BACKUP_DIR}/${timestamp}"
    
    cp -r "$APP_DIR" "$backup_path"
    
    log_info "备份完成：$backup_path"
  fi
}

# 函数：停止服务
stop_service() {
  log_info "停止服务..."
  
  if systemctl is-active --quiet "$APP_NAME"; then
    systemctl stop "$APP_NAME"
    log_info "服务已停止"
  else
    log_warn "服务未运行"
  fi
}

# 函数：编译应用
build_app() {
  log_info "编译应用..."
  
  cd /root/.openclaw/workspace/multipass-backend
  
  if [ "$ENV" = "prod" ]; then
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o "$APP_DIR/${APP_NAME}" cmd/server/main.go
  else
    go build -o "$APP_DIR/${APP_NAME}" cmd/server/main.go
  fi
  
  chmod +x "$APP_DIR/${APP_NAME}"
  
  log_info "编译完成"
}

# 函数：创建 systemd 服务
create_systemd_service() {
  log_info "创建 systemd 服务..."
  
  cat > "$SYSTEMD_SERVICE" <<EOF
[Unit]
Description=Multipass Backend API
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=${APP_DIR}
ExecStart=${APP_DIR}/${APP_NAME}
Restart=on-failure
RestartSec=5
Environment="SERVER_PORT=8080"
Environment="DB_HOST=localhost"
Environment="DB_USER=postgres"
Environment="DB_PASSWORD=postgres"
Environment="DB_NAME=multipass"
Environment="DB_SSLMODE=disable"
Environment="JWT_SECRET=your-secret-key-change-in-production"

# 安全限制
NoNewPrivileges=true
PrivateTmp=true

# 日志
StandardOutput=journal
StandardError=journal
SyslogIdentifier=${APP_NAME}

[Install]
WantedBy=multi-user.target
EOF

  systemctl daemon-reload
  
  log_info "systemd 服务创建完成"
}

# 函数：配置环境变量
setup_environment() {
  log_info "配置环境变量..."
  
  if [ ! -f "${APP_DIR}/.env" ]; then
    cp /root/.openclaw/workspace/multipass-backend/.env.example "${APP_DIR}/.env"
    log_warn "请编辑 ${APP_DIR}/.env 文件配置环境变量"
  fi
  
  log_info "环境配置完成"
}

# 函数：启动服务
start_service() {
  log_info "启动服务..."
  
  systemctl enable "$APP_NAME"
  systemctl start "$APP_NAME"
  
  sleep 2
  
  if systemctl is-active --quiet "$APP_NAME"; then
    log_info "服务启动成功"
  else
    log_error "服务启动失败"
    systemctl status "$APP_NAME"
    exit 1
  fi
}

# 函数：健康检查
health_check() {
  log_info "执行健康检查..."
  
  local max_attempts=10
  local attempt=1
  
  while [ $attempt -le $max_attempts ]; do
    if curl -s http://localhost:8080/health > /dev/null; then
      log_info "健康检查通过"
      return 0
    fi
    
    log_warn "健康检查失败，尝试 $attempt/$max_attempts"
    sleep 2
    attempt=$((attempt + 1))
  done
  
  log_error "健康检查失败"
  return 1
}

# 函数：Docker 部署
docker_deploy() {
  log_info "使用 Docker 部署..."
  
  cd /root/.openclaw/workspace/multipass-backend
  
  # 构建镜像
  docker build -t "${APP_NAME}:latest" .
  
  # 停止旧容器
  docker stop "${APP_NAME}" 2>/dev/null || true
  docker rm "${APP_NAME}" 2>/dev/null || true
  
  # 启动新容器
  docker run -d \
    --name "${APP_NAME}" \
    -p 8080:8080 \
    -e DB_HOST=postgres \
    -e DB_USER=postgres \
    -e DB_PASSWORD=postgres \
    -e DB_NAME=multipass \
    --restart unless-stopped \
    "${APP_NAME}:latest"
  
  log_info "Docker 部署完成"
}

# 主流程
main() {
  log_info "开始部署 Multipass Backend"
  
  check_dependencies
  create_directories
  backup_existing
  stop_service
  
  if [ "$ENV" = "docker" ]; then
    docker_deploy
  else
    build_app
    create_systemd_service
    setup_environment
    start_service
    health_check
  fi
  
  log_info "部署完成！"
  echo ""
  log_info "访问地址：http://localhost:8080"
  log_info "健康检查：http://localhost:8080/health"
  echo ""
  log_info "查看日志：journalctl -u ${APP_NAME} -f"
  log_info "停止服务：systemctl stop ${APP_NAME}"
  log_info "重启服务：systemctl restart ${APP_NAME}"
}

# 执行主流程
main
