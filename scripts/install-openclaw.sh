#!/bin/bash
# =============================================================================
# CinaSeek - OpenClaw 一键部署脚本
# 在 CinaClaw 虚拟机内执行
# =============================================================================
set -euo pipefail

VERSION="1.0.0"
OPENCLAW_VERSION="${OPENCLAW_VERSION:-latest}"
WORKSPACE="/root/.openclaw/workspace"
SERVICE_FILE="/etc/systemd/system/openclaw.service"
LOG_PREFIX="[CinaSeek/OpenClaw]"

# ── 日志函数 ──────────────────────────────────────────────────────────────────
log_info()    { echo "$LOG_PREFIX ✅ $*"; }
log_warn()    { echo "$LOG_PREFIX ⚠️  $*" >&2; }
log_error()   { echo "$LOG_PREFIX ❌ $*" >&2; }
log_step()    { echo "$LOG_PREFIX 🔄 $*"; }
log_section() { echo ""; echo "$LOG_PREFIX ═══════════════════════════════════════"; echo "$LOG_PREFIX $*"; echo "$LOG_PREFIX ═══════════════════════════════════════"; }

# ── 错误处理 ──────────────────────────────────────────────────────────────────
cleanup() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        log_error "部署失败 (exit code: $exit_code)"
        log_error "请检查上方日志排查问题"
    fi
}
trap cleanup EXIT

# ── 前置检查 ──────────────────────────────────────────────────────────────────
check_prerequisites() {
    log_section "步骤 1/8: 前置环境检查"

    # 检查 root 权限
    if [ "$(id -u)" -ne 0 ]; then
        log_error "此脚本需要 root 权限执行"
        exit 1
    fi
    log_info "Root 权限: OK"

    # 检查操作系统
    if [ -f /etc/os-release ]; then
        # shellcheck disable=SC1091
        source /etc/os-release
        log_info "操作系统: $NAME $VERSION_ID"
        if [[ "$ID" != "ubuntu" ]] && [[ "$ID" != "debian" ]]; then
            log_warn "未经测试的操作系统: $ID，继续执行但可能出现问题"
        fi
    else
        log_warn "无法检测操作系统版本"
    fi

    # 检查架构
    ARCH=$(uname -m)
    log_info "系统架构: $ARCH"

    # 检查内存（建议 >= 4G）
    TOTAL_MEM=$(grep MemTotal /proc/meminfo | awk '{print int($2/1024)}')
    if [ "$TOTAL_MEM" -lt 2048 ]; then
        log_error "内存不足: ${TOTAL_MEM}MB（最低要求 2048MB）"
        exit 1
    fi
    log_info "内存: ${TOTAL_MEM}MB"

    # 检查磁盘空间（建议 >= 20G）
    AVAILABLE_DISK=$(df -BG / | awk 'NR==2 {print $4}' | tr -d 'G')
    if [ "${AVAILABLE_DISK:-0}" -lt 10 ]; then
        log_error "磁盘空间不足: ${AVAILABLE_DISK}GB（最低要求 10GB）"
        exit 1
    fi
    log_info "可用磁盘: ${AVAILABLE_DISK}GB"
}

# ── 安装系统依赖 ──────────────────────────────────────────────────────────────
install_system_deps() {
    log_section "步骤 2/8: 安装系统依赖"

    log_step "更新软件源..."
    apt-get update -qq || {
        log_error "apt-get update 失败"
        exit 1
    }

    log_step "安装基础工具..."
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq \
        curl wget git build-essential python3 ca-certificates gnupg \
        >/dev/null 2>&1 || {
        log_error "系统依赖安装失败"
        exit 1
    }
    log_info "系统依赖安装完成"
}

