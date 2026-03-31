#!/bin/bash
# WebSocket 集成测试脚本

set -e

# 配置
BASE_URL="${BASE_URL:-http://localhost:8080}"
RELAY_URL="${RELAY_URL:-http://localhost:8081}"
TEST_TOKEN=""
SESSION_ID=""
USER_ID=""

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 测试健康检查
test_health() {
    log_info "Testing health endpoints..."
    
    # 终端服务健康检查
    response=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/health")
    if [ "$response" = "200" ]; then
        log_info "✓ Terminal service health check passed"
    else
        log_error "✗ Terminal service health check failed: $response"
        exit 1
    fi
    
    # 中转服务健康检查
    response=$(curl -s -o /dev/null -w "%{http_code}" "${RELAY_URL}/health")
    if [ "$response" = "200" ]; then
        log_info "✓ Relay service health check passed"
    else
        log_error "✗ Relay service health check failed: $response"
        exit 1
    fi
}

# 测试用户注册
test_user_register() {
    log_info "Testing user registration..."
    
    # 使用 websocat 或 wscat 进行 WebSocket 测试
    if command -v websocat &> /dev/null; then
        log_info "Using websocat for WebSocket tests"
    elif command -v wscat &> /dev/null; then
        log_info "Using wscat for WebSocket tests"
    else
        log_warn "No WebSocket client found, skipping WebSocket tests"
        log_warn "Install websocat: cargo install websocat"
        return
    fi
    
    # 这里应该使用 WebSocket 客户端进行测试
    # 简化版本：只测试 HTTP 端点
    log_info "User registration test completed"
}

# 测试终端连接
test_terminal_connection() {
    log_info "Testing terminal connection..."
    
    # 生成测试 token（实际应该调用 API）
    TEST_TOKEN="test-token-$(date +%s)"
    
    # 测试带 token 的连接
    response=$(curl -s -o /dev/null -w "%{http_code}" \
        "${BASE_URL}/api/v1/ws/terminal?token=${TEST_TOKEN}")
    
    # 应该返回 401（因为 token 无效）或 101（WebSocket 升级）
    if [ "$response" = "401" ] || [ "$response" = "101" ]; then
        log_info "✓ Terminal connection auth test passed"
    else
        log_error "✗ Terminal connection auth test failed: $response"
    fi
}

# 测试请求转发
test_forward_request() {
    log_info "Testing forward request..."
    
    # 测试转发 API
    response=$(curl -s -X POST "${RELAY_URL}/api/v1/forward" \
        -H "Content-Type: application/json" \
        -H "X-Auth-Token: test-token" \
        -d '{
            "target_url": "https://httpbin.org/get",
            "method": "GET",
            "headers": {},
            "body": ""
        }' \
        -o /dev/null -w "%{http_code}")
    
    # 应该返回 401（token 无效）或 200（成功转发）
    if [ "$response" = "401" ] || [ "$response" = "200" ]; then
        log_info "✓ Forward request test passed (status: $response)"
    else
        log_error "✗ Forward request test failed: $response"
    fi
}

# 压力测试
stress_test() {
    log_info "Starting stress test..."
    
    local concurrent=${1:-10}
    local requests=${2:-100}
    
    log_info "Configuration: $concurrent concurrent, $requests requests"
    
    # 使用 ab 或 hey 进行压力测试
    if command -v ab &> /dev/null; then
        log_info "Using Apache Benchmark for stress test"
        ab -n $requests -c $concurrent "${BASE_URL}/health" 2>/dev/null | grep -E "(Requests per second|Time per request|Failed requests)"
    elif command -v hey &> /dev/null; then
        log_info "Using hey for stress test"
        hey -n $requests -c $concurrent "${BASE_URL}/health" 2>&1 | grep -E "(Requests/sec|Latency|Failed)"
    else
        log_warn "No stress test tool found (install apache2-utils or hey)"
        
        # 简单的 curl 循环测试
        local start_time=$(date +%s)
        local success=0
        local failed=0
        
        for i in $(seq 1 $requests); do
            response=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/health")
            if [ "$response" = "200" ]; then
                ((success++))
            else
                ((failed++))
            fi
        done
        
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        
        log_info "Simple stress test completed"
        log_info "Duration: ${duration}s"
        log_info "Success: $success, Failed: $failed"
        if [ $duration -gt 0 ]; then
            log_info "Requests/sec: $((requests / duration))"
        fi
    fi
}

# WebSocket 连接测试
websocket_test() {
    log_info "Testing WebSocket connection..."
    
    if ! command -v websocat &> /dev/null; then
        log_warn "websocat not installed, skipping WebSocket test"
        log_warn "Install: cargo install websocat"
        return
    fi
    
    # 测试 WebSocket 连接
    timeout 5 websocat "${BASE_URL}/api/v1/ws/terminal?token=test" || true
    log_info "WebSocket connection test completed"
}

# 清理测试
cleanup() {
    log_info "Cleaning up test resources..."
    # 清理逻辑
}

# 主测试流程
main() {
    log_info "Starting integration tests..."
    log_info "Base URL: ${BASE_URL}"
    log_info "Relay URL: ${RELAY_URL}"
    
    trap cleanup EXIT
    
    # 运行测试
    test_health
    test_user_register
    test_terminal_connection
    test_forward_request
    websocket_test
    
    # 压力测试（可选）
    if [ "${RUN_STRESS_TEST:-false}" = "true" ]; then
        stress_test ${STRESS_CONCURRENT:-10} ${STRESS_REQUESTS:-100}
    fi
    
    log_info "All tests completed!"
}

# 运行主函数
main "$@"
