# 26 Stock Intelligence Multi-Agent Spec

## 1. Background

系统需要从“单次 AI 总结”升级为完整的股票分析智能体链路。一个可用的股票分析闭环必须同时处理公司/产品公开信息、实时交易行情、K线/涨跌样本、汇总判断、风险分析和 RAG/向量写入。

本模块服务持仓、股票池和股票助手。它不自动买入、不自动卖出、不自动下单，只提供研究、提醒、观察和风险提示。

## 2. Goals

- 建立 5 个明确职责的多智能体节点。
- 提供股票信息抓取 Skill，支持后续通过 MCP 调用开源股票信息服务。
- 获取交易行情上下文，包括最新价、涨跌幅、成交量和 K线样本摘要。
- 将公司/产品信息与交易行情汇总后交给分析 Agent。
- 将最终结构化结果写入 `rag_documents` 和 `rag_vectors`，供股票助手 RAG 使用。
- 所有节点记录工作流步骤、模型选择和操作日志。

## 3. Non-Goals

- 不实现自动买入、自动卖出或自动下单。
- 不绕过第三方网站反爬、登录、验证码或付费限制。
- 不在本阶段强制安装某个 MCP 服务；MCP 通过配置启用。
- 不把外部 MCP 结果当成唯一真实来源，必须保留行情源失败和 MCP 不可用降级。

## 4. User Scenarios

- 用户对高关注持仓运行工作流，系统抓取股票信息、读取行情、汇总分析并写入 RAG。
- 用户点击某只股票后，详情页可以看到基于最新行情和历史 RAG 的分析依据。
- 用户询问股票助手时，助手可以读取该股票的 RAG 文档、行情样本和公开信息摘要。
- 运维用户进入日志监控，可以看到每个 Agent 调用了什么模型、使用了什么信息源、输出了什么摘要。

## 5. Functional Scope

### Must Have

- 5 个 Agent 节点：
  - `stock_info_collect`：股票信息抓取 Skill Agent。
  - `trade_market_collect`：交易行情/K线 Agent。
  - `information_summarize`：信息汇总 Agent。
  - `investment_analysis`：分析 Agent。
  - `rag_vector_write`：RAG/向量写入 Agent。
- `stock_info_collect` 必须通过 Skill 封装外部来源。
- MCP 推荐配置使用开源 `yfinance-mcp-server`，用于股票信息、新闻和历史行情工具扩展。
- DeepSeek 模型按节点分配：
  - 抓取和行情节点：flash。
  - 汇总和向量写入：chat。
  - 分析节点：pro。
- RAG 文档 metadata 必须包含 agent 列表、模型路由、MCP 服务名、MCP 仓库、schema 版本。
- 如果配置了 `DEEPSEEK_API_KEY`，对话助手、信息汇总 Agent 和分析 Agent 必须真实调用 DeepSeek API；未配置时必须在日志 metadata 中标记 `missing-api-key` 并降级。
- Docker Compose 本地部署必须显式加载仓库根目录 `.env`，保证 `DEEPSEEK_API_KEY` 能传递到 `backend` 和 `agent` 容器；不允许出现宿主机 `.env` 已填写但容器内环境变量为空的情况。

### Should Have

- MCP 启用后优先使用 MCP 工具，失败后降级到后端 profile 和已保存行情。
- 支持后续扩展到多个 MCP 服务和多个信息源。

### Out Of Scope

- 券商交易接口下单。
- 需要账号登录的第三方爬取。
- 强依赖外部向量数据库；本阶段继续使用本地 JSONB embedding。

## 6. User Flow

### Normal Flow

```text
用户运行关注等级工作流
-> Go 后端选出持仓/股票池目标
-> Go 后端准备最新行情和历史样本
-> Python LangGraph 调用股票信息 Skill
-> Python Agent 汇总行情、公开信息和分析结果
-> Go 后端保存 workflow steps、operation logs、RAG 文档和向量
-> 前端展示工作流记录和股票助手可用上下文
```

### Error Flow

- MCP 未启用：记录 warning，使用后端 profile 降级。
- MCP 调用失败：记录 warning，继续使用已保存 profile 和行情。
- 行情源失败：优先使用备用行情源；仍失败则使用最近一次保存行情。
- 数据不足：输出 `insufficient_data` 风格的风险提示，不生成确定性买卖建议。
- Agent 服务不可用：Go 后端使用本地 fallback workflow，仍写入日志和 RAG。

