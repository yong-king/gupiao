# Stock Entry Monitor Calendar Spec

## 1. Background

Users need an obvious place to type a stock code, fetch real stock information, add multiple stocks to a personal stock pool, and attach buy/sell monitor prices. Reports also need a calendar view for daily price movement.

## 2. Goals

- The stock pool page has market, symbol, buy monitor price, and sell monitor price inputs.
- Users can fetch stock information before adding it to monitoring.
- Users can add or update multiple monitored stocks in one watchlist.
- Each monitored stock shows buy/sell monitor prices, latest quote, and current monitor status.
- Adding monitor prices creates corresponding alert rules.
- Reports include a daily change calendar.

## 3. Non-Goals

- Do not place trades or connect order APIs.
- Do not implement full portfolio P/L accounting in this plan.
- Do not bypass stock data source anti-scraping controls.

## 4. Functional Scope

### Must Have

- `watchlist.Symbol` stores `buy_price` and `sell_price`.
- Adding an existing symbol updates its monitor prices instead of failing.
- Frontend stock pool page accepts arbitrary market/symbol input.
- Frontend can collect and display stock quote/profile for the entered symbol.
- Frontend can add/update monitored symbol and create buy/sell alert rules.
- Report page renders daily change records as a calendar.

## 5. Testing

- Go watchlist tests cover duplicate update and monitor price validation.
- Go API tests verify monitor prices round-trip through watchlist API.
- Frontend tests cover monitor status text and calendar rendering.
- Smoke test fetches a real quote and adds a monitored stock.

## 6. Acceptance Criteria

- Given the user enters `US` and `AAPL`, clicking fetch shows current quote and company/profile text.
- Given buy/sell monitor prices are entered, clicking add/update shows the stock in the monitored table with those prices.
- Given the same stock is added again, the row is updated rather than duplicated.
- Given daily changes exist, the report page shows a calendar-style grid with up/down styling.
