# Stock Detail Spec

## 1. Background

Users need a single view of price snapshots, alerts, and analysis context for a monitored stock.

## 2. Goals

- Aggregate snapshots and alerts by symbol.
- Preserve data time and source.

## 3. Non-Goals

- Do not make trading decisions.

## 4. Functional Scope

### Must Have

- Stock detail DTO.
- Summary fields for latest price, latest alert, and risk level.

## 5. Testing

- Aggregation with data.
- Empty data behavior.

## 6. Acceptance Criteria

- Given snapshots and alerts When detail is requested Then latest data time is shown.

## 7. Definition Of Done

- Go tests pass.
