# 22 Portfolio Analysis Research RAG Spec

## Purpose

The portfolio console needs a richer analysis workflow. Users should see all holdings and stock-pool stocks in report lists, inspect a single stock detail, manage stock pools and alert rules, and drive public information collection by holding attention level.

## Requirements

- Analysis report must show all holding stocks in a holding-analysis list.
- Analysis report must show stocks from a selected stock pool in a stock-pool-analysis list.
- Each list row shows latest/today change. Chinese market color convention applies: rising is red, falling is green.
- Clicking a stock row selects the stock and shows detail: price curve, daily change calendar, daily change records, personal holding data, company/product summary, and LLM-style advice.
- Stock pools support create, delete, selection, adding stock, copying/moving one stock to another pool, and merging another pool into the current pool.
- Alert rules support create and delete.
- Refresh task supports refreshing selected stock pool and all holdings.
- Backend default automatic refresh interval is 30 minutes; the UI must make this default visible.
- Alert center and account monitoring sidebar entries must remain clickable and must render usable screens.
- Holdings include attention level: high, medium, low.
- Attention level controls research cadence:
  - high: collect public product/company information every 4 hours.
  - medium: every 6 hours.
  - low: every 24 hours.
- A research collection endpoint summarizes stock product/company information, combines it with latest price movement, and saves a RAG document for each stock.
- RAG storage must retain user, market, symbol, content, source type, metadata, and created time. If a true vector extension is unavailable locally, store vector-ready documents in PostgreSQL with metadata marking embedding status.

## Non-Goals

- No automatic trading.
- No broker credential storage.
- No aggressive crawling or bypassing paywalls/CAPTCHAs/anti-bot controls.
- No guarantee that MCP crawling works without valid MCP/Firecrawl credentials.

## Acceptance Criteria

- Given holdings exist, analysis report lists all holding stocks with red/green change values.
- Given a stock pool is selected, analysis report lists that pool's stocks with red/green change values.
- Given a list stock is clicked, detail panels update to that stock.
- Given a holding is saved with attention level, the backend persists it and the frontend displays it.
- Given a user deletes an alert rule, it disappears from the rule list after refresh.
- Given a user deletes a stock pool, it disappears and another pool is selected or created.
- Given a stock is copied or moved between pools, the target pool contains the stock.
- Given research collection runs for a holding, the backend saves a RAG document with attention level, refresh cadence, price movement, and summary metadata.
- Backend tests, frontend tests, frontend typecheck, and local smoke tests pass.
