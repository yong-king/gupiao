# Observability Spec

## 1. Background

The MVP needs consistent logs and simple counters before deployment.

## 2. Goals

- Provide structured log fields.
- Provide in-memory counters for jobs and errors.

## 3. Non-Goals

- Do not deploy external metrics stack.

## 4. Functional Scope

### Must Have

- Structured log event model.
- Counter registry.

## 5. Testing

- Log event fields.
- Counter increment.

## 6. Acceptance Criteria

- Given an operation occurs When logged Then service, action, status, and request id are present.

## 7. Definition Of Done

- Tests pass.
