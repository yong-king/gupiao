# Plan Index

本文是项目 plan 索引，不承载具体实现任务。每个可执行 plan 必须放在 `plans/` 目录下，并且只负责一个可以独立完成、独立测试、独立验收的范围。

全局约束见 [agent.md](/Users/youngking/Documents/jijin/agent.md)，spec 写作规范见 [spec.md](/Users/youngking/Documents/jijin/spec.md)。

## Plan Convention

每个 plan 文件必须包含：

- Status：`pending`、`spec_ready`、`in_progress`、`blocked`、`completed` 之一。
- Objective：该 plan 唯一目标。
- Scope：该 plan 负责的内容。
- Out Of Scope：明确不负责的内容。
- Dependencies：前置 plan 或外部依赖。
- Required Specs：该 plan 需要先产出的 spec。
- Tasks：可执行任务清单。
- Testing Gate：必须运行的测试和验收方式。
- Completion Criteria：完成标准。
- Delivery Notes：交付时记录实现文件、测试命令、测试结果和剩余风险。

## Execution Order

1. [00 Repository Foundation](</Users/youngking/Documents/jijin/plans/00-repository-foundation.plan.md>)
2. [01 Backend Foundation](</Users/youngking/Documents/jijin/plans/01-backend-foundation.plan.md>)
3. [02 Watchlist And Holdings](</Users/youngking/Documents/jijin/plans/02-watchlist-holdings.plan.md>)
4. [03 Market Refresh](</Users/youngking/Documents/jijin/plans/03-market-refresh.plan.md>)
5. [04 Alerts And Notifications](</Users/youngking/Documents/jijin/plans/04-alerts-notifications.plan.md>)
6. [05 Python Agent Analysis](</Users/youngking/Documents/jijin/plans/05-python-agent-analysis.plan.md>)
7. [06 Reports And Review](</Users/youngking/Documents/jijin/plans/06-reports-review.plan.md>)
8. [07 Account Monitoring](</Users/youngking/Documents/jijin/plans/07-account-monitoring.plan.md>)
9. [08 Frontend Console](</Users/youngking/Documents/jijin/plans/08-frontend-console.plan.md>)
10. [09 Hardening Release](</Users/youngking/Documents/jijin/plans/09-hardening-release.plan.md>)
11. [10 Configuration Management](</Users/youngking/Documents/jijin/plans/10-configuration-management.plan.md>)
12. [11 Auth Registration Login](</Users/youngking/Documents/jijin/plans/11-auth-registration-login.plan.md>)
13. [12 Stock Source Collection](</Users/youngking/Documents/jijin/plans/12-stock-source-collection.plan.md>)
14. [13 One Click Deployment](</Users/youngking/Documents/jijin/plans/13-one-click-deployment.plan.md>)
15. [14 Frontend Navigation Usability](</Users/youngking/Documents/jijin/plans/14-frontend-navigation-usability.plan.md>)
16. [15 Real Market Analytics And RAG](</Users/youngking/Documents/jijin/plans/15-real-market-analytics-rag.plan.md>)
17. [16 Operable Console Persistence And DeepSeek Config](</Users/youngking/Documents/jijin/plans/16-operable-console-persistence-config.plan.md>)
18. [17 Local Data And Repository Workflow](</Users/youngking/Documents/jijin/plans/17-local-data-and-repository-workflow.plan.md>)
19. [18 Stock Entry Monitor Calendar](</Users/youngking/Documents/jijin/plans/18-stock-entry-monitor-calendar.plan.md>)
20. [19 Operable Portfolio Alerts Accounts UI](</Users/youngking/Documents/jijin/plans/19-operable-portfolio-alerts-accounts-ui.plan.md>)
21. [20 Persistent Data MCP Research Premium UI](</Users/youngking/Documents/jijin/plans/20-persistent-data-mcp-research-premium-ui.plan.md>)
22. [21 Portfolio Workbench UX](</Users/youngking/Documents/jijin/plans/21-portfolio-workbench-ux.plan.md>)
23. [22 Portfolio Analysis Research RAG](</Users/youngking/Documents/jijin/plans/22-portfolio-analysis-research-rag.plan.md>)
24. [23 Multi-Agent Research Assistant](</Users/youngking/Documents/jijin/plans/23-multi-agent-research-assistant.plan.md>)
25. [24 LangGraph Agent Chat UX](</Users/youngking/Documents/jijin/plans/24-langgraph-agent-chat-ux.plan.md>)
26. [25 Observability Cadence Report Detail](</Users/youngking/Documents/jijin/plans/25-observability-cadence-report-detail.plan.md>)

## Global Testing Rule

每次代码生成或修改后必须运行相关测试。测试未达到 plan 的 Testing Gate 时，该 plan 不得标记为 `completed`。
