#!/bin/bash
# =============================================================================
# CinaSeek - 工作空间数据持久化挂载脚本
# 将宿主机目录挂载到虚拟机的 OpenClaw 工作空间
# 支持双向同步（rsync）
# =============================================================================
set -euo pipefail

VERSION="1.0.0"
LOG_PREFIX="[CinaSeek/Mount]"

# ── 默认配置 ──────────────────────────────────────────────────────────────────
HOST_PATH="${HOST_PATH:-}"          # 宿主机目录路径（必填）
VM_NAME="${VM_NAME:-}"              # 虚拟机名称（必填）
VM_PATH="${VM_PATH:-/root/.openclaw/workspace}"  # 虚拟机内挂载路径
MODE="${MODE:-sync}"                # sync | mount
DIRECTION="${DIRECTION:-bidirectional}"  # push | pull | bidirectional
EXCLUDE_FILE="${EXCLUDE_FILE:-}"    # rsync 排除文件路径
DRY_RUN="${DRY_RUN:-false}"         # 试运行模式

# ── 日志函数 ──────────────────────────────────────────────────────────────────
log_info()    { echo "$LOG_PREFIX ✅ $*"; }
log_warn()    { echo "$LOG_PREFIX ⚠️  $*" >&2; }
log_error()   { echo "$LOG_PREFIX ❌ $*" >&2; }
log_step()    { echo "$LOG_PREFIX 🔄 $*"; }

# ── 用法说明 ──────────────────────────────────────────────────────────────────
usage() {
    cat <<EOF
用法: $(basename "$0") [选项]

选项:
  --host-path PATH      宿主机目录路径（必填）
  --vm-name NAME        虚拟机名称（必填）
  --vm-path PATH        虚拟机内路径（默认: /root/.openclaw/workspace）
  --mode MODE           同步模式: sync 或 mount（默认: sync）
  --direction DIR       同步方向: push|pull|bidirectional（默认: bidirectional）
  --exclude FILE        rsync 排除文件
  --dry-run             试运行模式，不实际执行
  -h, --help            显示帮助信息

示例:
  # 双向同步
  $(basename "$0") --host-path /data/openclaw-workspace --vm-name my-vm

  # 仅推送到 VM
  $(basename "$0") --host-path /data/openclaw-workspace --vm-name my-vm --direction push

  # 通过 CinaClaw mount 挂载
  $(basename "$0") --host-path /data/openclaw-workspace --vm-name my-vm --mode mount
EOF
    exit 0
}

# ── 参数解析 ──────────────────────────────────────────────────────────────────
parse_args() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --host-path)  HOST_PATH="$2"; shift 2 ;;
            --vm-name)    VM_NAME="$2"; shift 2 ;;
            --vm-path)    VM_PATH="$2"; shift 2 ;;
            --mode)       MODE="$2"; shift 2 ;;
            --direction)  DIRECTION="$2"; shift 2 ;;
            --exclude)    EXCLUDE_FILE="$2"; shift 2 ;;
            --dry-run)    DRY_RUN="true"; shift ;;
            -h|--help)    usage ;;
            *)
                log_error "未知参数: $1"
                usage
                ;;
        esac
    done
}

# ── 前置检查 ──────────────────────────────────────────────────────────────────
validate() {
    if [ -z "$HOST_PATH" ]; then
        log_error "缺少 --host-path 参数"
        exit 1
    fi

    if [ -z "$VM_NAME" ]; then
        log_error "缺少 --vm-name 参数"
        exit 1
    fi

    if [ ! -d "$HOST_PATH" ]; then
        log_step "创建宿主机目录: $HOST_PATH"
        mkdir -p "$HOST_PATH"
    fi

    # 验证模式
    case "$MODE" in
        sync|mount) ;;
        *)
            log_error "不支持的模式: $MODE（可选: sync, mount）"
            exit 1
            ;;
    esac

    # 验证方向
    case "$DIRECTION" in
        push|pull|bidirectional) ;;
        *)
            log_error "不支持的方向: $DIRECTION（可选: push, pull, bidirectional）"
            exit 1
            ;;
    esac

    log_info "配置验证通过"
    log_info "  宿主机路径: $HOST_PATH"
    log_info "  虚拟机名称: $VM_NAME"
    log_info "  虚拟机路径: $VM_PATH"
    log_info "  同步模式:   $MODE"
    log_info "  同步方向:   $DIRECTION"
}

