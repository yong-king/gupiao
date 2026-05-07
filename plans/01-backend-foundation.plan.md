# 01 Backend Foundation Plan

## Status

completed

## Objective

建立 Go 后端基础能力，为后续股票池、刷新任务、提醒规则和审计日志提供稳定业务底座。

## Scope

- Go API 路由、中间件、错误响应和健康检查。
- PostgreSQL 连接与迁移管理。
- 基础 repository/service 分层。
- 用户配置和系统设置基础模型。
- 审计日志基础能力。

## Out Of Scope

- 不实现股票池业务。
- 不实现行情数据源。
- 不实现 Python Agent 调用。
- 不实现前端页面。

## Dependencies

- [00 Repository Foundation](</Users/youngking/Documents/jijin/plans/00-repository-foundation.plan.md>)

## Required Specs

- Go API Foundation Spec
- Database And Migration Spec
- User Settings Spec
- Audit Log Spec

## Tasks

- [x] 生成并确认后端基础 spec。
- [x] 实现健康检查 API。
- [x] 实现统一错误响应。
- [x] 接入数据库连接和迁移。
- [x] 实现用户设置基础模型。
- [x] 实现审计日志写入接口。

## Testing Gate

- `go test ./...` 必须通过。
- 健康检查 handler 测试必须通过。
- 数据库迁移测试或等价迁移验证必须通过。
- 审计日志 repository 测试必须通过。

## Completion Criteria

- 后端 API 基座可运行。
- 数据库迁移可重复执行。
- 审计日志可写入和查询。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/cmd/server/main.go`, `backend/internal/api/`, `backend/internal/database/`, `backend/internal/settings/`, `backend/internal/audit/`.
- Migration files: `backend/migrations/000001_foundation.sql`.
- Test files: `backend/internal/api/errors_test.go`, `backend/internal/api/router_test.go`, `backend/internal/database/migrations_test.go`, `backend/internal/settings/settings_test.go`, `backend/internal/audit/audit_test.go`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`.
- Test result: passed.
- Remaining risks: live PostgreSQL connection and third-party migration runner are deferred; this plan implements deterministic migration discovery and initial SQL only.
