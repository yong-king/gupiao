# Watchlist Holdings Link Spec

## 1. Background

Holdings should be usable as monitoring inputs alongside watchlists.

## 2. Goals

- Represent holdings using the same `(market, symbol)` identity convention as watchlists.
- Allow later refresh and alert modules to query monitored holdings.

## 3. Non-Goals

- Do not auto-create trading rules.
- Do not connect to broker APIs.

## 4. Functional Scope

### Must Have

- Shared symbol normalization.
- Holding records include market and symbol.
- Repository lookup by user id.

## 5. Testing

- Holding symbols are normalized the same way as watchlist symbols.
- Repository returns only a user's holdings.

## 6. Acceptance Criteria

- Given imported holdings When future refresh modules query by user Then the holdings are available as monitored symbols.

## 7. Definition Of Done

- Shared symbol behavior is tested.
