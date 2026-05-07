# 15 Real Market Analytics And RAG Plan

## Status

completed

## Objective

接入真实只读行情源，展示行情曲线，记录每日涨跌，并产出可供后续 RAG 分析的股票涨跌和公司/产品上下文。

## Scope

- Stooq CSV 行情 provider。
- 快照历史和每日涨跌 API。
- 公司/产品 profile 与保守分析 API。
- 前端分析报告页行情曲线、涨跌列表、公司产品分析。
- RAG-ready 日涨跌文本字段。

## Out Of Scope

- 不接入需要登录、验证码或绕过反爬的平台。
- 不实现向量数据库和完整 RAG 检索链。
- 不实现自动交易或确定性买卖建议。
- 不在本 plan 引入第三方前端图表库。

## Dependencies

- [03 Market Refresh](</Users/youngking/Documents/jijin/plans/03-market-refresh.plan.md>)
- [06 Reports And Review](</Users/youngking/Documents/jijin/plans/06-reports-review.plan.md>)
- [08 Frontend Console](</Users/youngking/Documents/jijin/plans/08-frontend-console.plan.md>)
- [10 Configuration Management](</Users/youngking/Documents/jijin/plans/10-configuration-management.plan.md>)

## Required Specs

- Real Market Analytics And RAG Spec

## Tasks

- [x] 编写真实行情、曲线、每日涨跌、RAG-ready 和公司/产品分析 spec。
- [x] 实现 Stooq CSV provider 和解析测试。
- [x] 扩展行情快照字段，支持 open/high/low/name。
- [x] 实现每日涨跌聚合和 RAG 文本。
- [x] 实现股票 profile 和保守分析。
- [x] 增加快照、每日涨跌和 profile API。
- [x] 前端报告页展示价格曲线、涨跌记录和公司产品分析。
- [x] 运行 Go、前端测试和浏览器验证。

## Testing Gate

- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test`
- `cd frontend && npm run typecheck`
- 浏览器验证：刷新真实 AAPL 行情后，报告页可展示曲线、每日涨跌和公司/产品分析。

## Completion Criteria

- AAPL 真实行情可以从公开只读源获取。
- 采集后的快照可用于曲线展示。
- 每日涨跌记录含 RAG-ready 文本。
- 公司/产品分析明确不构成投资建议。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/internal/marketdata/stooq.go`, `backend/internal/marketdata/marketdata.go`, `backend/internal/marketdata/profile.go`, `backend/internal/api/market_handlers.go`, `frontend/src/market.js`, `frontend/src/main.js`, `frontend/index.html`, `config/backend.example.json`.
- Test files: `backend/internal/marketdata/marketdata_test.go`, `backend/internal/api/market_handlers_test.go`, `frontend/tests/market.test.js`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`; `npm test`; `npm run typecheck`.
- Test result: passed.
- API smoke: registered `realtest@example.com`, created AAPL watchlist, refreshed real Stooq quote successfully, and verified snapshots, daily changes, and stock profile endpoints.
- Browser verification: report page modules for real quote curve, daily changes, and company/product analysis render in Chrome; after backend restart the browser must log in again before loading data.
- Remaining risks: Stooq is a public quote source and may throttle or change response fields; production should use a licensed market data vendor and persist snapshots in PostgreSQL.
