#!/usr/bin/env sh
set -eu

COMPOSE_FILE="${COMPOSE_FILE:-deploy/docker-compose.yml}"
DB_SERVICE="${DB_SERVICE:-postgres}"
DB_USER="${POSTGRES_USER:-jijin}"
DB_NAME="${POSTGRES_DB:-jijin}"

for file in backend/migrations/*.sql; do
  echo "applying ${file}"
  docker compose -f "${COMPOSE_FILE}" exec -T "${DB_SERVICE}" psql -U "${DB_USER}" -d "${DB_NAME}" < "${file}"
done