# ── 安装 Node.js ──────────────────────────────────────────────────────────────
install_nodejs() {
    log_section "步骤 3/8: 安装 Node.js 22.x"

    # 检查是否已安装
    if command -v node &>/dev/null; then
        NODE_VER=$(node -v)
        if [[ "$NODE_VER" == v22.* ]]; then
            log_info "Node.js 已安装: $NODE_VER，跳过"
            return 0
        else
            log_warn "现有 Node.js 版本 $NODE_VER 不符合要求，将安装 v22.x"
        fi
    fi

    log_step "添加 NodeSource 仓库..."
    curl -fsSL https://deb.nodesource.com/setup_22.x | bash - >/dev/null 2>&1 || {
        log_error "NodeSource 仓库添加失败"
        exit 1
    }

    log_step "安装 Node.js..."
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq nodejs >/dev/null 2>&1 || {
        log_error "Node.js 安装失败"
        exit 1
    }

    NODE_VER=$(node -v)
    NPM_VER=$(npm -v)
    log_info "Node.js 安装完成: $NODE_VER (npm $NPM_VER)"
}

# ── 安装 pnpm ─────────────────────────────────────────────────────────────────
install_pnpm() {
    log_section "步骤 4/8: 安装 pnpm"

    if command -v pnpm &>/dev/null; then
        PNPM_VER=$(pnpm -v)
        log_info "pnpm 已安装: $PNPM_VER，跳过"
        return 0
    fi

    log_step "通过 corepack 启用 pnpm..."
    corepack enable >/dev/null 2>&1 || {
        log_warn "corepack 启用失败，尝试 npm 安装..."
        npm install -g pnpm >/dev/null 2>&1 || {
            log_error "pnpm 安装失败"
            exit 1
        }
    }

    PNPM_VER=$(pnpm -v)
    log_info "pnpm 安装完成: $PNPM_VER"
}

# ── 安装 OpenClaw ─────────────────────────────────────────────────────────────
install_openclaw() {
    log_section "步骤 5/8: 安装 OpenClaw"

    if command -v openclaw &>/dev/null; then
        OC_VER=$(openclaw --version 2>/dev/null || echo "unknown")
        log_info "OpenClaw 已安装: $OC_VER"
        if [ "$OPENCLAW_VERSION" = "latest" ]; then
            log_step "更新到最新版本..."
            npm install -g "openclaw@latest" >/dev/null 2>&1 || {
                log_warn "OpenClaw 更新失败，保持现有版本"
            }
        fi
        return 0
    fi

    log_step "安装 OpenClaw ($OPENCLAW_VERSION)..."
    if [ "$OPENCLAW_VERSION" = "latest" ]; then
        npm install -g openclaw >/dev/null 2>&1 || {
            log_error "OpenClaw 安装失败"
            exit 1
        }
    else
        npm install -g "openclaw@$OPENCLAW_VERSION" >/dev/null 2>&1 || {
            log_error "OpenClaw $OPENCLAW_VERSION 安装失败"
            exit 1
        }
    fi

    OC_VER=$(openclaw --version 2>/dev/null || echo "unknown")
    log_info "OpenClaw 安装完成: $OC_VER"
}

# ── 初始化工作空间 ─────────────────────────────────────────────────────────────
init_workspace() {
    log_section "步骤 6/8: 初始化工作空间"

    log_step "创建工作空间目录: $WORKSPACE"
    mkdir -p "$WORKSPACE"
    log_info "工作空间目录已创建"

    # 创建基础文件
    if [ ! -f "$WORKSPACE/AGENTS.md" ]; then
        cat > "$WORKSPACE/AGENTS.md" <<'AGENTSEOF'
# AGENTS.md - OpenClaw Workspace

This is your OpenClaw AI agent workspace.
AGENTSEOF
        log_info "AGENTS.md 已创建"
    fi

    if [ ! -f "$WORKSPACE/SOUL.md" ]; then
        cat > "$WORKSPACE/SOUL.md" <<'SOULEOF'
# SOUL.md

You are an AI assistant running in a CinaSeek-managed OpenClaw environment.
SOULEOF
        log_info "SOUL.md 已创建"
    fi

    log_info "工作空间初始化完成"
}

