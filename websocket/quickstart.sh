#!/bin/bash
# 快速启动脚本

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_dependencies() {
    log_info "Checking dependencies..."
    
    local missing=()
    
    if ! command -v docker &> /dev/null; then
        missing+=("docker")
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        missing+=("docker-compose")
    fi
    
    if ! command -v curl &> /dev/null; then
        missing+=("curl")
    fi
    
    if [ ${#missing[@]} -ne 0 ]; then
        log_error "Missing dependencies: ${missing[*]}"
        log_info "Install with:"
        log_info "  Ubuntu: sudo apt-get install -y ${missing[*]}"
        log_info "  macOS: brew install ${missing[*]}"
        exit 1
    fi
    
    log_success "All dependencies installed"
}

check_env() {
    log_info "Checking environment configuration..."
    
    if [ ! -f deploy/.env ]; then
        log_warn "deploy/.env not found, creating from example..."
        cp deploy/.env.example deploy/.env
        log_warn "Please edit deploy/.env with your configuration"
        log_warn "At minimum, set CLOUDFLARE_TUNNEL_ID and CLOUDFLARE_HOSTNAME"
        
        read -p "Continue without configuration? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

start_services() {
    log_info "Starting services..."
    
    cd deploy
    
    if docker-compose ps | grep -q "Up"; then
        log_warn "Services already running, restarting..."
        docker-compose restart
    else
        docker-compose up -d
    fi
    
    cd ..
    
    log_success "Services started"
}

wait_health() {
    log_info "Waiting for services to be healthy..."
    
    local max_attempts=30
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            log_success "Terminal service is healthy"
            break
        fi
        
        attempt=$((attempt + 1))
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        log_error "Terminal service failed to become healthy"
        return 1
    fi
    
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:8081/health > /dev/null 2>&1; then
            log_success "Relay service is healthy"
            break
        fi
        
        attempt=$((attempt + 1))
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        log_error "Relay service failed to become healthy"
        return 1
    fi
}

show_status() {
    log_info "Service status:"
    echo ""
    
    cd deploy
    docker-compose ps
    cd ..
    
    echo ""
    log_info "Ports:"
    echo "  Terminal Service: http://localhost:8080"
    echo "  Relay Service:    http://localhost:8081"
    echo "  Nginx (if enabled): http://localhost:80"
    echo ""
    log_info "Logs:"
    echo "  View all logs: docker-compose logs -f"
    echo "  Terminal only: docker-compose logs -f websocket-terminal"
    echo "  Relay only:    docker-compose logs -f cloud-relay"
    echo ""
    log_info "Testing:"
    echo "  Health check: curl http://localhost:8080/health"
    echo "  WebSocket test: websocat ws://localhost:8080/api/v1/ws/terminal?token=test"
    echo "  Run tests: cd tests && ./integration_test.sh"
}

main() {
    echo "========================================"
    echo "  MultiPass WebSocket Quick Start"
    echo "========================================"
    echo ""
    
    check_dependencies
    check_env
    start_services
    wait_health
    show_status
    
    echo ""
    log_success "Setup completed!"
    echo ""
    log_info "Next steps:"
    echo "  1. Configure Cloudflare Tunnel (see docs/DEPLOYMENT.md)"
    echo "  2. Generate auth tokens for production use"
    echo "  3. Run integration tests: cd tests && ./integration_test.sh"
    echo "  4. Monitor logs: docker-compose logs -f"
}

# 处理命令行参数
case "${1:-start}" in
    start)
        main
        ;;
    stop)
        log_info "Stopping services..."
        cd deploy && docker-compose down
        log_success "Services stopped"
        ;;
    restart)
        log_info "Restarting services..."
        cd deploy && docker-compose restart
        log_success "Services restarted"
        ;;
    logs)
        cd deploy && docker-compose logs -f
        ;;
    status)
        cd deploy && docker-compose ps
        ;;
    test)
        cd tests && ./integration_test.sh
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|logs|status|test}"
        exit 1
        ;;
esac
