# Deploy

Deployment files are introduced by the hardening and release plan. This directory is reserved so future plans have a stable location for Docker, compose, and environment-specific assets.

Start local middleware:

```bash
docker compose -f deploy/docker-compose.yml up -d postgres redis
```

If you want backend and agent to use DeepSeek locally, put `DEEPSEEK_API_KEY=...` in the repo-root `.env`. The compose file loads that `.env` into the `backend` and `agent` containers.

Apply local database migrations:

```bash
sh deploy/scripts/apply-migrations.sh
```

Configuration files:

- Backend: `config/backend.example.json`
- Agent: `agent/config/agent.example.json`
- Repository workflow: `config/repository.example.json`
- Database and Redis fields: `db/schema.sql`, `db/redis-keys.md`

Commit workflow:

```bash
sh scripts/commit-changes.sh
```
