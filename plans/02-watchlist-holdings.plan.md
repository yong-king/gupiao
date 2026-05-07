# 02 Watchlist And Holdings Plan

## Status

completed

## Objective

支持用户创建股票池，并通过手动录入或 CSV 导入持仓数据，为后续刷新、规则和提醒提供监控对象。

## Scope

- watchlist CRUD。
- watchlist symbol 添加、删除、去重和市场字段校验。
- 手动录入持仓。
- CSV 导入持仓。
- 持仓与股票池或监控范围关联。
- 前端股票池和持仓校验逻辑。

## Out Of Scope

- 不直接连接券商账户。
- 不保存交易密码。
- 不实现行情刷新。
- 不实现提醒规则触发。
- 不实现自动交易。

## Dependencies

- [01 Backend Foundation](</Users/youngking/Documents/jijin/plans/01-backend-foundation.plan.md>)

## Required Specs

- Watchlist Spec
- Holdings Import Spec
- Watchlist Holdings Link Spec

## Tasks

- [x] 生成并确认股票池和持仓 spec。
- [x] 实现 watchlist 数据模型和迁移。
- [x] 实现 watchlist API。
- [x] 实现持仓录入和 CSV 导入。
- [x] 实现前端股票池和持仓校验逻辑。
- [x] 添加审计日志。

## Testing Gate

- Go：watchlist CRUD、重复 symbol、字段校验、持仓导入测试通过。
- Python：本 plan 不涉及 Python Agent；无需 Python 测试。
- Vue：前端股票代码校验、持仓导入校验测试或等价验证通过；完整页面在 `08 Frontend Console` 统一实现。
- 数据库迁移验证通过。

## Completion Criteria

- 用户可以创建股票池。
- 用户可以添加股票代码。
- 用户可以手动录入或 CSV 导入持仓。
- 股票池和持仓可作为后续监控输入。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/internal/watchlist/`, `backend/internal/holdings/`, `backend/internal/api/watchlist_handlers.go`, `backend/internal/api/holdings_handlers.go`, `frontend/src/watchlists.js`, `frontend/src/holdings.js`.
- Migration files: `backend/migrations/000002_watchlists_holdings.sql`.
- Test files: `backend/internal/watchlist/watchlist_test.go`, `backend/internal/holdings/holdings_test.go`, `backend/internal/api/watchlist_handlers_test.go`, `frontend/tests/watchlists.test.js`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`; `npm test`; `npm run typecheck`.
- Test result: passed.
- Remaining risks: full Vue pages are deferred to `08 Frontend Console`; repositories are in-memory until database-backed repositories are implemented in later persistence work.
