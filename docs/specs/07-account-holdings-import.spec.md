# Account Holdings Import Spec

## 1. Background

Account monitoring uses CSV or manual holdings import first.

## 2. Goals

- Reuse holdings import.
- Associate holdings with account alias.

## 3. Non-Goals

- Do not connect broker API in MVP.

## 4. Functional Scope

### Must Have

- Account holdings wrapper.
- Import validation.

## 5. Testing

- Import valid holdings.
- Invalid rows are preserved.

## 6. Acceptance Criteria

- Given CSV holdings When imported Then account risk can be calculated.

## 7. Definition Of Done

- Tests pass.