# ── 配置 systemd 服务 ─────────────────────────────────────────────────────────
setup_systemd() {
    log_section "步骤 7/8: 配置 systemd 服务"

    # 确定 openclaw 二进制路径
    OPENCLAW_BIN=$(command -v openclaw)
    if [ -z "$OPENCLAW_BIN" ]; then
        log_error "找不到 openclaw 二进制文件"
        exit 1
    fi
    log_info "OpenClaw 路径: $OPENCLAW_BIN"

    # 如果不在标准路径，创建软链接
    if [[ "$OPENCLAW_BIN" != /usr/local/bin/openclaw ]]; then
        ln -sf "$OPENCLAW_BIN" /usr/local/bin/openclaw
        OPENCLAW_BIN="/usr/local/bin/openclaw"
    fi

    log_step "写入 systemd 服务文件: $SERVICE_FILE"
    cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=OpenClaw AI Agent Runtime
After=network.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=${OPENCLAW_BIN} gateway start
ExecReload=/bin/kill -HUP \$MAINPID
Restart=always
RestartSec=5
TimeoutStopSec=30

Environment=NODE_ENV=production
Environment=HOME=/root

WorkingDirectory=${WORKSPACE}

StandardOutput=journal
StandardError=journal
SyslogIdentifier=openclaw

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable openclaw.service >/dev/null 2>&1 || {
        log_warn "systemctl enable 失败（在容器环境中可忽略）"
    }
    log_info "systemd 服务配置完成"
}

# ── 启动并验证 ─────────────────────────────────────────────────────────────────
start_and_verify() {
    log_section "步骤 8/8: 启动 OpenClaw 并验证"

    log_step "启动 OpenClaw 服务..."
    systemctl start openclaw.service 2>/dev/null || {
        log_warn "systemctl start 失败，尝试直接启动..."
        # 容器环境回退
        nohup openclaw gateway start >/var/log/openclaw.log 2>&1 &
        sleep 3
    }

    # 等待服务启动
    log_step "等待服务就绪..."
    local retries=0
    local max_retries=30
    while [ $retries -lt $max_retries ]; do
        if curl -sf http://127.0.0.1:3271/health >/dev/null 2>&1; then
            log_info "OpenClaw 服务已就绪！"
            break
        fi
        retries=$((retries + 1))
        sleep 2
    done

    if [ $retries -ge $max_retries ]; then
        log_warn "服务健康检查超时，可能需要手动检查"
    fi

    # 验证安装
    log_step "验证安装..."
    local all_ok=true

    if command -v node &>/dev/null; then
        log_info "Node.js: $(node -v)"
    else
        log_error "Node.js: 未找到"
        all_ok=false
    fi

    if command -v pnpm &>/dev/null; then
        log_info "pnpm: $(pnpm -v)"
    else
        log_error "pnpm: 未找到"
        all_ok=false
    fi

    if command -v openclaw &>/dev/null; then
        log_info "OpenClaw: $(openclaw --version 2>/dev/null || echo 'installed')"
    else
        log_error "OpenClaw: 未找到"
        all_ok=false
    fi

    if [ -d "$WORKSPACE" ]; then
        log_info "工作空间: $WORKSPACE ✓"
    else
        log_error "工作空间: 不存在"
        all_ok=false
    fi

    # 最终报告
    echo ""
    log_section "部署完成"
    if [ "$all_ok" = true ]; then
        log_info "所有组件安装成功！"
        echo ""
        echo "  管理命令:"
        echo "    systemctl status openclaw   # 查看状态"
        echo "    systemctl restart openclaw  # 重启服务"
        echo "    systemctl stop openclaw     # 停止服务"
        echo "    journalctl -u openclaw -f   # 查看日志"
        echo ""
        echo "  工作空间: $WORKSPACE"
        echo "  Gateway:  http://0.0.0.0:3271"
    else
        log_warn "部分组件安装异常，请检查上方日志"
        exit 1
    fi
}

# ── 主流程 ─────────────────────────────────────────────────────────────────────
main() {
    echo ""
    echo "🐾 CinaSeek - OpenClaw 一键部署 v$VERSION"
    echo "   时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "   目标版本: $OPENCLAW_VERSION"
    echo ""

    check_prerequisites
    install_system_deps
    install_nodejs
    install_pnpm
    install_openclaw
    init_workspace
    setup_systemd
    start_and_verify
}

main "$@"
