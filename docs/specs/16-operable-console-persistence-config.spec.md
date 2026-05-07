# Operable Console Persistence And DeepSeek Config Spec

## 1. Background

The console must show visible effects for every clickable MVP action. The project also needs safer DeepSeek API key configuration, local PostgreSQL/Redis setup, database tables for stock monitoring data, and a direct "test real stock data then save" workflow.

## 2. Goals

- Make watchlist, holdings, alert rules, refresh, alerts, reports, and settings pages visibly useful.
- Add a direct real quote collection endpoint that tests provider access and saves the snapshot before returning analysis data.
- Configure DeepSeek through environment variables, not literal API keys in files.
- Expose backend dependency status for PostgreSQL, Redis, LLM config, and stock source config.
- Add database tables for stock profiles, RAG documents, and LLM analysis runs.

## 3. Non-Goals

- Do not commit real API keys.
- Do not implement automatic buy/sell/order submission.
- Do not bypass captchas, platform logins, or anti-scraping controls.
- Do not complete full vector search RAG in this plan.

## 4. Functional Scope

### Must Have

- `DEEPSEEK_API_KEY` is the configured secret env var for backend and agent examples.
- Backend can list watchlists and holdings for the current user.
- Backend can collect one real quote by `(market, symbol)`, save it, and return snapshot, daily changes, and profile.
- Backend exposes system dependency/config status without leaking secrets.
- Frontend pages show specific records after each action.
- Frontend report page can collect and save real AAPL data before rendering curve and analysis.

### Should Have

- Settings page should show PostgreSQL, Redis, LLM provider/model, key env name, and stock source.
- Rule page should show created rules.
- Holdings page should show imported holdings rows.
- User-facing messages should say what changed.

## 5. Testing

- Go tests for list endpoints, collect-and-save endpoint, and dependency status.
- Frontend tests for display helpers.
- Python config tests confirm DeepSeek env configuration.
- Run the real provider smoke test against AAPL.

## 6. Acceptance Criteria

- Given the user clicks stock pool actions Then the page shows created watchlist and symbols.
- Given the user clicks holdings import Then the page shows imported rows.
- Given the user clicks alert rule creation Then the page shows the rule.
- Given the user clicks real quote collection Then the backend fetches and saves the quote, and the report page can show chart, daily changes, and profile.
- Given settings loads Then it shows database, Redis, LLM, and stock source configuration status without displaying secret values.

## 7. Definition Of Done

- Spec and plan are updated.
- Relevant Go, Python, and frontend tests pass.
- Local API smoke verifies real quote collection and saved data retrieval.
