# Real Market Analytics And RAG Spec

## 1. Background

The MVP needs to verify whether real stock data collection is feasible and expose collected data visually. It also needs to keep daily up/down records that can later be indexed for RAG analysis, and collect company/product information so price movement analysis has business context.

## 2. Goals

- Collect real read-only quote data for user-provided stock symbols.
- Display a price curve in the frontend.
- Record daily price movement from collected snapshots.
- Generate RAG-ready text for daily change records.
- Provide stock company/product context and a conservative monitoring analysis.
- Keep all outputs as research and alerting signals, not trading instructions.

## 3. Non-Goals

- Do not bypass captchas, logins, paywalls, or anti-scraping controls.
- Do not submit buy/sell orders or create automatic trading controls.
- Do not claim investment advice or guaranteed prediction.
- Do not require a vector database in this plan; RAG-ready records are enough for this step.

## 4. Functional Scope

### Must Have

- Configurable real market provider, initially Stooq CSV for public quote testing.
- Backend endpoint for collected snapshots.
- Backend endpoint for daily change records.
- Backend endpoint for company/product profile and monitoring analysis.
- Frontend report page with an SVG price curve.
- Frontend daily up/down list.
- Frontend company/product analysis panel.

### Should Have

- Snapshot records include source, data time, open, high, low, close, volume, and company name when available.
- Daily records include `RAGText` so future RAG ingestion can chunk/index the content.
- Known companies can be enriched locally while the provider supplies live name/quote data.

## 5. Testing

- Go unit tests for Stooq CSV parsing, daily change aggregation, and profile analysis.
- Go API tests for snapshots, daily changes, and profile endpoints.
- Frontend tests for chart path rendering and daily change formatting.
- Browser smoke test for the report page curve and analysis panel.

## 6. Acceptance Criteria

- Given AAPL is in the watchlist When manual refresh runs Then the backend fetches a real Stooq quote and saves a snapshot.
- Given snapshots exist When requesting daily changes Then the API returns per-day change records with RAG text.
- Given the frontend report page loads market analysis Then it displays a curve, daily change list, and company/product context.
- Given analysis is displayed Then it is framed as monitoring/research and not as an automatic trade recommendation.

## 7. Definition Of Done

- Spec and plan are updated.
- Backend and frontend tests pass.
- Local app can call the new real-data path.
