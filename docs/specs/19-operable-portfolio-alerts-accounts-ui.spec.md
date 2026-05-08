# Operable Portfolio Alerts Accounts UI Spec

## 1. Background

The console must be usable after a stock is added to the stock pool. Holdings, alert rules, alert center, analysis reports, and account monitoring need visible configuration and feedback instead of static placeholder panels.

## 2. Goals

- Holdings page can manually configure a holding with market, symbol, quantity, and cost basis.
- Alert rules page can configure price and change-percent rules for any monitored symbol.
- Alert center can run a monitoring check, list triggered alerts, show notifications, and mark notifications read.
- Analysis report shows quote numbers and a visible chart/value even when there is only one saved snapshot.
- Account monitoring page can save read-only account configuration and list configured accounts.
- Frontend styling should feel more modern, denser, and operationally useful.

## 3. Non-Goals

- Do not implement automatic trading.
- Do not store broker passwords, tokens, secrets, or write-enabled account credentials.
- Do not implement full broker synchronization in this plan.

## 4. Functional Scope

### Must Have

- Backend supports single holding upsert through `POST /api/holdings`.
- Backend supports read-only account config through `GET/POST /api/accounts`.
- Backend supports marking notifications read through `POST /api/notifications/read`.
- Frontend holdings page has a form and latest-price/PnL table.
- Frontend rules page has a rule creation form and rules table.
- Frontend alert center has buttons that run checks, refresh alerts, and mark notifications read.
- Frontend report page includes latest price/open/high/low/change numbers.

## 5. Testing

- Go tests for holding upsert, account repository/API, and notification read behavior.
- Frontend tests for expired token handling and report number/chart helpers.
- API smoke creates user, adds stock, saves holding/rule/account, runs monitoring, fetches report data.

## 6. Acceptance Criteria

- Given a stock exists in the stock pool, the user can add a matching holding from the Holdings page.
- Given a stock symbol and threshold, the user can save an alert rule from the Rules page.
- Given matching rules exist, running Alert Center check refreshes quotes and lists triggered notifications when thresholds match.
- Given at least one snapshot exists, Analysis Report shows numeric quote summary and visible chart output.
- Given account config is read-only and has no sensitive metadata, it can be saved and listed.
