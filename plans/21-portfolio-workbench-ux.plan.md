# 21 Portfolio Workbench UX Plan

Status: completed

## Objective

Complete the user-facing portfolio workbench flow for multiple stock pools, holdings CRUD, alert severity logs, holdings/custom analysis, account integration settings, and modernized UI.

## Scope

- Add missing backend holding delete support.
- Update frontend stock pool flow to create and select multiple user-scoped pools.
- Update frontend holdings flow to support create, update, list, delete, and stock-pool selection.
- Update frontend alert rules with target selection and severity levels.
- Update refresh and alert center screens so they use selected pools/holdings and Chinese product copy.
- Update reports into holdings analysis and custom analysis with quote curve, daily calendar, records, and profile summary.
- Update account monitoring to model read-only external app integrations.
- Update settings with dependency checks, theme selection, and user information.
- Refresh visual styling to a more polished operational console.

## Out Of Scope

- Real broker login automation.
- Automatic trading.
- Storing broker secrets.
- Guaranteed MCP research unless valid MCP credentials are configured.

## Dependencies

- Plans 16-20 persistence, auth, real market collection, and dependency checks.
- PostgreSQL and Redis local Docker services.
- DeepSeek key supplied via local environment/config for LLM checks.

## Required Specs

- `docs/specs/21-portfolio-workbench-ux.spec.md`

## Tasks

- [x] Implement backend holding delete in memory repository, PostgreSQL store, and API handler.
- [x] Replace Demo/AAPL primary frontend flows with selected stock pool and selected stock controls.
- [x] Add multi-pool create/select/load UI.
- [x] Add holdings edit/delete/select-from-pool UI.
- [x] Add rule target and severity UI.
- [x] Add severity-colored alert center logs and notifications.
- [x] Add holdings/custom analysis report targets and stock detail panel.
- [x] Add read-only external app account configuration UI.
- [x] Add settings theme/user/dependency UI.
- [x] Modernize CSS while keeping responsive behavior.
- [x] Run backend tests, frontend tests, frontend typecheck, and a local smoke check.

## Testing Gate

- `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
- `cd frontend && npm test`
- `cd frontend && npm run typecheck`
- Local API smoke must create/select a pool, add stock `CN:000821`, save/update/delete a holding, and collect market/profile data when the configured stock source allows it.

## Completion Criteria

- All required tests pass.
- The frontend exposes clickable controls for all sidebar pages.
- User-scoped data remains backend-persistent after login/logout.
- No primary UI copy mentions Demo/AAPL.
- Git commit and push are attempted after completion.

## Delivery Notes

- Backend holding delete added in memory repository, PostgreSQL store, and `/api/holdings` DELETE handler.
- Frontend console now uses selected stock pool and selected stock instead of primary Demo/AAPL flows.
- Added multi-pool create/select, holdings save/update/delete, rule severity selection, colored alert logs, holdings/custom analysis, read-only external app account config, theme setting, and user info setting.
- Updated visual styling in `frontend/index.html` and bumped frontend module cache version to `v=22`.
- Tests passed:
  - `cd backend && env GOCACHE=$PWD/.gocache go test ./...`
  - `cd frontend && npm test`
  - `cd frontend && npm run typecheck`
- Local smoke passed:
  - Created a smoke user and stock pool.
  - Added `CN:000821` to the stock pool.
  - Saved and updated a holding.
  - Collected real quote data for `CN:000821`; received price `10.54` from the configured source.
  - Deleted the holding and verified the list returned to zero.
