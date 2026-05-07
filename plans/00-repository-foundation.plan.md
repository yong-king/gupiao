# 00 Repository Foundation Plan

## Status

completed

## Objective

建立项目仓库基础结构，使 Go 后端、Python Agent、Vue 前端后续可以独立开发、测试和运行。

## Scope

- 创建 `backend/`、`agent/`、`frontend/`、`deploy/`、`docs/` 目录约定。
- 定义本地开发环境变量模板。
- 定义基础 README 或开发启动说明。
- 建立三端最小健康检查或空测试入口。
- 定义公共契约位置，例如错误码、时间格式、市场代码和股票代码格式。

## Out Of Scope

- 不实现业务 API。
- 不实现行情采集。
- 不实现提醒规则。
- 不实现 Python 分析逻辑。
- 不实现前端业务页面。

## Dependencies

- `agent.md`
- `spec.md`

## Required Specs

- Repository Foundation Spec
- Local Development Environment Spec
- Shared Contract Spec

## Tasks

- [x] 生成并确认 Repository Foundation Spec。
- [x] 创建三端目录结构。
- [x] 初始化 Go、Python、Vue 的最小项目骨架。
- [x] 添加环境变量模板。
- [x] 添加最小健康检查或测试入口。
- [x] 记录开发启动方式。

## Testing Gate

- Go：能运行后端最小测试命令。
- Python：能运行 Agent 最小测试命令。
- Vue：能运行前端最小测试或类型检查命令。
- 如果某端尚未初始化，必须有明确的跳过原因，且本 plan 不得标记 completed，除非 spec 明确允许该端延后。

## Completion Criteria

- 三端目录和最小工程骨架存在。
- 基础测试命令可执行并通过。
- 后续 plan 能基于该结构继续实现。
- Delivery Notes 记录测试命令和结果。

## Delivery Notes

- Implementation files: `README.md`, `.env.example`, `.gitignore`, `docs/development.md`, `docs/contracts.md`, `backend/`, `agent/`, `frontend/`, `deploy/`.
- Test files: `backend/internal/health/health_test.go`, `backend/internal/contracts/contracts_test.go`, `agent/tests/test_health.py`, `frontend/tests/contracts.test.js`.
- Test commands: `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`; `PYTHONPATH=src python3 -m unittest discover -s tests`; `npm test`; `npm run typecheck`.
- Test result: all passed. Initial Go test without workspace-local `GOCACHE` failed due sandbox cache permission; rerun with local `GOCACHE` passed.
- Remaining risks: Vue runtime dependencies are intentionally deferred to the frontend implementation plan; Python FastAPI dependency is intentionally deferred to the Python Agent plan.
