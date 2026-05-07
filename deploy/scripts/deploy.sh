#!/usr/bin/env sh
set -eu

ROOT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/deploy/docker-compose.yml"

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required" >&2
  exit 1
fi

docker compose -f "$COMPOSE_FILE" config >/dev/null
docker compose -f "$COMPOSE_FILE" up -d --build

echo "Jijin deployed."
echo "Frontend: http://127.0.0.1:8081"
echo "Backend:  http://127.0.0.1:8080/healthz"
echo "Agent:    http://127.0.0.1:8090/healthz"
