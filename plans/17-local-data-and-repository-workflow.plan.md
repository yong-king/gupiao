# 17 Local Data And Repository Workflow Plan

## Status

completed

## Objective

补齐本地 PostgreSQL/Redis 数据层说明、配置文件入口，以及修改后询问提交的仓库工作流。

## Scope

- `.env.example` 和 JSON 配置文件补充 DeepSeek、数据库、Redis、仓库设置。
- `db/` 本地数据说明文件。
- 可提示用户确认的提交脚本。
- Docker Compose 和脚本语法验证。

## Out Of Scope

- 不自动下单。
- 不自动提交或推送未经确认的修改。
- 不在仓库中保存真实 API key、交易账户密码或 token。

## Dependencies

- [10 Configuration Management](</Users/youngking/Documents/jijin/plans/10-configuration-management.plan.md>)
- [13 One Click Deployment](</Users/youngking/Documents/jijin/plans/13-one-click-deployment.plan.md>)
- [16 Operable Console Persistence And DeepSeek Config](</Users/youngking/Documents/jijin/plans/16-operable-console-persistence-config.plan.md>)

## Required Specs

- Local Data And Repository Workflow Spec

## Tasks

- [x] 增加仓库配置文件。
- [x] 增加 `.env.example` 中 DeepSeek、数据库、Redis、仓库配置项。
- [x] 增加本地 PostgreSQL schema 文件。
- [x] 增加 Redis key/字段说明。
- [x] 增加提交前询问的脚本。
- [x] 验证 Docker PostgreSQL/Redis。
- [x] 跑通相关测试。

## Testing Gate

- `sh -n scripts/commit-changes.sh`
- `sh -n deploy/scripts/apply-migrations.sh`
- `docker compose -f deploy/docker-compose.yml config`
- `docker compose -f deploy/docker-compose.yml up -d postgres redis`
- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test && npm run typecheck`

## Completion Criteria

- 配置文件位置清晰。
- 数据库表和 Redis 字段有本地文件记录。
- 提交脚本不会绕过用户确认。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `.env.example`, `.gitignore`, `config/backend.example.json`, `config/backend.local.json`, `config/repository.example.json`, `config/repository.local.json`, `agent/config/agent.local.json`, `db/README.md`, `db/schema.sql`, `db/redis-keys.md`, `docs/configuration.md`, `deploy/docker-compose.yml`, `deploy/README.md`, `scripts/commit-changes.sh`.
- Repository workflow: `scripts/commit-changes.sh` asks before commit and asks again before push; it can initialize a Git repository only after confirmation.
- Docker validation passed: `docker compose -f deploy/docker-compose.yml config`; `docker compose -f deploy/docker-compose.yml up -d postgres redis`; `sh deploy/scripts/apply-migrations.sh`; Redis `PING`.
- Test commands passed: `sh -n scripts/commit-changes.sh`; `sh -n deploy/scripts/apply-migrations.sh`; `cd backend && env GOCACHE=$PWD/.gocache go test ./...`; `cd frontend && npm test && npm run typecheck`.
- Remaining setup: fill `REPOSITORY_REMOTE_URL` in `config/repository.local.json` or env before pushing.
