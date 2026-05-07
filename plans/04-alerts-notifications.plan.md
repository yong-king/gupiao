# 04 Alerts And Notifications Plan

## Status

completed

## Objective

实现确定性提醒规则、提醒事件和通知闭环，让用户在刷新后收到可解释、可追溯的提醒。

## Scope

- alert_rule CRUD。
- 价格、涨跌幅、成交量、止盈、止损、浮盈浮亏规则。
- 规则启停、冷却时间、适用范围。
- 刷新后规则判断。
- alert_event 结构化保存。
- 提醒去重和冷却。
- 前端提醒中心。
- 前端实时推送或通知抽象。

## Out Of Scope

- 不实现自动买入或卖出。
- 不把提醒写成确定性交易指令。
- 不实现 LLM 解释。
- 不实现复杂策略回测。

## Dependencies

- [03 Market Refresh](</Users/youngking/Documents/jijin/plans/03-market-refresh.plan.md>)

## Required Specs

- Alert Rules Spec
- Alert Events Spec
- Notification Spec

## Tasks

- [x] 生成并确认提醒规则 spec。
- [x] 实现 alert_rule 模型和 API。
- [x] 实现确定性规则判断引擎。
- [x] 实现 alert_event 生成和查询。
- [x] 实现提醒去重和冷却。
- [x] 实现前端提醒中心。
- [x] 实现前端实时提醒或通知抽象。

## Testing Gate

- Go：规则校验、阈值边界、禁用规则、数据不足、重复提醒、冷却窗口测试通过。
- 数据库：alert_rule 和 alert_event 迁移验证通过。
- Vue：提醒列表、风险等级、触发原因、已读状态测试或浏览器验证通过。
- Python：本 plan 不涉及 Python Agent；无需 Python 测试。

## Completion Criteria

- 确定性规则可以从刷新数据触发提醒。
- 提醒保存结构化字段，不只保存文本。
- 提醒可查询、可去重、可冷却、可审计。
- 前端展示数据时间、触发原因和风险等级。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/internal/alerts/`, `backend/internal/notifications/`, `backend/internal/api/alert_handlers.go`, `backend/internal/api/refresh_handlers.go`.
- Migration files: `backend/migrations/000004_alerts_notifications.sql`.
- Test files: `backend/internal/alerts/alerts_test.go`, `backend/internal/notifications/notifications_test.go`, `backend/internal/api/alert_handlers_test.go`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`.
- Test result: passed.
- Remaining risks: external notification channels and WebSocket/SSE delivery are deferred; MVP uses a queryable in-memory notification center.
