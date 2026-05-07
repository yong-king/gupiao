# 07 Account Monitoring Plan

## Status

completed

## Objective

在只读、低频、可审计的边界内支持账户级监控，优先通过 CSV 或手动导入持仓数据实现。

## Scope

- 账户别名和刷新模式配置。
- 账户只读能力描述。
- CSV 或手动导入账户持仓。
- 账户刷新冷却策略。
- 持仓集中度、单票浮盈浮亏、组合回撤观察。
- 账户监控前端页面。

## Out Of Scope

- 不保存交易密码。
- 不自动登录券商。
- 不自动买入、卖出或下单。
- 不绕过券商风控。
- 不做高频账户刷新。

## Dependencies

- [02 Watchlist And Holdings](</Users/youngking/Documents/jijin/plans/02-watchlist-holdings.plan.md>)
- [03 Market Refresh](</Users/youngking/Documents/jijin/plans/03-market-refresh.plan.md>)
- [04 Alerts And Notifications](</Users/youngking/Documents/jijin/plans/04-alerts-notifications.plan.md>)

## Required Specs

- Account Settings Spec
- Account Holdings Import Spec
- Portfolio Risk Spec

## Tasks

- [x] 生成并确认账户监控 spec。
- [x] 实现账户配置模型和 API。
- [x] 实现账户持仓导入。
- [x] 实现账户刷新冷却。
- [x] 实现组合风险计算。
- [x] 实现账户监控页面。
- [x] 添加审计日志。

## Testing Gate

- Go：账户配置校验、敏感字段不落库、持仓导入、账户冷却、组合风险计算测试通过。
- Vue：账户配置页、导入错误、风险展示测试或浏览器验证通过。
- Python：如果使用 Agent 解释组合风险，必须通过 schema 和 missing_data 测试；否则无需 Python 测试。
- 安全：确认没有自动交易 API 或下单字段。

## Completion Criteria

- 账户监控为只读。
- 账户数据来源和数据时间可追溯。
- 刷新低频、可冷却、可审计。
- 不出现任何下单能力。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/internal/accounts/`, `frontend/src/accounts.js`.
- Migration files: `backend/migrations/000007_account_monitoring.sql`.
- Test files: `backend/internal/accounts/accounts_test.go`, `frontend/tests/accounts.test.js`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`; `npm test`; `npm run typecheck`.
- Test result: passed.
- Remaining risks: real broker read-only APIs are deferred; MVP account monitoring uses manual/CSV holdings and rejects sensitive credential metadata.