## 7. Backend Design

### Go Responsibilities

- 选择工作流目标。
- 获取真实行情和历史样本。
- 调用 Python Agent。
- 保存 RAG、向量、工作流步骤和操作日志。
- 提供查询 API。

### API

```http
POST /api/workflows/research/run
GET /api/workflows
GET /api/operation-logs
POST /api/assistant/chat
POST /api/assistant/chat/stream
```

### Error Types

- validation_error
- data_source_error
- agent_error
- storage_error
- insufficient_data

### Audit

每个节点必须记录：

- user_id、market、symbol。
- step_name、agent_name。
- model。
- input_summary、output_summary。
- workflow_job_id。
- agent_engine。
- MCP 服务名和信息源。

## 8. Python Agent Design

### Agent Input Schema

```json
{
  "user_id": "user-1",
  "job_id": "wf-1",
  "market": "CN",
  "symbol": "000821",
  "attention_level": "high",
  "interval": "1h",
  "profile": {},
  "latest_snapshot": {},
  "snapshots": [],
  "snapshots_count": 0
}
```

### Agent Output Schema

```json
{
  "engine": "langgraph",
  "market": "CN",
  "symbol": "000821",
  "content": "结构化研究总结文本",
  "metadata": {
    "rag_schema": "stock_intelligence_v2",
    "agents": "stock_info_collect,trade_market_collect,information_summarize,investment_analysis,rag_vector_write",
    "stock_research_mcp": "yfinance-mcp-server"
  },
  "steps": []
}
```

### Analysis Steps

- 股票信息抓取 Skill Agent：通过 MCP 配置或 profile 降级获取公司/产品/新闻摘要。
- 交易行情/K线 Agent：整理最新价、涨跌幅、成交量、K线样本摘要。
- 信息汇总 Agent：合并公开信息和行情信息。
- 分析 Agent：输出研究、提醒和风险提示，不输出下单指令。
- RAG/向量写入 Agent：生成可落库内容和 metadata。

### LLM Tasks

- 公开信息归纳。
- 多源信息摘要。
- 风险提示表达。
- RAG 文档文本生成。

### Insufficient Data Behavior

缺少行情、公开信息或 RAG 时，必须在输出中明确“数据不足”，并建议人工补充或稍后刷新。

## 9. Database Design

复用：

- `agent_workflow_jobs`
- `agent_workflow_steps`
- `rag_documents`
- `rag_vectors`
- `operation_logs`

新增 metadata 约定：

- `rag_schema=stock_intelligence_v2`
- `stock_research_mcp`
- `stock_research_mcp_repository`
- `model_stock_info_collect`
- `model_trade_market_collect`
- `model_information_summarize`
- `model_investment_analysis`

## 10. Configuration

Agent 配置文件：`agent/config/agent.example.json`

```json
{
  "mcp": {
    "stock_research": {
      "enabled": false,
      "provider": "yfinance",
      "name": "yfinance-mcp-server",
      "repository": "https://github.com/barvhaim/yfinance-mcp-server",
      "command": "uvx",
      "args": ["yfinance-mcp-server"],
      "tools": ["get_stock_info", "get_news", "get_history"]
    }
  }
}
```

本地 Docker 部署要求：

- 仓库根目录 `.env` 用于保存 `DEEPSEEK_API_KEY` 等敏感变量。
- `deploy/docker-compose.yml` 必须通过 `env_file` 或等效显式方式把 `.env` 注入 `backend` 与 `agent` 容器。
- 启动后，系统依赖检查或容器内环境检查应能确认 `DEEPSEEK_API_KEY` 非空。

## 11. Testing Gate

- `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test`
- `cd frontend && npm run typecheck`

## 12. Acceptance Criteria

- 工作流返回 5 个新节点名称。
- RAG metadata 包含 `stock_intelligence_v2`。
- MCP 配置能从 agent config 读取。
- MCP 未启用时，工作流仍通过 profile 和行情样本降级完成。
- 操作日志能看到每个 Agent 的模型和输出摘要。
- 当仓库根目录 `.env` 配置了 `DEEPSEEK_API_KEY` 且通过 Docker Compose 启动时，股票助手返回的 `llm_status` 必须为真实模型调用状态，例如 `deepseek-api`，而不是 `missing-api-key` 或 `local-fallback`。
