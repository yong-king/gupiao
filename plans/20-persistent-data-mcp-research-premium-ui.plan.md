# 20 Persistent Data MCP Research Premium UI Plan

## Status

in_progress

## Objective

将核心业务数据落到 PostgreSQL，并规划 MCP 互联网研究和高级前端体验。

## Scope

- PostgreSQL 持久化：用户、session、股票池、持仓、提醒规则、提醒事件、通知、只读账户配置、行情快照。
- 登出/登录和后端重启后的数据保留。
- MCP/Firecrawl 股票公司与产品公开信息采集 spec 化。
- 前端视觉和状态反馈升级。

## Out Of Scope

- 不自动交易。
- 不保存券商密码/token/secret。
- 不绕过登录、验证码、付费墙或反爬限制。
- 本 plan 不完成生产级向量数据库检索。

## Dependencies

- [19 Operable Portfolio Alerts Accounts UI](</Users/youngking/Documents/jijin/plans/19-operable-portfolio-alerts-accounts-ui.plan.md>)

## Required Specs

- Persistent Data MCP Research Premium UI Spec

## Tasks

- [x] 增加 PostgreSQL driver 和持久化 store。
- [x] Auth 用户/session 接 PostgreSQL。
- [x] 股票池和股票池 symbol 接 PostgreSQL。
- [x] 持仓接 PostgreSQL。
- [x] 提醒规则、提醒事件、通知接 PostgreSQL。
- [x] 只读账户配置接 PostgreSQL。
- [x] 行情快照读取/保存接 PostgreSQL。
- [x] 端到端验证后端重启后数据仍可读取。
- [ ] MCP/Firecrawl 股票研究工作流实现。
- [ ] 研究结果写入 `rag_documents`。
- [ ] DeepSeek 分析运行记录写入 `llm_analysis_runs`。
- [ ] 前端进一步高级化：全局 loading、toast、空状态、数据密度优化。

## Testing Gate

- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test && npm run typecheck`
- `sh deploy/scripts/apply-migrations.sh`
- Persistence smoke：注册用户，创建股票池/持仓/规则/账户/行情，重启后端，登录同一用户，确认数据仍存在。
- MCP research smoke：对 `CN:000821` 采集公开公司/产品信息并记录 source。

## Completion Criteria

- 核心用户数据不因退出登录或后端重启丢失。
- MCP 研究链路有可测试入口和数据落库。
- 前端操作状态清晰且专业。
- Delivery Notes 记录测试结果和剩余风险。

## Delivery Notes

- Persistence implementation files: `backend/internal/persistence/postgres.go`, `backend/internal/auth/auth.go`, `backend/internal/api/router.go`, `backend/internal/api/watchlist_handlers.go`, `backend/internal/api/holdings_handlers.go`, `backend/internal/api/alert_handlers.go`, `backend/internal/api/refresh_handlers.go`, `backend/internal/api/market_handlers.go`, `backend/internal/api/accounts_handlers.go`.
- Tests passed: `cd backend && env GOCACHE=$PWD/.gocache go test ./...`; `cd frontend && npm test && npm run typecheck`; `sh deploy/scripts/apply-migrations.sh`.
- Persistence smoke passed: after backend restart, login with the same account returned persisted watchlists, holdings, alert rules, alert events, notifications, account config, and snapshots.
- MCP/Firecrawl smoke attempted for `000821 京山轻机 产品 业务 股票 公司 2026`, but the MCP call returned `Unauthorized: Invalid token`. MCP research remains blocked until Firecrawl/MCP credentials are configured.
