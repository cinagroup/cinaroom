#!/bin/bash
# Start script for CinaSeek WebSocket service
set -e

WORKDIR="/root/.openclaw/workspace/cinaroom/websocket/backend"
cd "$WORKDIR"

export WS_PORT=8081
exec ./bin/cinaseek-ws
