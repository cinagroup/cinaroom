#!/bin/bash
# Start script for CinaSeek backend
# Usage: ./start-backend.sh [--env-file /path/to/.env]

set -e

WORKDIR="/root/.openclaw/workspace/cinaroom/backend"
ENV_FILE="${1:-$WORKDIR/.env}"

cd "$WORKDIR"

# Load env vars
if [ -f "$ENV_FILE" ]; then
    set -a
    source "$ENV_FILE"
    set +a
fi

exec ./bin/cinaseek-backend
