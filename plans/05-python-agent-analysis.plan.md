# 05 Python Agent Analysis Plan

## Status

completed

## Objective

建立 Python Agent 分析能力，用结构化 JSON 返回技术指标、文本摘要、风险点和提醒解释。

## Scope

- FastAPI Agent 服务。
- Agent 输入输出 schema。
- 健康检查。
- 基础技术指标计算。
- 新闻公告摘要和风险点提取。
- 基于 Go 规则触发结果生成解释。
- Go 调用 Python Agent。
- agent_runs 保存输入、输出、版本和错误。

## Out Of Scope

- 不让 Agent 自动交易。
- 不让 LLM 作为事实来源。
- 不实现复杂多 Agent 编排。
- 不实现完整回测平台。

## Dependencies

- [03 Market Refresh](</Users/youngking/Documents/jijin/plans/03-market-refresh.plan.md>)
- [04 Alerts And Notifications](</Users/youngking/Documents/jijin/plans/04-alerts-notifications.plan.md>)

## Required Specs

- Python Agent API Spec
- Technical Indicator Spec
- News Filing Analysis Spec
- Agent Signal Explanation Spec
- Agent Runs Spec

## Tasks

- [x] 生成并确认 Python Agent 相关 spec。
- [x] 实现 Agent 服务健康检查。
- [x] 定义并实现输入输出 schema。
- [x] 实现基础技术指标计算。
- [x] 实现文本摘要和风险点提取接口。
- [x] 实现 signal explanation 输出。
- [x] 实现 Go 调用 Agent 的 client。
- [x] 实现 agent_runs 记录。

## Testing Gate

- Python：`pytest` 通过。
- Python：输入 schema、输出 schema、missing_data、无效输入测试通过。
- Python：指标计算边界窗口和缺失数据测试通过。
- Go：Agent client 成功、超时、无效 JSON、agent_runs 保存测试通过。
- Vue：如果展示 Agent 输出，必须完成相关组件测试或浏览器验证。

## Completion Criteria

- Python Agent 输出稳定结构化 JSON。
- 数据不足时返回 `data_issue` 或 `missing_data`。
- Go 能调用 Agent 并保存 agent_runs。
- LLM 失败时有降级路径。
- 测试结果记录在 Delivery Notes。

## Delivery Notes

- Implementation files: `agent/src/agent_core/schema.py`, `agent/src/agent_core/indicators.py`, `agent/src/agent_core/text_analysis.py`, `agent/src/agent_core/analyzer.py`, `agent/src/agent_core/server.py`, `backend/internal/agentclient/`.
- Migration files: `backend/migrations/000005_agent_runs.sql`.
- Test files: `agent/tests/test_analyzer.py`, `backend/internal/agentclient/agentclient_test.go`.
- Test commands: `PYTHONPATH=src python3 -m unittest discover -s tests`; `env GOCACHE=/Users/youngking/Documents/jijin/backend/.gocache go test ./...`.
- Test result: passed. Go Agent client tests use a mock RoundTripper because sandboxed tests cannot bind local ports.
- Remaining risks: FastAPI and real LLM integration are deferred; MVP Agent uses deterministic standard-library analysis and HTTP server.
