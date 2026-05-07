# Security Hardening Spec

## 1. Background

The project must preserve non-trading boundaries and avoid sensitive credential storage.

## 2. Goals

- Detect forbidden trading API labels.
- Detect sensitive metadata keys.

## 3. Non-Goals

- Do not implement full RBAC.

## 4. Functional Scope

### Must Have

- Security scanner helpers.
- Tests for forbidden terms.

## 5. Testing

- Sensitive key detection.
- Trading action detection.

## 6. Acceptance Criteria

- Given a route name contains order submission When scanned Then it is rejected.

## 7. Definition Of Done

- Security tests pass.
