# 24 LangGraph Agent Chat UX Plan

Status: completed

## Objective

Replace the simulated multi-agent path with an inspectable Python LangGraph workflow, add model routing and multi-turn streaming stock assistant, and improve the frontend financial console layout.

## Scope

- Add LangGraph-compatible workflow code in the Python agent service.
- Add DeepSeek task-to-model routing in agent configuration.
- Add agent workflow and assistant endpoints to the Python service.
- Extend Go agent client and workflow handler to call Python agent first.
- Add assistant session history persistence and streaming API.
- Update frontend stock assistant to a chat box with streaming output and session context.
- Add stock-detail refresh cooldown for clicked stock entries.
- Refresh frontend CSS/layout for cleaner financial dashboard presentation.
- Add tests for agent workflow/chat, backend workflow/chat streaming, and frontend behavior.
- Run test gates, smoke test, commit, and push.

## Out Of Scope

- Production scheduler daemon.
- Real broker integration or trading execution.
- Bypassing website scraping protections.
- Full external vector database replacement beyond existing local vector-ready storage.

## Dependencies

- [23 Multi-Agent Research Assistant](</Users/youngking/Documents/jijin/plans/23-multi-agent-research-assistant.plan.md>)
- Local backend config with `agent_url`.
- Optional `DEEPSEEK_API_KEY` for runtime LLM calls; tests must pass without calling external LLMs.

## Required Specs

- `docs/specs/24-langgraph-agent-chat-ux.spec.md`

## Tasks

- [x] Write spec and plan.
- [x] Implement Python LangGraph/fallback workflow and chat endpoint.
- [x] Add DeepSeek model routing config.
- [x] Extend Go agent client and workflow handler integration.
- [x] Add assistant session persistence and streaming endpoint.
- [x] Rework frontend assistant chat, streaming, cooldown, and layout.
- [x] Add/adjust tests.
- [x] Run backend, agent, and frontend test gates.
- [x] Smoke test local workflow and streaming chat.
- [x] Commit and push.

## Testing Gate

- `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test`
- `cd frontend && npm run typecheck`
- Local smoke:
  - start Python agent service.
  - start Go backend.
  - run high-attention workflow for `CN:000821`.
  - verify workflow step metadata names agent engine/model.
  - call streaming assistant endpoint and verify chunked answer.

## Completion Criteria

- LangGraph workflow code exists and is used by the agent service.
- Backend uses the Python agent when reachable and falls back safely.
- Assistant supports persisted multi-turn chat and streaming frontend output.
- Clickable stocks load usable detail data with cooldown protection.
- Tests and smoke checks pass.
- Changes are committed and pushed.

## Delivery Notes

- Added `agent/src/agent_core/langgraph_workflow.py` with LangGraph graph construction and deterministic sequential fallback using the same node functions.
- Added task model routing: collection uses `deepseek-v4-flash`, synthesis/chat uses `deepseek-chat`, and risk review uses `deepseek-v4-pro`.
- Added Python agent endpoints `/workflow/research` and `/assistant/chat`.
- Backend workflow now calls the Python agent first and stores returned step/model/engine metadata; it keeps the previous Go fallback path.
- Added `session_id` to assistant message persistence and `POST /api/assistant/chat/stream` for SSE streaming.
- Frontend `股票助手` is now a chat panel with multi-turn session id, streaming answer updates, and refreshed finance-console layout.
- Stock detail clicks use a 5 minute frontend cooldown unless the user explicitly presses a refresh/action button.
- Migration applied with `sh deploy/scripts/apply-migrations.sh`.
- Tests passed:
  - `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
  - `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
  - `cd frontend && npm test`
  - `cd frontend && npm run typecheck`
- Local smoke passed:
  - started Python agent and Go backend.
  - saved `CN:000821` high-attention holding.
  - ran high-attention workflow.
  - verified `langgraph_agent_workflow` RAG document and agent engine/model metadata.
  - verified streaming assistant endpoint emits SSE chunks and final response.
