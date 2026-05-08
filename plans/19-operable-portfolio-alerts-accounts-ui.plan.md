# 19 Operable Portfolio Alerts Accounts UI Plan

## Status

completed

## Objective

完成持仓、提醒规则、提醒中心、分析报告和账户监控的可操作配置，并提升前端界面质感。

## Scope

- 单条持仓保存 API 和前端表单。
- 提醒规则配置表单和列表展示。
- 提醒中心运行监控检查、刷新、通知已读。
- 分析报告数值摘要和单点曲线可见性。
- 只读账户配置 API 和前端表单。
- 现代化控制台视觉样式。

## Out Of Scope

- 不自动交易。
- 不保存券商密码/token/secret。
- 不实现完整券商同步。

## Dependencies

- [18 Stock Entry Monitor Calendar](</Users/youngking/Documents/jijin/plans/18-stock-entry-monitor-calendar.plan.md>)

## Required Specs

- Operable Portfolio Alerts Accounts UI Spec

## Tasks

- [x] 增加 `POST /api/holdings` 单条持仓保存。
- [x] 增加 `GET/POST /api/accounts` 只读账户配置。
- [x] 增加 `POST /api/notifications/read`。
- [x] 改造持仓页为配置表单和盈亏表格。
- [x] 改造提醒规则页为配置表单和规则表格。
- [x] 改造提醒中心为运行检查、刷新和已读操作。
- [x] 改造分析报告为数值摘要和可见曲线。
- [x] 改造账户监控页为只读配置入口。
- [x] 前端视觉现代化。
- [x] 跑通测试和 smoke。

## Testing Gate

- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test && npm run typecheck`
- API smoke：注册登录，添加 `CN:000821`，保存持仓，保存提醒规则，运行提醒检查，保存只读账户，读取报告数据。

## Completion Criteria

- 持仓、提醒规则、提醒中心、分析报告、账户监控都有实际按钮和可见结果。
- 股票池已加入的股票可以在其他页面继续配置。
- 测试和 smoke 结果写入 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/internal/holdings/holdings.go`, `backend/internal/api/holdings_handlers.go`, `backend/internal/accounts/accounts.go`, `backend/internal/api/accounts_handlers.go`, `backend/internal/api/alert_handlers.go`, `backend/internal/api/router.go`, `frontend/src/main.js`, `frontend/src/market.js`, `frontend/index.html`.
- Test files: `backend/internal/holdings/holdings_test.go`, `backend/internal/api/watchlist_handlers_test.go`, `backend/internal/accounts/accounts_test.go`, `backend/internal/api/accounts_handlers_test.go`, `backend/internal/notifications/notifications_test.go`, `frontend/tests/market.test.js`.
- Test commands passed: `cd backend && env GOCACHE=$PWD/.gocache go test ./...`; `cd frontend && npm test && npm run typecheck`.
- API smoke passed: register/login, collect `CN:000821`, create watchlist, add monitored symbol, save holding, save alert rule, run monitor refresh, list alerts/notifications, mark notification read, save/list read-only account, and read report snapshots/daily changes.
- Frontend served at `http://127.0.0.1:5173/index.html?v=21`.
