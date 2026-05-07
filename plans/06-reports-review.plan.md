# 06 Reports And Review Plan

## Status

completed

## Objective

提供股票详情分析和每日复盘报告，帮助用户查看提醒、风险点、数据来源和需要人工确认的问题。

## Scope

- 股票详情聚合接口。
- analysis_report 模型。
- 每日盘后复盘任务。
- Python Agent 结构化日报生成。
- 前端股票详情页。
- 前端分析报告页。

## Out Of Scope

- 不实现自动交易建议。
- 不实现复杂投研知识图谱。
- 不实现长周期回测。

## Dependencies

- [04 Alerts And Notifications](</Users/youngking/Documents/jijin/plans/04-alerts-notifications.plan.md>)
- [05 Python Agent Analysis](</Users/youngking/Documents/jijin/plans/05-python-agent-analysis.plan.md>)

## Required Specs

- Stock Detail Spec
- Daily Review Report Spec
- Analysis Report Page Spec

## Tasks

- [x] 生成并确认报告相关 spec。
- [x] 实现股票详情聚合接口。
- [x] 实现 analysis_report 数据模型。
- [x] 实现每日复盘任务。
- [x] 实现 Python Agent 日报输出。
- [x] 实现前端股票详情页。
- [x] 实现前端分析报告页。

## Testing Gate

- Go：股票详情聚合、日报任务、无数据报告、Agent 失败降级测试通过。
- Python：日报输出 schema、missing_data、空输入测试通过。
- Vue：详情页和报告页 loading、empty、error、正常展示测试或浏览器验证通过。
- 数据库：analysis_report 迁移验证通过。

## Completion Criteria

- 用户可以查看单股详情。
- 用户可以查看每日复盘。
- 报告明确展示数据时间、来源、风险点和人工确认边界。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/internal/reports/`, `frontend/src/reports.js`.
- Migration files: `backend/migrations/000006_analysis_reports.sql`.
- Test files: `backend/internal/reports/reports_test.go`, `frontend/tests/reports.test.js`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`; `npm test`; `npm run typecheck`.
- Test result: passed.
- Remaining risks: full Vue report pages and scheduled report UI integration are deferred to `08 Frontend Console`; current plan provides tested report domain and formatting.
