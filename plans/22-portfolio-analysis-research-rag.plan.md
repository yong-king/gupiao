# 22 Portfolio Analysis Research RAG Plan

Status: completed

## Objective

Complete the analysis/report, stock-pool management, alert-rule deletion, holding attention level, and research-to-RAG workflow requested after the portfolio workbench release.

## Scope

- Add holding attention level to backend model, migration, persistence, API, and UI.
- Add stock pool deletion API and UI.
- Add stock copy/move/merge UI using existing stock insert API plus delete where available.
- Add alert rule deletion API and UI.
- Add research summary/RAG endpoint and PostgreSQL persistence.
- Update analysis report to list all holdings and selected-pool stocks, then show detailed stock analysis.
- Update refresh UI with selected pool/holdings actions and default 30-minute auto refresh language.
- Fix alert center and account monitoring click/render path by ensuring their render functions do not throw and all bindings remain valid.
- Run tests and smoke checks, then commit.

## Out Of Scope

- Production crawler scheduling daemon.
- Real vector embedding generation without configured embedding provider/vector extension.
- Broker automation and auto trading.

## Dependencies

- [21 Portfolio Workbench UX](</Users/youngking/Documents/jijin/plans/21-portfolio-workbench-ux.plan.md>)
- Local PostgreSQL and Redis services.
- DeepSeek key for future LLM provider checks; deterministic local summary is acceptable when LLM/MCP credentials are unavailable.

## Required Specs

- `docs/specs/22-portfolio-analysis-research-rag.spec.md`

## Tasks

- [x] Backend holding attention level and migration.
- [x] Backend stock pool delete support.
- [x] Backend alert rule delete support.
- [x] Backend research summary and RAG document persistence endpoint.
- [x] Frontend holdings attention level field and display.
- [x] Frontend stock pool delete/copy/move/merge controls.
- [x] Frontend alert rule delete button.
- [x] Frontend analysis report lists and stock detail panel.
- [x] Frontend refresh 30-minute default copy.
- [x] API smoke checks and browser-readiness review for alerts/accounts/report clickability.
- [x] Test gate.

## Testing Gate

- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test`
- `cd frontend && npm run typecheck`
- Apply migrations locally.
- Smoke API: save `CN:000821` holding with high attention, collect market data, run research/RAG collection, delete alert rule, delete a stock pool.
- Browser smoke: alerts and accounts navigation render their pages after reload.

## Completion Criteria

- Tests pass.
- Report lists holdings and selected-pool stocks.
- Holding attention level persists.
- RAG summary endpoint returns and stores a document.
- Stock pool and rule deletion work.
- Commit is created.

## Delivery Notes

- Added `holdings.attention_level` and `rag_vectors` migration in `backend/migrations/000012_attention_research_rag.sql`.
- Added stock pool delete, stock delete-from-pool, alert rule delete, and research summary/RAG collection APIs.
- Added `POST /api/research/collect`, which summarizes company/product/profile data with latest price movement and stores vector-ready RAG records in PostgreSQL.
- Updated analysis report to show holding stocks and selected stock-pool stocks as lists with Chinese market colors: red for up, green for down.
- Added stock detail panels for curve, daily calendar, daily changes, personal holding summary, and research/LLM-style recommendation summary.
- Added holdings attention level UI and research collection button.
- Added stock pool delete/copy/move/merge controls and alert rule delete controls.
- Updated refresh copy to show the backend default 30-minute automatic refresh interval.
- Tests passed:
  - `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
  - `cd frontend && npm test`
  - `cd frontend && npm run typecheck`
- Local smoke passed:
  - Saved `CN:000821` holding with `high` attention level.
  - Collected market data.
  - Generated RAG document `rag-dd87eb85b184c666670f` with interval `4h0m0s`.
  - Created and deleted an alert rule; verified rule count returned to 0.
  - Created and deleted a stock pool; verified pool count returned to 1.
- Browser automation through the in-app browser Node runtime was unavailable in this session, so clickability was covered by code path review plus local API/page reload readiness. The user-facing frontend is served at `http://127.0.0.1:5173/index.html?v=23`.
