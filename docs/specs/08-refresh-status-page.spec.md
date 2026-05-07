# Refresh Status Page Spec

## 1. Background

Users need to see refresh status and rate-limit states.

## 2. Goals

- Represent refresh status text.
- Include rate-limited and error states.

## 3. Non-Goals

- Do not implement live polling in this plan.

## 4. Functional Scope

### Must Have

- Refresh status formatter.

## 5. Testing

- Status formatting tests.

## 6. Acceptance Criteria

- Given status is rate_limited Then UI text explains cooldown.

## 7. Definition Of Done

- Tests pass.
