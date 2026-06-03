#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REMOTE_HOST="${DEPLOY_HOST:-qiqiqu.cn}"
REMOTE_PORT="${DEPLOY_PORT:-6000}"
REMOTE_USER="${DEPLOY_USER:-shynin}"
REMOTE_ROOT="${DEPLOY_ROOT:-/home/shynin/billSoftware}"
REMOTE_COMPOSE_FILE="$REMOTE_ROOT/docker-compose.prod.yml"
DOMAIN="${DEPLOY_DOMAIN:-www.qiqiqu.cn}"
FRONTEND_BASE="${BILL_PUBLIC_BASE:-/pc/}"
FORCE_CONFIG="${DEPLOY_FORCE_CONFIG:-0}"
DEPLOY_PASSWORD="${DEPLOY_PASSWORD:-}"
TMP_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}

require_cmd npm
require_cmd python3
require_cmd rsync
require_cmd ssh
require_cmd scp
require_cmd zsh

SSH_TARGET="$REMOTE_USER@$REMOTE_HOST"
RSYNC_RSH="ssh -o StrictHostKeyChecking=no -p $REMOTE_PORT"
if [[ -n "$DEPLOY_PASSWORD" ]]; then
  require_cmd sshpass
  SSH_CMD=(sshpass -p "$DEPLOY_PASSWORD" ssh -o StrictHostKeyChecking=no -p "$REMOTE_PORT" "$SSH_TARGET")
  SCP_CMD=(sshpass -p "$DEPLOY_PASSWORD" scp -P "$REMOTE_PORT" -o StrictHostKeyChecking=no)
  RSYNC_RSH="sshpass -p $DEPLOY_PASSWORD $RSYNC_RSH"
else
  SSH_CMD=(ssh -o StrictHostKeyChecking=no -p "$REMOTE_PORT" "$SSH_TARGET")
  SCP_CMD=(scp -P "$REMOTE_PORT" -o StrictHostKeyChecking=no)
fi

render_backend_config() {
  local source_config="$ROOT_DIR/backend/configs/app.yaml"
  if [[ ! -f "$source_config" ]]; then
    source_config="$ROOT_DIR/deploy/backend/app.prod.example.yaml"
  fi

  cp "$source_config" "$TMP_DIR/app.yaml"
  python3 - "$TMP_DIR/app.yaml" <<'PY'
from pathlib import Path
import sys

path = Path(sys.argv[1])
lines = path.read_text().splitlines()
section = None
result = []

for line in lines:
    stripped = line.strip()
    if not line.startswith(" "):
        section = stripped[:-1] if stripped.endswith(":") else None
        result.append(line)
        continue

    if section == "server" and stripped.startswith("port:"):
        result.append("  port: 8080")
        continue
    if section == "server" and stripped.startswith("mode:"):
        result.append("  mode: release")
        continue
    if section == "database" and stripped.startswith("host:"):
        result.append("  host: mysql")
        continue
    if section == "database" and stripped.startswith("port:"):
        result.append("  port: 3306")
        continue

    result.append(line)

path.write_text("\n".join(result) + "\n")
PY
}

render_backend_env() {
  local source_env="$ROOT_DIR/backend/.env.example"
  if [[ -f "$ROOT_DIR/backend/.env" ]]; then
    source_env="$ROOT_DIR/backend/.env"
  fi

  cp "$source_env" "$TMP_DIR/backend.env"
}

echo "==> Building frontend with base $FRONTEND_BASE"
(
  cd "$ROOT_DIR/frontend"
  mkdir -p "$TMP_DIR/frontend-dist"
  BILL_PUBLIC_BASE="$FRONTEND_BASE" npx webpack --config webpack.config.js --mode production --output-path "$TMP_DIR/frontend-dist"
)

echo "==> Rendering deployment config"
render_backend_config
render_backend_env

echo "==> Building backend linux binary"
mkdir -p "$TMP_DIR/backend-app"
rsync -a --delete \
  --exclude '.git' \
  --exclude 'configs/app.yaml' \
  --exclude 'tmp' \
  --exclude 'bin' \
  "$ROOT_DIR/backend/" \
  "$TMP_DIR/backend-app/"
mkdir -p "$TMP_DIR/backend-app/bin"
(
  cd "$TMP_DIR/backend-app"
  zsh -ic 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bin/billserver ./cmd/server'
)

