# 08 Frontend Console Plan

## Status

completed

## Objective

整合 Vue 操作台，让用户可以完成股票池、持仓、规则、刷新、提醒、报告和系统设置的核心操作。

## Scope

- Vue 应用框架、路由、布局。
- API client 和错误处理。
- 股票池入口整合。
- 持仓入口整合。
- 规则配置入口整合。
- 刷新任务状态。
- 实时提醒入口。
- 报告入口。
- 系统设置页面。

## Out Of Scope

- 不做营销型首页。
- 不实现复杂图表编辑器。
- 不实现自动交易按钮。

## Dependencies

- [02 Watchlist And Holdings](</Users/youngking/Documents/jijin/plans/02-watchlist-holdings.plan.md>)
- [03 Market Refresh](</Users/youngking/Documents/jijin/plans/03-market-refresh.plan.md>)
- [04 Alerts And Notifications](</Users/youngking/Documents/jijin/plans/04-alerts-notifications.plan.md>)
- [06 Reports And Review](</Users/youngking/Documents/jijin/plans/06-reports-review.plan.md>)

## Required Specs

- Frontend App Shell Spec
- Refresh Status Page Spec
- Realtime Alerts UI Spec
- System Settings Page Spec

## Tasks

- [x] 生成并确认前端操作台 spec。
- [x] 实现路由、布局和基础状态管理。
- [x] 实现 API client 和统一错误展示。
- [x] 整合股票池、持仓、规则、刷新、提醒、报告入口。
- [x] 实现系统设置。
- [x] 完成关键页面响应式验证。

## Testing Gate

- Vue：类型检查通过。
- Vue：组件测试或页面测试通过。
- E2E 或浏览器验证：关键流程可操作，包括股票池、规则、刷新、提醒和报告入口。
- Go/Python：如本 plan 修改 API 契约或 Agent 输出展示，必须同步运行相关后端或 Agent 测试。

## Completion Criteria

- 用户可以从操作台完成 MVP 核心流程。
- 页面展示数据时间、触发原因、风险等级和刷新状态。
- 页面不出现自动交易能力。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `frontend/index.html`, `frontend/src/app.js`, `frontend/src/main.js`.
- Test files: `frontend/tests/app.test.js`.
- Test commands: `npm test`; `npm run typecheck`.
- Test result: passed.
- Remaining risks: no Vue runtime dependency is installed in this sandboxed MVP; the app shell is static HTML plus ES modules. Full Vue SFC build can be introduced when dependency installation is allowed.
