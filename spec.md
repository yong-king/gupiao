# Spec Standard

## Purpose

本文定义本项目功能规格文档的写作规范。所有模块 spec 都应基于 `agent.md` 的全局约束编写，并保持可实现、可测试、可审计。

本项目是基于 Go 后端、Python Agent、Vue 前端的股票监控、分析与提醒系统。系统不自动买入、不自动卖出、不自动下单，只提供提醒、观察、风险提示和人工确认问题。

## Spec Rules

编写 spec 时必须遵守：

- 使用中文。
- 优先描述 MVP 可落地方案。
- 所有投资相关表达必须使用“提醒”“观察”“风险提示”，不得写成确定性交易指令。
- 明确哪些行为需要人工确认。
- 明确数据来源、数据时间、刷新限制和审计要求。
- 明确 Go、Python、Vue 的职责边界。
- Agent 输出必须是结构化 schema。
- 数据不足时必须定义失败或降级行为。
- 不编造不存在的数据源、券商能力或账户能力。
- 每条验收标准必须可验证。
- 每个模块必须定义测试门禁，代码生成或修改后必须测试通过才算完成。

## Spec Template

复制下面模板生成具体模块规格。

```markdown
# {{MODULE_NAME}} Spec

## 1. Background

说明模块背景、业务价值和它在系统中的位置。

必须回答：

- 这个模块解决什么问题。
- 服务哪些用户场景。
- 和股票池、持仓、提醒、刷新任务或 Agent 分析之间的关系。
- 是否涉及账户数据或外部数据源。

## 2. Goals

列出本模块必须实现的目标。

- Goal 1
- Goal 2
- Goal 3

## 3. Non-Goals

列出本阶段明确不做的内容。

- Non-goal 1
- Non-goal 2
- 不实现自动买入、自动卖出或自动下单。

## 4. User Scenarios

用用户视角描述核心场景。

- 用户希望……
- 用户点击……
- 用户收到……
- 用户看到……

## 5. Functional Scope

### Must Have

- 必须实现的功能点。

### Should Have

- 有价值但可以稍后完成的功能点。

### Out Of Scope

- 本模块不处理的功能点。

## 6. User Flow

### Normal Flow

```text
用户动作
-> 前端请求
-> Go 后端处理
-> 必要时调用 Python Agent
-> 数据落库
-> 前端展示或发送通知
```

### Error Flow

描述数据源失败、限频命中、Agent 失败、数据不足、重复提醒等情况。

### Manual Refresh Flow

如果模块涉及刷新，必须描述手动刷新流程。

### Scheduled Refresh Flow

如果模块涉及自动刷新，必须描述定时刷新流程和冷却规则。

## 7. Backend Design

### Go Responsibilities

- API。
- Service。
- Worker。
- Scheduler。
- Repository。
- Rate limiter。
- Audit log。

### API

```http
POST /api/...
GET /api/...
PATCH /api/...
DELETE /api/...
```

### Request Schema

```json
{}
```

### Response Schema

```json
{}
```

### Error Types

- validation_error
- not_found
- rate_limited
- data_source_error
- agent_error
- insufficient_data
- conflict

### Rate Limit And Retry

说明：

- 请求超时。
- 重试次数。
- 退避策略。
- 同账户刷新冷却。
- 同数据源限频。

### Audit

说明需要记录：

- 操作人。
- 请求参数摘要。
- 数据来源。
- 数据时间。
- 触发规则。
- 输出结果。
- 错误信息。

## 8. Python Agent Design

如果模块不需要 Python Agent，明确写“不涉及 Python Agent”。

### Agent Input Schema

```json
{}
```

### Agent Output Schema

```json
{
  "signal": "risk_warning",
  "confidence": 0.0,
  "risk_level": "low",
  "triggered_rules": [],
  "summary": "",
  "reasoning": [],
  "data_time": "",
  "source_refs": [],
  "missing_data": [],
  "recommended_action": "继续观察"
}
```

### Analysis Steps

- Step 1
- Step 2
- Step 3

### Deterministic Rules

列出由规则和结构化数据直接判断的内容。

### LLM Tasks

列出 LLM 只负责的摘要、解释、归纳或文本理解内容。

### Insufficient Data Behavior

定义数据不足时的返回结构，不允许强行输出买入或卖出观察。

## 9. Database Design

### Tables

#### table_name

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| id | uuid | yes | 主键 |
| created_at | timestamptz | yes | 创建时间 |
| updated_at | timestamptz | yes | 更新时间 |

### Indexes

- index_name(fields)

### Status Fields

列出状态枚举和状态流转。

### Retention

说明行情、任务、提醒、Agent 运行记录是否需要保留策略。

## 10. Frontend Design

### Pages

- 页面 1
- 页面 2

### Components

- 组件 1
- 组件 2

### User Actions

- 创建。
- 编辑。
- 删除。
- 手动刷新。
- 查看详情。
- 确认已读。

### Display Fields

- 数据时间。
- 触发原因。
- 风险等级。
- 规则来源。
- 刷新状态。
- 错误信息。

### States

- Loading。
- Empty。
- Error。
- Rate limited。
- Agent failed。
- Insufficient data。

## 11. Notifications And Alerts

### Trigger Conditions

列出触发提醒的条件。

### Deduplication

定义去重维度，例如：

- user_id。
- symbol。
- rule_id。
- signal。
- cooldown window。

### Cooldown

定义冷却时间和重复推送策略。

### Message Schema

```json
{
  "title": "",
  "summary": "",
  "signal": "",
  "risk_level": "",
  "symbol": "",
  "data_time": "",
  "source_refs": []
}
```

### Channels

- 前端提醒中心。
- WebSocket 或 SSE。
- 邮件。
- 企业微信或其他后续渠道。

## 12. Security And Risk Control

必须说明：

- 是否涉及账户数据。
- 是否保存敏感信息。
- 如何避免频繁刷新账户。
- 如何处理只读权限。
- 如何避免自动交易能力。
- 如何记录审计日志。

## 13. Testing

### Unit Tests

- 测试点 1
- 测试点 2

### Integration Tests

- 测试点 1
- 测试点 2

### Agent Schema Tests

- 输入 schema 校验。
- 输出 schema 校验。
- 数据不足时输出校验。

### End-To-End Tests

- 用户创建配置。
- 用户手动刷新。
- 系统触发提醒。
- 前端展示提醒。

### Edge Cases

- 数据源超时。
- 行情数据延迟。
- 重复提醒。
- 刷新冷却命中。
- Agent 返回无效 JSON。

## 14. Acceptance Criteria

使用 Given/When/Then 编写，每条必须可验证。

- Given 用户已创建股票池 When 用户点击手动刷新 Then 系统创建 refresh_job 并展示刷新状态。
- Given 数据源返回超时 When 刷新任务执行 Then 系统记录失败原因并不生成无依据提醒。
- Given Agent 返回 missing_data When 后端保存结果 Then 前端显示数据不足而不是买入或卖出观察。

## 15. Definition Of Done

模块只有满足以下条件才算完成：

- Must Have 功能已实现。
- Go、Python、Vue 相关测试已按影响范围新增或更新。
- 相关测试命令已执行并通过。
- Agent 输入输出 schema 已验证。
- 数据不足、数据源失败、限频命中和重复提醒等边界场景已验证。
- 前端涉及页面或交互时，已完成浏览器验证或等价端到端验证。
- 交付说明包含已执行测试、测试结果和剩余风险。
- 如果任何测试无法运行，必须记录原因，模块不得标记为完成。

## 16. Open Questions

列出仍需确认的问题，不要在 spec 中编造答案。

- 问题 1
- 问题 2
```

