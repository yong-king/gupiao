# 21 Portfolio Workbench UX Spec

## Purpose

Users need the console to behave like a real portfolio workbench instead of a demo page. The system must let each authenticated user create multiple stock pools, maintain holdings, configure watch-only alert rules, review alert logs, analyze holdings and custom stocks, configure read-only account integrations, and validate backend/LLM/database/Redis readiness from settings.

## Requirements

- Stock pools are user-scoped and support creating multiple named pools.
- A selected stock pool controls add/update stock actions, manual refresh, and alert checks.
- Stock pool stock entries support market, symbol, optional buy-watch price, and optional sell-watch price.
- Holdings are user-scoped and persisted in the backend database.
- Holdings support create, update, list, and delete by `(user_id, market, symbol)`.
- Holdings can be created from a stock pool selection or by entering a new stock code; saving a holding also keeps that stock available for monitoring.
- Demo-specific copy must not be used as the primary product flow.
- Alert rules target either holdings or stock-pool stocks by `(market, symbol)`.
- Alert rules support severity levels: low, medium, high, and critical.
- Alert center displays triggered rule logs and notifications with severity color, read status, and warning language.
- Refresh task UI is fully Chinese and lets the user refresh a selected stock pool or holding list; it must not hardcode `AAPL` in the primary action.
- Analysis report has two areas: holdings analysis and custom stock analysis.
- Holdings are automatically available in analysis. If a holding is deleted, the UI should make clear that the custom analysis target can still be kept or changed manually.
- Analysis detail shows price curve, daily change calendar, daily change records, and company/product summary for the selected stock.
- Account monitoring stores read-only integration configuration for external tools such as 同花顺、东方财富、雪球 or manual/CSV import. It must not store broker passwords, tokens, secrets, or implement automated trading.
- System settings shows backend/database/Redis/LLM/stock-source checks, configuration file paths, theme selection, and current user information.
- Theme selection supports light and dark modes and persists locally.
- All pages reachable from the sidebar must have visible, clickable controls that either perform an action or clearly report why the action is unavailable.

## Non-Goals

- No automatic buy/sell order placement.
- No storage of external broker trading credentials.
- No bypassing data-source or broker anti-abuse controls.
- No guaranteed prediction of future stock prices.
- MCP/web research is optional unless valid MCP credentials are configured; failure must be visible as a readiness/configuration issue, not a broken UI.

## Acceptance Criteria

- Given a logged-in user, when they create two stock pools, both pools are listed and the selected pool controls stock additions and refresh.
- Given a stock is in a pool, when the user enters quantity and cost, the holding can be saved, updated, and deleted.
- Given holdings exist, when the user opens reports, holdings appear as analysis targets and can load quote curve, daily change calendar, records, and profile summary.
- Given an alert rule with high or critical severity triggers, the alert center displays a high-emphasis warning row.
- Given the user opens refresh tasks, all primary labels are Chinese and no hardcoded `AAPL` demo action appears.
- Given the user opens settings, they can test backend dependencies and change light/dark theme.
- Backend Go tests and frontend unit/type tests pass after the implementation.
