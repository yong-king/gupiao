# 12 Stock Source Collection Plan

## Status

completed

## Objective

为股票平台信息采集建立可替换的数据源接口，支持 API 数据源和爬虫解析器两种路径。

## Scope

- 股票平台数据源接口。
- HTML/JSON 解析器。
- Mock 平台源。
- 超时、限频和 source_refs。
- 文档说明优先使用官方 API，爬虫仅作为补充。

## Out Of Scope

- 不绕过登录、验证码或反爬限制。
- 不抓取需要授权的账户页面。

## Dependencies

- [03 Market Refresh](</Users/youngking/Documents/jijin/plans/03-market-refresh.plan.md>)
- [10 Configuration Management](</Users/youngking/Documents/jijin/plans/10-configuration-management.plan.md>)

## Required Specs

- Stock Source Collection Spec

## Tasks

- [x] 生成采集 spec。
- [x] 实现数据源接口。
- [x] 实现 HTML/JSON 解析器。
- [x] 接入 mock 平台源测试。
- [x] 补限频和合规文档。

## Testing Gate

- Go 解析器测试通过。
- 数据源错误和缺失字段测试通过。

## Completion Criteria

- 可以从平台响应解析股票行情。
- 保留数据来源和数据时间。
- 不包含绕过风控行为。

## Delivery Notes

- Implementation files: `backend/internal/stocksource/`, `docs/stock-sources.md`.
- Test files: `backend/internal/stocksource/stocksource_test.go`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`.
- Test result: passed.
- Remaining risks: real platform adapters are not implemented; parser only supports allowed public JSON/HTML shapes and must respect platform terms/rate limits.
