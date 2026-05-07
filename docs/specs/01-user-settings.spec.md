# User Settings Spec

## 1. Background

Users need safe default runtime preferences before account monitoring, refresh jobs, and notifications are implemented.

## 2. Goals

- Define refresh modes.
- Define notification preference fields.
- Provide validation and repository interface.

## 3. Non-Goals

- Do not implement authentication.
- Do not implement notification delivery.
- Do not connect to broker accounts.

## 4. Functional Scope

### Must Have

- User settings model.
- Refresh mode validation.
- In-memory repository for foundation tests.

## 5. Testing

- Default settings validation.
- Invalid refresh mode rejection.
- Repository save and lookup.

## 6. Acceptance Criteria

- Given a new user When default settings are created Then refresh mode is conservative.
- Given an invalid refresh mode When settings are validated Then validation fails.

## 7. Definition Of Done

- Settings model and tests exist.
- Go tests pass.
