# Watchlist Spec

## 1. Background

Users need a monitored symbol set before refresh jobs and alert rules can run.

## 2. Goals

- Create watchlists.
- Add and remove symbols.
- Normalize symbol values.
- Reject unsupported market codes.

## 3. Non-Goals

- Do not fetch market data.
- Do not trigger alerts.
- Do not implement automatic trading.

## 4. Functional Scope

### Must Have

- Watchlist model.
- Symbol model with `(market, symbol)`.
- In-memory repository for MVP foundation.
- Duplicate symbol rejection within the same watchlist.
- Frontend validation helper.

## 5. Testing

- Go tests for create, add symbol, duplicate rejection, invalid market rejection.
- Frontend tests for symbol normalization and validation.

## 6. Acceptance Criteria

- Given a user creates a watchlist When it is saved Then it can be loaded by id.
- Given a symbol is added twice When the second add is attempted Then a conflict-style error is returned.
- Given an unsupported market When validation runs Then validation fails.

## 7. Definition Of Done

- Watchlist model and tests exist.
- Frontend validation tests pass.
- No automatic trading behavior exists.
