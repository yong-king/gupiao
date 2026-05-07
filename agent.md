# Agent Instructions

## Project

本项目是一个股票监控、分析与提醒系统。系统面向用户自选股票、股票池、手动导入持仓或只读账户数据，进行低频、可控、可追溯的监控分析，并在满足规则或分析条件时提醒用户关注买入、卖出、止盈、止损、减仓或继续观察。

系统不是自动交易机器人。任何买入、卖出、减仓、加仓和调仓动作都必须由用户人工确认。

## Product Principles

- 系统只提供“提醒”“观察”“风险提示”“需要人工确认的问题”，不得输出确定性交易指令。
- 所有提醒都必须可解释、可追溯、可审计。
- LLM 不能作为事实来源，只能用于摘要、归纳、解释和文本理解。
- 确定性规则和结构化数据判断优先于 LLM 判断。
- 数据不足时必须明确返回数据不足，不得强行生成买入或卖出观察。
- 不承诺收益，不使用“稳赚”“必涨”“必跌”等表达。

## Technical Stack

后端业务系统使用 Go：

- API 服务。
- 用户、股票池、持仓、账户只读配置。
- 提醒规则、刷新任务、任务调度。
- 数据源限频、账户访问风控、审计日志。
- WebSocket 或 SSE 实时推送。
- 通知发送和任务状态管理。

Agent 与分析系统使用 Python：

- 股票分析。
- 技术指标计算。
- 财报和公告解析。
- 新闻摘要和情绪辅助分析。
- LLM 报告生成。
- 策略信号解释。
- 风险点提取。

前端使用 Vue：

- 股票池管理。
- 持仓监控。
- 提醒规则配置。
- 实时提醒中心。
- Agent 分析报告。
- 刷新任务状态。
- 系统设置。

推荐基础设施：

- PostgreSQL 作为主数据库。
- TimescaleDB 可作为行情时间序列扩展。
- Redis 用于缓存、限频、任务队列和短期状态。
- Go 任务队列优先考虑 Asynq 或同类方案。
- Python Agent 服务优先使用 FastAPI，必要时配合 Celery 或 RQ。

## Architecture Boundaries

Go 负责稳定业务系统：

- 业务 API。
- 规则 CRUD。
- 确定性规则判断。
- 调度和异步任务。
- 数据持久化。
- 限频、重试、退避。
- 通知和审计。

Python 负责分析和 Agent 能力：

- 指标计算。
- 文本材料理解。
- 风险摘要。
- 策略信号解释。
- 结构化分析结果生成。

Python Agent 不直接访问交易密码，不负责下单，不绕过券商或数据源限制，只接收 Go 后端提供的结构化输入并返回结构化输出。

## Prohibited Scope

不得实现：

- 自动买入。
- 自动卖出。
- 自动下单。
- 高频交易。
- 绕过券商、行情源或数据平台风控的采集。
- 保存交易密码。
- 伪造或补全不存在的数据源结论。

可以实现：

- 买入观察提醒。
- 卖出观察提醒。
- 止盈观察提醒。
- 止损观察提醒。
- 风险升高提醒。
- 异动原因分析。
- 数据缺失提醒。
- 数据异常提醒。
- 人工确认问题清单。

## Workflows

### Manual Refresh

```text
用户点击刷新
-> Vue 调用 Go API
-> Go 创建 refresh_job
-> Go 检查限频和刷新冷却
-> Go 拉取行情、持仓、新闻或公告
-> Go 执行确定性规则判断
-> 必要时调用 Python Agent
-> Go 保存 alert_event、agent_run 和 audit_log
-> Go 推送前端状态或通知
```

### Scheduled Monitoring

```text
Go Scheduler 定时触发
-> 加载启用的股票池、持仓和规则
-> 检查市场时间、数据源限频和账户冷却
-> 拉取必要数据
-> 先执行确定性规则
-> 必要时调用 Python Agent 解释信号
-> 生成结构化提醒
-> 推送通知
```

### Daily Review

```text
盘后任务触发
-> 汇总股票池、持仓、行情、公告和新闻
-> Python Agent 生成结构化日报
-> Go 存储 analysis_report
-> Vue 展示日报、风险点和关键提醒
```

## Refresh And Rate Limit Rules

系统必须支持：

- 手动模式：只有用户点击时刷新。
- 保守模式：账户数据尽量不自动刷新，股票代码低频监控。
- 标准模式：股票池定时刷新，账户数据低频刷新。

默认建议：

- 行情数据：1 到 15 分钟，取决于数据源限制和用户配置。
- 新闻和公告：10 到 30 分钟。
- 账户持仓：30 到 120 分钟，或仅手动刷新。
- 每日复盘：交易日盘后执行。

必须实现：

- 同一账户刷新冷却时间。
- 同一数据源请求限频。
- 连续失败自动退避。
- 手动刷新也经过最小冷却检查。
- 每次刷新记录请求时间、数据时间、数据来源和错误信息。

