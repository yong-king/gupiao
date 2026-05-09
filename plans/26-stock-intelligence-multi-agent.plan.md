# 26 Stock Intelligence Multi-Agent Plan

Status: completed

## Objective

Upgrade the existing research workflow into a complete stock intelligence multi-agent system with stock information skill, trading market context, summarization, analysis, and RAG/vector persistence.

## Scope

- Add a stock research Skill under the Python Agent.
- Add MCP configuration for an open-source stock MCP server.
- Rename and restructure LangGraph nodes around the complete stock intelligence flow.
- Pass historical snapshots from Go backend into Python Agent.
- Store workflow metadata with schema version and MCP provenance.
- Update specs and plan index.
- Run backend, agent, and frontend tests.
- Commit and push after tests pass.

## Out Of Scope

- Automatic trading.
- Broker order placement.
- Login-required crawling or anti-bot bypass.
- Mandatory runtime installation of the external MCP server.

## Dependencies

- [23 Multi-Agent Research Assistant](</Users/youngking/Documents/jijin/plans/23-multi-agent-research-assistant.plan.md>)
- [25 Observability Cadence Report Detail](</Users/youngking/Documents/jijin/plans/25-observability-cadence-report-detail.plan.md>)
- Optional open-source MCP candidate: `https://github.com/barvhaim/yfinance-mcp-server`

## Required Specs

- [26 Stock Intelligence Multi-Agent Spec](</Users/youngking/Documents/jijin/docs/specs/26-stock-intelligence-multi-agent.spec.md>)

## Tasks

- [x] Add `StockResearchSkill` abstraction.
- [x] Add agent MCP config for `yfinance-mcp-server`.
- [x] Update LangGraph workflow nodes and model routing.
- [x] Include snapshot/Kline context in workflow input.
- [x] Update agent tests for new node names and metadata.
- [x] Update spec and plan docs.
- [x] Run full test gate.
- [x] Commit and push.

## Testing Gate

- `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test`
- `cd frontend && npm run typecheck`

## Completion Criteria

- All tests pass.
- Workflow returns the 5 stock intelligence nodes.
- RAG metadata includes MCP provenance and `stock_intelligence_v2`.
- Backend compiles with expanded workflow request schema.
- Changes are committed and pushed.

## Delivery Notes

- Added `agent/src/agent_core/skills/stock_research.py` as the stock information Skill abstraction.
- Added MCP config for `yfinance-mcp-server` in `agent/config/agent.example.json`.
- Updated LangGraph workflow to the five-node stock intelligence chain:
  - `stock_info_collect`
  - `trade_market_collect`
  - `information_summarize`
  - `investment_analysis`
  - `rag_vector_write`
- Expanded backend workflow request payload with historical snapshots for K-line context.
- Added `rag_schema=stock_intelligence_v2` and MCP provenance metadata.
- Tests passed:
  - `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
  - `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
  - `cd frontend && npm test`
  - `cd frontend && npm run typecheck`
- Local workflow smoke returned all five new nodes and `stock_intelligence_v2`.
