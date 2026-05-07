# Frontend App Shell Spec

## 1. Background

The MVP needs a usable operations console entry point before a full Vue build is introduced.

## 2. Goals

- Provide local app shell.
- Provide navigation model.
- Avoid marketing homepage.

## 3. Non-Goals

- Do not install runtime dependencies in this plan.
- Do not add automatic trading controls.

## 4. Functional Scope

### Must Have

- `frontend/index.html`.
- App state model.
- Navigation labels for watchlists, holdings, rules, refresh, alerts, reports, accounts, settings.

## 5. Testing

- Node tests for navigation and safety.

## 6. Acceptance Criteria

- Given the app opens Then operation views are visible.

## 7. Definition Of Done

- Frontend tests pass.
