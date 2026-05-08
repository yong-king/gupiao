# Persistent Data MCP Research Premium UI Spec

## 1. Background

The current MVP used in-memory repositories for several business objects. Users lose data after backend restart or session reset. The console also needs internet research about stock companies/products and a more polished professional interface.

## 2. Goals

- Persist account users, sessions, watchlists, holdings, alert rules, alert events, notifications, broker account configs, and market snapshots in PostgreSQL.
- Preserve user data after logout/login and backend restart.
- Keep Redis for cache/rate-limit/session acceleration later, but PostgreSQL is the source of truth.
- Add an MCP-backed internet research workflow for stock company/product information.
- Store collected research text in `rag_documents` and summarize it for decision support.
- Upgrade frontend visual quality toward a premium operations console.

## 3. Non-Goals

- Do not implement automatic trading.
- Do not store broker passwords, tokens, secrets, or write-enabled credentials.
- Do not scrape sites that require login, captcha bypass, paywall bypass, or prohibited bot behavior.
- Do not treat MCP/research output as investment advice.

## 4. Functional Scope

### Persistence

- Auth registration/login uses PostgreSQL `auth_users` and `auth_sessions`.
- Watchlists and symbols use PostgreSQL `watchlists` and `watchlist_symbols`.
- Holdings use PostgreSQL `holdings`.
- Alert rules/events/notifications use PostgreSQL `alert_rules`, `alert_events`, and `notifications`.
- Read-only account config uses PostgreSQL `broker_accounts`.
- Market snapshots use PostgreSQL `price_snapshots`.

### MCP Research

- Research workflow accepts `(market, symbol, company name)`.
- MCP provider should use approved web research tools such as Firecrawl to collect public company/product/news pages.
- Research output is normalized into source URL, title, summary, content, timestamp, and risk notes.
- Research documents are stored in `rag_documents`.
- LLM analysis runs are tracked in `llm_analysis_runs` when DeepSeek is used.

### UI

- Use a polished operations-console visual style.
- Make loading, empty, error, and success states visible.
- Avoid decorative landing-page treatment; prioritize dense actionable views.

## 5. Testing

- Go tests for persistence-backed auth and business APIs.
- API smoke: create data, restart backend, login again, verify data remains.
- MCP research smoke for one stock, with source links captured.
- Frontend tests for critical display helpers and expired-token behavior.

## 6. Acceptance Criteria

- Given a user logs out and logs in again, watchlists, holdings, rules, alerts, reports, and account config remain visible.
- Given backend restarts, the user can log in and still see persisted data.
- Given stock research is requested, public source summaries are stored and available for later RAG analysis.
- Given the frontend is opened through the local server, the console presents a modern, professional interface.
