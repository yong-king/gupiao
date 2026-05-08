# Configuration

## Files

- Backend local config: `config/backend.local.json`
- Backend example: `config/backend.example.json`
- Agent local config: `agent/config/agent.local.json`
- Agent example: `agent/config/agent.example.json`
- Repository local config: `config/repository.local.json`
- Repository commit workflow example: `config/repository.example.json`
- Deployment env example: `.env.example`
- Local database schema and Redis key plan: `db/schema.sql`, `db/redis-keys.md`

## Sensitive Values

Do not commit real secrets.

- User passwords: never stored directly; store salted password hashes only.
- Session tokens: store token hashes only.
- LLM API keys: configure through the env var named by `llm.api_key_env`, for example `DEEPSEEK_API_KEY`.
- Database passwords: inject through `DATABASE_URL` or deployment secrets.
- Broker credentials: not supported in MVP; do not store trading passwords.

## Backend

Set `JIJIN_BACKEND_CONFIG` to point to a JSON config file. Environment variables override file values:

```bash
export JIJIN_BACKEND_CONFIG=config/backend.local.json
```

- `GO_BACKEND_ADDR`
- `DATABASE_URL`
- `REDIS_ADDR`
- `AGENT_URL`
- `REPOSITORY_REMOTE_URL`
- `REPOSITORY_BRANCH`

`stock_sources[0]` controls the quote provider. The local MVP defaults to:

- `type: "stooq"`
- `base_url: "https://stooq.com/q/l/"`
- low-frequency refresh only

Use a licensed vendor or account-safe read-only integration before production use.

DeepSeek is configured by provider/model plus an environment variable name. The Python agent uses the task-specific model fields for routing:

- `llm.provider`: `deepseek`
- `llm.model`: `deepseek-chat`
- `llm.chat_model`: `deepseek-chat` for dialogue and synthesis
- `llm.flash_model`: `deepseek-v4-flash` for lightweight collection/crawling tasks
- `llm.pro_model`: `deepseek-v4-pro` for risk review and heavier reasoning tasks
- `llm.api_key_env`: `DEEPSEEK_API_KEY`

Do not put the actual key in any JSON file. Set it in the shell or server secret manager:

```bash
export DEEPSEEK_API_KEY="..."
```

## Local PostgreSQL And Redis

Start middleware:

```bash
docker compose -f deploy/docker-compose.yml up -d postgres redis
```

Apply tables:

```bash
sh deploy/scripts/apply-migrations.sh
```

Inspect tables:

```bash
docker compose -f deploy/docker-compose.yml exec postgres psql -U jijin -d jijin -c "\\dt"
```

Local field documentation is also committed for review:

- PostgreSQL tables: `db/schema.sql`
- Redis keys and fields: `db/redis-keys.md`

## Repository Commit Workflow

Set the remote in `config/repository.local.json` or through env vars:

```bash
export REPOSITORY_REMOTE_URL="git@github.com:your-org/jijin.git"
export REPOSITORY_BRANCH="main"
```

After a code change, run:

```bash
sh scripts/commit-changes.sh
```

The script asks before staging/committing and asks again before pushing. If the directory is not yet a Git repository, it initializes one only after you confirm.

## Agent

Set `JIJIN_AGENT_CONFIG` to point to a JSON config file. Environment variables override file values:

```bash
export JIJIN_AGENT_CONFIG=agent/config/agent.local.json
```

- `AGENT_HOST`
- `AGENT_PORT`
