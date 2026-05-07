# Local Data And Repository Workflow Spec

## 1. Background

The project needs a clearly documented local data setup and a repeatable repository commit workflow. Configuration must tell developers where to set database, Redis, DeepSeek, agent, stock source, and repository values.

## 2. Goals

- Provide local PostgreSQL and Redis startup instructions.
- Record PostgreSQL tables and Redis key fields under `db/`.
- Add repository configuration without hard-coding a remote.
- Add a script that asks before staging, committing, and optionally pushing changes.
- Keep secrets out of committed files.

## 3. Non-Goals

- Do not auto-commit without explicit confirmation.
- Do not auto-push unless the user confirms.
- Do not store real API keys or brokerage credentials.

## 4. Functional Scope

### Must Have

- `.env.example` includes `DEEPSEEK_API_KEY`, `DATABASE_URL`, `REDIS_ADDR`, `REPOSITORY_REMOTE_URL`, and `REPOSITORY_BRANCH`.
- `config/backend.example.json` includes repository metadata.
- `config/repository.example.json` exists for the commit workflow.
- `db/schema.sql` lists local PostgreSQL tables.
- `db/redis-keys.md` lists Redis keys, fields, TTLs, and purpose.
- `scripts/commit-changes.sh` asks before committing.

## 5. Testing

- Shell syntax check for scripts.
- Go config tests for repository config.
- Docker compose config validation.
- Attempt to start PostgreSQL and Redis locally.

## 6. Acceptance Criteria

- A developer can find all editable config files from `docs/configuration.md`.
- A developer can run one command to start PostgreSQL/Redis and one command to apply migrations.
- A developer can run one script that asks whether to commit changes.
