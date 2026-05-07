# Technical Indicator Spec

## 1. Background

The Agent needs deterministic calculations before natural language explanation.

## 2. Goals

- Compute simple moving average.
- Compute change percent.
- Compute volatility estimate.

## 3. Non-Goals

- Do not implement a full technical analysis library.

## 4. Functional Scope

### Must Have

- Moving average.
- Latest change percent.
- Missing data behavior.

## 5. Testing

- Indicator calculation.
- Empty price list.

## 6. Acceptance Criteria

- Given prices When indicators run Then SMA and latest change are deterministic.

## 7. Definition Of Done

- Indicator tests pass.