# ── Rsync 同步 ────────────────────────────────────────────────────────────────
sync_push() {
    log_step "推送到虚拟机 $VM_NAME:$VM_PATH ..."

    local rsync_opts=(-avz --delete)
    if [ -n "$EXCLUDE_FILE" ] && [ -f "$EXCLUDE_FILE" ]; then
        rsync_opts+=(--exclude-from="$EXCLUDE_FILE")
    fi
    if [ "$DRY_RUN" = "true" ]; then
        rsync_opts+=(--dry-run)
    fi

    # 使用 CinaClaw 的 SSH 信息获取连接参数
    # 这里假设 SSH 已配置或通过 cinaclaw CLI 获取
    local ssh_target
    ssh_target="${VM_NAME}:${VM_PATH}"

    rsync "${rsync_opts[@]}" "$HOST_PATH/" "$ssh_target/" 2>/dev/null || {
        log_error "rsync push 失败"
        log_error "请确认虚拟机 SSH 连接已配置"
        exit 1
    }

    log_info "推送完成"
}

sync_pull() {
    log_step "从虚拟机 $VM_NAME:$VM_PATH 拉取..."

    local rsync_opts=(-avz)
    if [ -n "$EXCLUDE_FILE" ] && [ -f "$EXCLUDE_FILE" ]; then
        rsync_opts+=(--exclude-from="$EXCLUDE_FILE")
    fi
    if [ "$DRY_RUN" = "true" ]; then
        rsync_opts+=(--dry-run)
    fi

    local ssh_target
    ssh_target="${VM_NAME}:${VM_PATH}/"

    rsync "${rsync_opts[@]}" "$ssh_target" "$HOST_PATH/" 2>/dev/null || {
        log_error "rsync pull 失败"
        log_error "请确认虚拟机 SSH 连接已配置"
        exit 1
    }

    log_info "拉取完成"
}

sync_bidirectional() {
    log_step "双向同步..."

    # 先拉取（避免覆盖远程新文件）
    sync_pull

    # 再推送
    sync_push

    log_info "双向同步完成"
}

do_sync() {
    case "$DIRECTION" in
        push)          sync_push ;;
        pull)          sync_pull ;;
        bidirectional) sync_bidirectional ;;
    esac
}

# ── CinaClaw 挂载 ─────────────────────────────────────────────────────────────
do_mount() {
    log_step "通过 CinaClaw 挂载 $HOST_PATH → $VM_NAME:$VM_PATH ..."

    if ! command -v cinaclaw &>/dev/null; then
        log_error "cinaclaw CLI 未安装"
        exit 1
    fi

    local mount_args=()
    mount_args+=(mount "$HOST_PATH")
    mount_args+=("--target-instance" "$VM_NAME")
    mount_args+=("--target-path" "$VM_PATH")

    if [ "$DRY_RUN" = "true" ]; then
        log_info "试运行模式: cinaclaw ${mount_args[*]}"
        return 0
    fi

    cinaclaw "${mount_args[@]}" || {
        log_error "CinaClaw 挂载失败"
        exit 1
    }

    log_info "挂载成功"
    log_info "  $HOST_PATH → $VM_NAME:$VM_PATH"
}

# ── 主流程 ─────────────────────────────────────────────────────────────────────
main() {
    parse_args "$@"
    validate

    case "$MODE" in
        sync)  do_sync ;;
        mount) do_mount ;;
    esac

    log_info "操作完成"
}

main "$@"