## Example Spec Prompt

```text
请按照 spec.md 的规范，为“提醒规则模块”生成完整 spec。

模块目标：
用户可以为股票池、单只股票或持仓配置价格、涨跌幅、成交量、止盈、止损、技术指标和新闻公告相关提醒。系统在手动刷新或定时刷新后判断是否触发提醒，并生成可追溯的 alert_event。

用户场景：
- 用户为 AAPL 设置价格低于 160 时提醒买入观察。
- 用户为持仓股票设置浮亏超过 8% 时提醒止损观察。
- 用户希望同一个提醒 30 分钟内不要重复推送。
- 用户希望看到提醒触发时的数据时间和触发原因。

输入：
- 股票代码。
- 市场。
- 规则类型。
- 阈值。
- 启用状态。
- 冷却时间。
- 适用范围：单只股票、股票池、持仓。

输出：
- alert_rule。
- alert_event。
- 前端提醒。
- 通知消息。

约束：
- 不允许生成自动买卖指令。
- 规则判断必须可追溯。
- 同一规则需要支持冷却和去重。
- 数据不足时不得触发确定性提醒。
```

## Implementation Prompt

```text
请阅读 agent.md 和 spec.md，并根据指定 spec 在当前仓库中实现功能。

实现要求：
- 先阅读现有代码结构和约定。
- 遵循项目已有目录、命名、错误处理和测试风格。
- 保持改动范围聚焦。
- 为核心逻辑添加测试。
- 每次代码生成或修改后必须运行相关测试，测试达到 spec 的验收标准后才算完成。
- Go 后端负责业务 API、规则、任务、存储、通知和审计。
- Python Agent 负责分析、指标、文本摘要和结构化分析输出。
- Vue 前端负责操作台、配置页面、状态展示和提醒中心。
- 外部数据源访问必须有超时、错误处理和限频入口。
- Python Agent 输出必须是结构化 schema。
- 前端必须展示数据时间、触发原因和风险等级。
- 不实现自动买入、自动卖出或自动下单。

本次要实现的 spec：
{{SPEC}}

请完成：
1. 代码实现。
2. 数据库迁移或模型更新，如需要。
3. 单元测试或集成测试。
4. 执行相关测试并确认通过。
5. 简短说明改动点、测试命令、测试结果和剩余风险。
```
