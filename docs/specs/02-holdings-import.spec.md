# Holdings Import Spec

## 1. Background

Users can monitor positions without connecting a broker by manually entering or importing holdings.

## 2. Goals

- Define holding records.
- Parse CSV holdings.
- Validate market, symbol, quantity, and cost basis.

## 3. Non-Goals

- Do not connect to broker accounts.
- Do not save trading passwords.
- Do not place trades.

## 4. Functional Scope

### Must Have

- Holding model.
- CSV parser for `market,symbol,quantity,cost_basis`.
- Validation errors with row numbers.
- In-memory repository.

## 5. Testing

- CSV parse success.
- Invalid row handling.
- Repository save and lookup.

## 6. Acceptance Criteria

- Given valid CSV When parsed Then holdings are returned with normalized symbols.
- Given an invalid quantity When parsed Then the row error includes the row number.

## 7. Definition Of Done

- Holdings parser and tests exist.
- No account login or trading capability exists.