## Core Data Models

优先设计以下核心模型或等价结构：

- users
- watchlists
- watchlist_symbols
- broker_accounts
- holdings
- price_snapshots
- market_data_points
- alert_rules
- alert_events
- analysis_reports
- agent_runs
- refresh_jobs
- data_sources
- notification_channels
- audit_logs

`agent_runs` 必须记录：

- 输入数据快照。
- Agent 类型和版本。
- Prompt 版本。
- 分析结果。
- 触发规则。
- 置信度。
- 数据时间。
- 错误信息。

## Alert Output

提醒必须结构化保存，不允许只保存自然语言文本。

推荐结构：

```json
{
  "symbol": "AAPL",
  "market": "US",
  "signal": "sell_watch",
  "risk_level": "high",
  "confidence": 0.72,
  "triggered_rules": [
    "price_below_stop_loss",
    "negative_news_detected"
  ],
  "summary": "价格跌破用户设置的止损线，同时出现负面新闻。",
  "recommended_action": "观察是否减仓",
  "data_time": "2026-05-05T15:30:00+08:00",
  "source_refs": [
    {
      "type": "price",
      "source": "market_data_provider",
      "time": "2026-05-05T15:30:00+08:00"
    }
  ]
}
```

推荐 `signal` 枚举：

- buy_watch
- sell_watch
- hold_watch
- take_profit_watch
- stop_loss_watch
- risk_warning
- abnormal_movement
- data_issue

## Python Agent Output

Python Agent 必须返回 JSON 或等价结构化对象。自然语言摘要只能作为字段之一，不能作为唯一输出。

Agent 输出必须包含：

- signal
- confidence
- risk_level
- triggered_rules
- summary
- reasoning
- data_time
- source_refs
- missing_data
- recommended_action

Agent 必须区分：

- 确定性规则触发。
- 数据驱动的统计判断。
- LLM 基于文本材料的解释。
- 推测性观察。

如果数据不足，必须返回 `data_issue` 或 `missing_data`，不得强行给出买入或卖出观察。

## Frontend Principles

Vue 前端应构建实际操作台，不做营销型首页。

主要页面：

- 股票池页面。
- 股票详情页。
- 持仓监控页面。
- 提醒规则页面。
- 实时提醒中心。
- Agent 分析报告页面。
- 刷新任务和数据源状态页面。
- 系统设置页面。

前端必须清晰展示：

- 数据时间。
- 触发原因。
- 规则来源。
- 风险等级。
- 刷新状态。
- 数据源状态。
- 系统提醒和用户最终操作之间的边界。

## Coding Rules

- 遵循现有项目结构、命名、错误处理和测试风格。
- Go 和 Python 服务边界必须清晰。
- 业务规则优先放在 Go 或确定性规则引擎中。
- LLM 只做摘要、解释、归纳和文本理解。
- 外部数据访问必须有超时、错误处理、限频入口和审计记录。
- 不要把数据源调用、策略判断、通知发送混在一个函数里。
- 每个重要工作流都要有状态表或日志，便于排查。
- 为核心规则、Agent schema、刷新限频和提醒去重添加测试。
- 优先完成 MVP，不要过早引入复杂多 Agent 编排。

## Testing Gate

每次代码生成、修改或模块实现都必须经过测试门禁，测试达到预期后才算完成。

执行要求：

- 生成或修改 Go 代码后，必须运行相关 Go 单元测试；模块完成前必须运行覆盖该模块的集成测试或等价验证。
- 生成或修改 Python Agent 代码后，必须运行 Python 单元测试、schema 校验测试和数据不足场景测试。
- 生成或修改 Vue 前端代码后，必须运行类型检查、单元测试或组件测试；涉及页面交互时必须进行浏览器验证或等价端到端验证。
- 修改数据库迁移、任务调度、刷新限频、提醒去重、通知发送或 Agent 输出 schema 时，必须补充或更新对应测试。
- 如果测试无法运行，必须明确记录阻塞原因、缺失依赖和未验证风险；该模块不得标记为完成。
- 不允许只靠人工阅读代码就把模块标记为完成。

模块完成标准：

- 代码已实现 spec 中的 Must Have。
- 相关测试已新增或更新。
- 相关测试命令已执行并通过。
- 失败场景和边界条件已覆盖。
- 测试结果和剩余风险已在交付说明中写清楚。

## MVP Scope

第一阶段优先实现：

- 用户股票池。
- 股票代码监控。
- 手动刷新。
- 定时低频刷新。
- 基础提醒规则。
- 行情数据落库。
- Python 基础指标分析。
- Agent 生成结构化分析结果。
- 提醒中心。
- 每日复盘报告。
- 账户数据 CSV 或手动导入。

暂不实现：

- 自动下单。
- 高频交易。
- 复杂多账户券商登录。
- 完整投研知识图谱。
- 复杂回测平台。
