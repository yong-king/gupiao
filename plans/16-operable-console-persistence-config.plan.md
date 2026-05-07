# 16 Operable Console Persistence And DeepSeek Config Plan

## Status

completed

## Objective

完善前端可操作性、真实行情测试保存链路、DeepSeek 安全配置，以及 PostgreSQL/Redis 和后续 RAG/LLM 分析所需表结构。

## Scope

- DeepSeek provider/model/API key env 配置。
- 移除配置文件中的真实 API key。
- watchlist/holdings/rules 的可见查询和前端展示。
- 真实行情 collect-and-save API。
- 系统依赖状态 API。
- stock_profiles、rag_documents、llm_analysis_runs 数据库表。
- 前端页面操作反馈和设置页配置展示。

## Out Of Scope

- 不实现自动交易。
- 不实现完整向量库 RAG 检索。
- 不绕过需要登录、验证码或反爬的平台。

## Dependencies

- [10 Configuration Management](</Users/youngking/Documents/jijin/plans/10-configuration-management.plan.md>)
- [11 Auth Registration Login](</Users/youngking/Documents/jijin/plans/11-auth-registration-login.plan.md>)
- [15 Real Market Analytics And RAG](</Users/youngking/Documents/jijin/plans/15-real-market-analytics-rag.plan.md>)

## Required Specs

- Operable Console Persistence And DeepSeek Config Spec

## Tasks

- [x] 更新 DeepSeek 配置，API key 改为 `DEEPSEEK_API_KEY` 环境变量。
- [x] 增加数据库表迁移。
- [x] 增加 watchlist/holdings 查询 API。
- [x] 增加真实行情测试并保存 API。
- [x] 增加系统依赖状态 API。
- [x] 完善前端每个页面的操作反馈和数据展示。
- [x] 启动/验证本地 PostgreSQL 和 Redis。
- [x] 跑通测试和真实行情 smoke。

## Testing Gate

- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
- `cd frontend && npm test`
- `cd frontend && npm run typecheck`
- API smoke：collect AAPL real quote, then read snapshots/daily changes/profile.
- Docker smoke：PostgreSQL and Redis containers are available, and migrations can be inspected/applied.

## Completion Criteria

- 页面上的主要按钮都有可见效果。
- 真实行情可以先测试、再保存、再展示。
- DeepSeek key 不出现在仓库配置文件。
- 数据库表结构覆盖股票池、持仓、规则、快照、每日涨跌、profile、RAG 文档和 LLM 分析。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/internal/api/`, `backend/internal/config/`, `backend/internal/watchlist/`, `backend/migrations/000009_market_analytics.sql`, `backend/migrations/000010_rag_llm_profiles.sql`, `frontend/src/main.js`, `config/backend.example.json`, `agent/config/agent.example.json`, `docs/configuration.md`.
- Docker middleware: `postgres:16-alpine` and `redis:7-alpine` started through `docker compose -f deploy/docker-compose.yml up -d postgres redis`.
- Database migration smoke: `sh deploy/scripts/apply-migrations.sh` created 20 PostgreSQL tables; Redis smoke returned `PONG`.
- Test commands passed: `cd backend && env GOCACHE=$PWD/.gocache go test ./...`; `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`; `cd frontend && npm test && npm run typecheck`.
- API smoke passed: register/login, create watchlist, add AAPL, import holdings, create alert rule, collect and save AAPL real quote from Stooq, read snapshots, daily changes, profile, and dependency status.
- Remaining risk: current API repositories are still in-memory; PostgreSQL/Redis schema and config are ready for replacing in-memory repositories with persistent implementations.
