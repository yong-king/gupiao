# Portfolio Risk Spec

## 1. Background

Users need account-level observations such as concentration and unrealized gain/loss.

## 2. Goals

- Calculate position values.
- Calculate concentration.
- Emit risk observations only.

## 3. Non-Goals

- Do not generate sell orders.

## 4. Functional Scope

### Must Have

- Portfolio risk summary.
- High concentration warning.
- Missing price behavior.

## 5. Testing

- Concentration.
- Missing price.

## 6. Acceptance Criteria

- Given one holding is over threshold When risk is calculated Then risk warning is returned.

## 7. Definition Of Done

- Risk tests pass.
