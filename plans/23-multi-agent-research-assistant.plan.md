# 23 Multi-Agent Research Assistant Plan

Status: completed

## Objective

Implement an attention-level multi-agent research workflow and a RAG-backed stock chat assistant.

## Scope

- Add workflow job/step and assistant chat database tables.
- Store deterministic local vector embeddings in `rag_vectors`.
- Add backend workflow APIs for attention-level research runs and workflow listing.
- Add backend stock assistant chat API using holdings, stock pools, market snapshots, profiles, and RAG documents.
- Add frontend controls for high/medium/low AI workflow runs.
- Add frontend stock assistant page.
- Update spec and plan index.
- Run backend, frontend, and agent tests.
- Apply migrations and run API smoke.
- Commit and push.

## Out Of Scope

- Production scheduler daemon.
- Real external web crawler integration without valid MCP/Firecrawl credentials.
- Trading execution.
- External embedding provider integration.

## Dependencies

- [22 Portfolio Analysis Research RAG](</Users/youngking/Documents/jijin/plans/22-portfolio-analysis-research-rag.plan.md>)
- PostgreSQL and Redis local services.
- Optional DeepSeek key for future LLM replacement; this plan keeps deterministic fallback behavior.

## Required Specs

- `docs/specs/23-multi-agent-research-assistant.spec.md`

## Tasks

- [x] Add migration for workflow jobs, workflow steps, and assistant messages.
- [x] Extend persistence for workflow, RAG lookup, vector indexing, and chat history.
- [x] Add workflow run/list API.
- [x] Add assistant chat API.
- [x] Add frontend AI workflow controls by attention level.
- [x] Add frontend stock assistant page.
- [x] Add backend and frontend tests.
- [x] Run tests and local smoke.
- [x] Commit and push.

## Testing Gate

- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
- `cd frontend && npm test`
- `cd frontend && npm run typecheck`
- `sh deploy/scripts/apply-migrations.sh`
- Smoke API:
  - Save `CN:000821` high-attention holding.
  - Run high-attention workflow.
  - Verify workflow jobs/steps are returned.
  - Verify assistant chat returns an answer using RAG context.

## Completion Criteria

- Tests pass.
- Workflow and assistant APIs work against local backend.
- Frontend exposes attention-level workflow and assistant UI.
- Changes are committed and pushed.

## Delivery Notes

- Added `backend/migrations/000013_agent_workflows_assistant.sql` for workflow jobs, workflow steps, and assistant chat messages.
- Added deterministic local embeddings in `rag_vectors` using `deterministic-hash-v1`; this gives a local vector-ready database path without requiring pgvector or an external embedding provider.
- Added `POST /api/workflows/research/run` to run the multi-agent workflow by holding attention level.
- Added `GET /api/workflows` to list workflow jobs and steps.
- Added `POST /api/assistant/chat` for RAG-backed stock assistant answers.
- Added frontend attention-level workflow controls in analysis reports.
- Added frontend `股票助手` page for stock-code questions.
- Tests passed:
  - `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
  - `cd agent && PYTHONPATH=src python3 -m unittest discover -s tests`
  - `cd frontend && npm test`
  - `cd frontend && npm run typecheck`
- Migration applied with `sh deploy/scripts/apply-migrations.sh`.
- Local smoke passed:
  - Saved `CN:000821` high-attention holding.
  - Ran high-attention workflow.
  - Workflow completed with `5` steps and `1` RAG/vector document.
  - Assistant chat returned an answer using `1` RAG document.
