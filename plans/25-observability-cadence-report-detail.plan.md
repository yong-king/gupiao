# 25 Observability Cadence Report Detail Plan

Status: completed

## Objective

Add AI/crawler operation logging, configurable attention-level cadences, and split the analysis report list from stock detail analysis.

## Scope

- Add operation log database table, persistence methods, backend API, and frontend page.
- Log market/crawler calls, workflow agent steps, RAG writes, and assistant model calls.
- Add configurable product research and realtime quote refresh cadence by attention level.
- Update research/workflow paths to use configured product cadence.
- Update frontend reports to show only today's change and buy/sell reminder status.
- Add stock detail view for charts, daily records, holding detail, profile/RAG summary, and workflow controls.
- Update agent workflow comments and metadata to make node/agent count explicit.
- Add tests and run gates.
- Commit and push.

## Out Of Scope

- Production background scheduler.
- Real broker trading integration.
- External crawler bypass logic.

## Dependencies

- [24 LangGraph Agent Chat UX](</Users/youngking/Documents/jijin/plans/24-langgraph-agent-chat-ux.plan.md>)
- PostgreSQL migration runner.

## Required Specs

- `docs/specs/25-observability-cadence-report-detail.spec.md`

## Tasks

- [x] Write spec and plan.
- [x] Add backend cadence config and docs.
- [x] Add operation log migration, persistence, API, and tests.
- [x] Instrument market collect, research workflow, RAG writes, and assistant chat.
- [x] Refactor report/detail frontend and add logs page.
- [x] Update comments and agent metadata/tests.
- [x] Run all test gates and local smoke.
- [x] Commit and push.

## Testing Gate

- `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test`
- `cd frontend && npm run typecheck`
- `sh deploy/scripts/apply-migrations.sh`
- Smoke:
  - collect `CN:000821`.
  - run high-attention workflow.
  - ask stock assistant.
  - verify operation logs contain crawler and AI model entries.

## Completion Criteria

- Operation logs are persisted and visible in the frontend.
- Configurable cadence appears in config examples and is used by backend research paths.
- Report page is list-only; stock detail page contains deeper analysis.
- Test gates and smoke pass.

## Delivery Notes

- Added `backend/migrations/000015_operation_logs.sql` and `operation_logs` persistence/API.
- Added `/api/operation-logs` and frontend `日志监控` page.
- Market quote collection, public/profile information collection, LangGraph workflow steps, RAG writes, and stock assistant calls now write operation logs.
- Added configurable `cadence.product_research` and `cadence.realtime_quote` to backend config.
- Updated default product research cadence to high `1h`, medium `2h`, low `4h`.
- Added realtime quote cadence defaults high `2m`, medium `5m`, low `10m`.
- Analysis report page now only shows latest price, today's change, and buy/sell reminder status.
- Clicking a stock code opens a stock detail page with chart, daily calendar, daily records, holding detail, product/profile summary, RAG workflow controls, and research actions.
- Added Chinese comments around LangGraph node/model routing and log-write failure isolation.
- Tests passed:
  - `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
  - `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
  - `cd frontend && npm test`
  - `cd frontend && npm run typecheck`
- Migration applied with `sh deploy/scripts/apply-migrations.sh` after starting Docker Postgres/Redis.
- Local smoke passed:
  - registered smoke user.
  - saved `CN:000821` high-attention holding.
  - ran market collect; provider returned an error in this environment, and the failure was still logged as `crawler_quote_collect`.
  - ran high-attention workflow.
  - asked stock assistant.
  - verified operation logs include `crawler_quote_collect`, `crawler_public_info_collect`, `ai_agent_step`, `rag_vector_write`, and `ai_assistant_chat`, with DeepSeek model names.
