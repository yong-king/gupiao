# 24 LangGraph Agent 聊天体验与界面优化计划

状态：已完成

## 目标

优化股票助手的真实回答体验，减少模板化复读；优化前端多个页面的小组件布局，让主要面板尽量纵向铺满；同步补齐中文 spec、中文 plan、适当代码注释，以及助手相关接口的 Swagger 文档。

## 范围

- 优化 Python Agent 股票助手提示词与兜底回答逻辑。
- 优化 Go 后端股票助手响应结构，返回模型与调用状态元信息。
- 优化前端股票助手消息展示、模型状态展示和页面提示。
- 调整前端页面网格布局，改为以纵向整行为主。
- 将本模块 spec 和 plan 改为中文。
- 补充助手接口 Swagger/OpenAPI 字段描述。
- 增加或更新相关测试。
- 完成后重启验证本地服务。

## 非目标

- 不重构整套前端框架。
- 不新增交易能力。
- 不一次性翻译仓库所有历史 plan。

## 依赖

- [23 多智能体研究助手计划](</Users/youngking/Documents/jijin/plans/23-multi-agent-research-assistant.plan.md>)
- [26 股票智能体多 Agent 规范](</Users/youngking/Documents/jijin/docs/specs/26-stock-intelligence-multi-agent.spec.md>)

## 对应规范

- [24 LangGraph Agent 聊天体验与界面优化规范](</Users/youngking/Documents/jijin/docs/specs/24-langgraph-agent-chat-ux.spec.md>)

## 任务

- [x] 将本模块 spec 改为中文并补充新体验要求。
- [x] 将本模块 plan 改为中文。
- [x] 优化 Python Agent 助手提示词与兜底回答。
- [x] 优化 Go 后端助手接口返回字段与兜底文案。
- [x] 优化前端助手聊天区展示模型/状态。
- [x] 调整页面网格布局，减少小组件横向并排。
- [x] 更新 Swagger/OpenAPI 助手接口契约。
- [x] 更新或补充测试。
- [x] 重启本地服务并做冒烟验证。
- [ ] 提交代码。

## 测试门槛

- `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test`
- `cd frontend && npm run typecheck`

## 完成标准

- 助手回答更像真实问答，而不是上下文复读。
- 页面布局以纵向主面板为主，阅读体验明显改善。
- 前端能看到模型名、提供方和调用状态。
- Swagger 与实际响应一致。
- spec 与 plan 使用中文。

## 完成说明

- 股票助手新增“模型问题直答”短路逻辑，问模型时不再输出股票分析模板。
- 股票分析回答优化为“先结论、再依据、再缺口、再建议”的结构。
- 助手接口响应补充 `provider` 和 `llm_status`，前端同步展示模型状态。
- 前端主页面网格改为默认单列纵向铺满，减少多个小组件并排造成的割裂感。
- `POST /api/assistant/chat` 与 `POST /api/assistant/chat/stream` 的 Swagger 契约已同步更新。
- 已完成测试与重启冒烟验证。
