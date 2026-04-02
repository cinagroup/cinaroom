#!/bin/bash
# CinaSeek 全平台编译脚本
# 输出到 build/output/ 目录
set -e

VERSION=${1:-"1.0.0"}
OUTPUT_DIR="build/output"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"
mkdir -p "$OUTPUT_DIR"

echo "========================================"
echo " CinaSeek Build All - v${VERSION}"
echo "========================================"
echo ""

# 后端 API 服务 (from backend/)
echo "[1/5] Building server (Linux amd64)..."
cd backend
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$VERSION" -o "../$OUTPUT_DIR/cinaseek-server-linux-amd64" ./cmd/server/
cd "$PROJECT_DIR"

echo "[2/5] Building server (Linux arm64)..."
cd backend
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$VERSION" -o "../$OUTPUT_DIR/cinaseek-server-linux-arm64" ./cmd/server/
cd "$PROJECT_DIR"

# WebSocket 中转服务 (from websocket/backend/)
echo "[3/5] Building relay (Linux amd64)..."
cd websocket/backend
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$VERSION" -o "../../$OUTPUT_DIR/cinaseek-relay-linux-amd64" .
cd "$PROJECT_DIR"

# 用户端客户端（三平台, from backend/）
echo "[4/5] Building client (Windows amd64)..."
cd backend
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$VERSION" -o "../$OUTPUT_DIR/cinaseek-client-windows-amd64.exe" ./cmd/client/ 2>/dev/null || echo "  Note: client windows-amd64 build skipped"

echo "[5/5] Building client (macOS arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$VERSION" -o "../$OUTPUT_DIR/cinaseek-client-darwin-arm64" ./cmd/client/ 2>/dev/null || echo "  Note: client darwin-arm64 build skipped"
cd "$PROJECT_DIR"

echo ""
echo "========================================"
echo " Build complete! Output:"
echo "========================================"
ls -lh "$OUTPUT_DIR/"
echo ""
