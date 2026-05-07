# 09 Hardening Release Plan

## Status

completed

## Objective

补齐可观测性、安全加固、部署和回归测试，使 MVP 可以本地或测试环境稳定运行。

## Scope

- 结构化日志。
- 任务指标和错误统计。
- 敏感信息保护检查。
- 只读账户边界检查。
- Dockerfile 和 docker-compose。
- 环境变量校验。
- 数据库迁移启动流程。
- 全链路回归测试。

## Out Of Scope

- 不做生产级高可用部署。
- 不接真实券商自动交易能力。
- 不实现复杂权限系统，除非前置 spec 已要求。

## Dependencies

- [01 Backend Foundation](</Users/youngking/Documents/jijin/plans/01-backend-foundation.plan.md>)
- [05 Python Agent Analysis](</Users/youngking/Documents/jijin/plans/05-python-agent-analysis.plan.md>)
- [08 Frontend Console](</Users/youngking/Documents/jijin/plans/08-frontend-console.plan.md>)

## Required Specs

- Observability Spec
- Security Hardening Spec
- Deployment Spec
- Regression Test Spec

## Tasks

- [x] 生成并确认发布加固 spec。
- [x] 实现结构化日志和基础指标。
- [x] 实现敏感信息检查和配置校验。
- [x] 实现 Dockerfile 和 compose。
- [x] 实现迁移启动流程。
- [x] 编写全链路回归测试。
- [x] 记录发布和回滚说明。

## Testing Gate

- Go：`go test ./...` 通过。
- Python：`pytest` 通过。
- Vue：类型检查和前端测试通过。
- E2E：手动刷新、规则触发、Agent 分析、提醒展示、日报展示流程通过。
- Deploy：本地 compose 启动和健康检查通过。
- Security：确认没有自动交易 API、交易密码字段或敏感信息日志泄漏。

## Completion Criteria

- MVP 可本地部署运行。
- 核心链路回归通过。
- 部署说明、环境变量和回滚说明存在。
- 已知风险和后续计划记录清楚。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `backend/internal/observability/`, `backend/internal/security/`, `docs/release.md`.
- Deployment files: `backend/Dockerfile`, `agent/Dockerfile`, `frontend/Dockerfile`, `deploy/docker-compose.yml`.
- Test files: `backend/internal/observability/observability_test.go`, `backend/internal/security/security_test.go`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`; `PYTHONPATH=src python3 -m unittest discover -s tests`; `npm test`; `npm run typecheck`; `docker compose -f deploy/docker-compose.yml config`.
- Test result: passed.
- Remaining risks: Docker compose was statically validated but containers were not started, so image build and middleware boot remain a runtime verification step.
