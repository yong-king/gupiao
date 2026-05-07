# 18 Stock Entry Monitor Calendar Plan

## Status

completed

## Objective

补齐用户可输入股票代码、添加多只股票监控、维护买入/卖出关注价，以及每日涨跌日历展示。

## Scope

- Watchlist symbol 增加买入/卖出关注价。
- Watchlist API 支持关注价写入和重复股票更新。
- 前端股票池页改为输入驱动的股票添加和信息获取。
- 添加股票时按关注价创建提醒规则。
- 报告页增加每日涨跌日历。
- 数据库迁移和本地 schema 同步关注价字段。

## Out Of Scope

- 不自动买入、卖出或下单。
- 不接入券商交易接口。
- 不实现完整持仓盈亏核算。

## Dependencies

- [16 Operable Console Persistence And DeepSeek Config](</Users/youngking/Documents/jijin/plans/16-operable-console-persistence-config.plan.md>)
- [17 Local Data And Repository Workflow](</Users/youngking/Documents/jijin/plans/17-local-data-and-repository-workflow.plan.md>)

## Required Specs

- Stock Entry Monitor Calendar Spec

## Tasks

- [x] 扩展 watchlist symbol 关注价模型。
- [x] 更新 watchlist API 和测试。
- [x] 增加数据库迁移和本地 schema 字段。
- [x] 改造股票池页面为输入股票代码和关注价。
- [x] 加入股票时创建买入/卖出提醒规则。
- [x] 增加每日涨跌日历渲染和测试。
- [x] 跑通测试和真实行情 smoke。

## Testing Gate

- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test && npm run typecheck`
- `sh deploy/scripts/apply-migrations.sh`
- API smoke：register/login, fetch AAPL quote, add/update AAPL monitor prices, verify watchlist row.

## Completion Criteria

- 股票池页可以输入并获取任意支持市场的股票信息。
- 股票池可以保存多只股票和每只股票的买入/卖出关注价。
- 重复添加同一股票会更新关注价。
- 报告页展示每日涨跌日历。
- 测试和 smoke 结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/internal/watchlist/watchlist.go`, `backend/internal/api/watchlist_handlers.go`, `backend/migrations/000002_watchlists_holdings.sql`, `backend/migrations/000011_watchlist_monitor_prices.sql`, `db/schema.sql`, `frontend/src/main.js`, `frontend/src/market.js`, `frontend/index.html`.
- Test files: `backend/internal/watchlist/watchlist_test.go`, `backend/internal/api/watchlist_handlers_test.go`, `frontend/tests/market.test.js`.
- Test commands passed: `cd backend && env GOCACHE=$PWD/.gocache go test ./...`; `cd frontend && npm test && npm run typecheck`; `sh deploy/scripts/apply-migrations.sh`.
- API smoke passed: register/login, collect AAPL quote, create watchlist, add AAPL monitor prices, update AAPL monitor prices, create buy/sell alert rules, and read daily changes.
