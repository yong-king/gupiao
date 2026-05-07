# 03 Market Refresh Plan

## Status

completed

## Objective

建立低频、可控、可追溯的行情和数据刷新能力，支持手动刷新和定时刷新。

## Scope

- 数据源接口抽象。
- 行情快照模型和落库。
- refresh_job 创建、执行、失败和状态流转。
- 手动刷新入口。
- 定时刷新入口。
- 同数据源限频、同账户冷却、连续失败退避。

## Out Of Scope

- 不实现复杂实时行情。
- 不实现账户自动登录。
- 不实现提醒规则业务。
- 不实现 Python Agent 分析。

## Dependencies

- [01 Backend Foundation](</Users/youngking/Documents/jijin/plans/01-backend-foundation.plan.md>)
- [02 Watchlist And Holdings](</Users/youngking/Documents/jijin/plans/02-watchlist-holdings.plan.md>)

## Required Specs

- Market Data Source Spec
- Refresh Job Spec
- Rate Limit And Backoff Spec

## Tasks

- [x] 生成并确认刷新相关 spec。
- [x] 实现数据源接口和 mock provider。
- [x] 实现行情快照模型和 repository。
- [x] 实现 refresh_job 状态机。
- [x] 实现手动刷新 API。
- [x] 实现定时刷新 worker 或 scheduler。
- [x] 实现限频、冷却和退避。
- [x] 写入审计日志。

## Testing Gate

- Go：数据源 mock、refresh_job 状态流转、失败重试、冷却命中、限频拒绝测试通过。
- 数据库：行情快照和刷新任务迁移验证通过。
- Vue：如果实现刷新入口或状态展示，必须完成组件测试或浏览器验证。
- Python：本 plan 不涉及 Python Agent；无需 Python 测试。

## Completion Criteria

- 用户可以手动刷新股票池或持仓。
- 系统可以执行低频定时刷新。
- 刷新过程有状态、有错误、有审计。
- 限频和冷却生效。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/internal/marketdata/`, `backend/internal/refresh/`, `backend/internal/ratelimit/`, `backend/internal/api/refresh_handlers.go`.
- Migration files: `backend/migrations/000003_refresh_jobs_market_data.sql`.
- Test files: `backend/internal/marketdata/marketdata_test.go`, `backend/internal/refresh/refresh_test.go`, `backend/internal/ratelimit/ratelimit_test.go`, `backend/internal/api/refresh_handlers_test.go`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`.
- Test result: passed.
- Remaining risks: real external market data providers and Redis-backed distributed rate limiting are deferred; current implementation uses deterministic mock provider and in-memory limiter for MVP testing.