echo "==> Preparing remote directories"
"${SSH_CMD[@]}" "
  set -e
  mkdir -p '$REMOTE_ROOT/frontend/dist' \
           '$REMOTE_ROOT/frontend/nginx' \
           '$REMOTE_ROOT/backend/app' \
           '$REMOTE_ROOT/backend/configs' \
           '$REMOTE_ROOT/nginx/conf.d' \
           '$REMOTE_ROOT/nginx/certs'
"

echo "==> Syncing frontend assets"
rsync -az --delete -e "$RSYNC_RSH" \
  "$TMP_DIR/frontend-dist/" \
  "$SSH_TARGET:$REMOTE_ROOT/frontend/dist/"

echo "==> Syncing backend source"
rsync -az --delete \
  --exclude '.git' \
  --exclude 'configs/app.yaml' \
  --exclude 'tmp' \
  -e "$RSYNC_RSH" \
  "$TMP_DIR/backend-app/" \
  "$SSH_TARGET:$REMOTE_ROOT/backend/app/"

echo "==> Syncing deployment manifests"
rsync -az -e "$RSYNC_RSH" \
  "$ROOT_DIR/docker-compose.prod.yml" \
  "$SSH_TARGET:$REMOTE_COMPOSE_FILE"
rsync -az -e "$RSYNC_RSH" \
  "$ROOT_DIR/backend/.env.example" \
  "$SSH_TARGET:$REMOTE_ROOT/backend/.env.example"
rsync -az -e "$RSYNC_RSH" \
  "$ROOT_DIR/deploy/nginx/gateway.conf" \
  "$SSH_TARGET:$REMOTE_ROOT/nginx/conf.d/default.conf"
rsync -az -e "$RSYNC_RSH" \
  "$ROOT_DIR/deploy/frontend-nginx/default.conf" \
  "$SSH_TARGET:$REMOTE_ROOT/frontend/nginx/default.conf"

echo "==> Uploading runtime config"
if [[ "$FORCE_CONFIG" == "1" ]]; then
  "${SCP_CMD[@]}" "$TMP_DIR/app.yaml" "$SSH_TARGET:$REMOTE_ROOT/backend/configs/app.yaml"
  "${SCP_CMD[@]}" "$TMP_DIR/backend.env" "$SSH_TARGET:$REMOTE_ROOT/backend/.env"
else
  "${SCP_CMD[@]}" "$TMP_DIR/app.yaml" "$SSH_TARGET:$REMOTE_ROOT/backend/configs/app.yaml.tmp"
  "${SCP_CMD[@]}" "$TMP_DIR/backend.env" "$SSH_TARGET:$REMOTE_ROOT/backend/.env.tmp"
  "${SSH_CMD[@]}" "
    set -e
    if [ ! -f '$REMOTE_ROOT/backend/configs/app.yaml' ]; then
      mv '$REMOTE_ROOT/backend/configs/app.yaml.tmp' '$REMOTE_ROOT/backend/configs/app.yaml'
    else
      rm -f '$REMOTE_ROOT/backend/configs/app.yaml.tmp'
    fi
    if [ ! -f '$REMOTE_ROOT/backend/.env' ]; then
      mv '$REMOTE_ROOT/backend/.env.tmp' '$REMOTE_ROOT/backend/.env'
    else
      rm -f '$REMOTE_ROOT/backend/.env.tmp'
    fi
  "
fi

echo "==> Ensuring TLS certificates"
"${SSH_CMD[@]}" "
  set -e
  if [ ! -f '$REMOTE_ROOT/nginx/certs/www.qiqiqu.cn.pem' ] || [ ! -f '$REMOTE_ROOT/nginx/certs/www.qiqiqu.cn.key' ]; then
    unzip -o '$REMOTE_ROOT/24395009_www.qiqiqu.cn_nginx.zip' -d '$REMOTE_ROOT/nginx/certs' >/dev/null
  fi
"

echo "==> Deploying containers"
"${SSH_CMD[@]}" "
  set -e
  cd '$REMOTE_ROOT'
  docker compose -f '$REMOTE_COMPOSE_FILE' up -d --build
  docker compose -f '$REMOTE_COMPOSE_FILE' ps
"

echo "==> Deployment completed for $DOMAIN"
