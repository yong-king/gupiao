# Local Data Layer

This directory records the local PostgreSQL tables and Redis key plan for development and deployment.

## Configure

- Backend config: `config/backend.example.json`
- Agent config: `agent/config/agent.example.json`
- Repository commit config: `config/repository.example.json`
- Environment overrides: `.env.example`

Never put real API keys or account passwords in committed JSON files. Use environment variables such as `DEEPSEEK_API_KEY`.

## Start Middleware

```bash
docker compose -f deploy/docker-compose.yml up -d postgres redis
sh deploy/scripts/apply-migrations.sh
```

## Inspect

```bash
docker compose -f deploy/docker-compose.yml exec postgres psql -U jijin -d jijin -c "\\dt"
docker compose -f deploy/docker-compose.yml exec redis redis-cli INFO keyspace
```

## Files

- `db/schema.sql`: consolidated PostgreSQL schema for local inspection.
- `db/redis-keys.md`: Redis keys and fields used by the project.
