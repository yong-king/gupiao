# Account Settings Spec

## 1. Background

Users may monitor account-level holdings without granting trading capability.

## 2. Goals

- Store account alias and refresh mode.
- Mark account capability as read-only.
- Reject sensitive credential fields.

## 3. Non-Goals

- Do not save trading passwords.
- Do not log in to brokers.
- Do not place orders.

## 4. Functional Scope

### Must Have

- Account config model.
- Refresh mode validation.
- Sensitive field rejection helper.

## 5. Testing

- Valid config.
- Invalid mode.
- Sensitive field rejection.

## 6. Acceptance Criteria

- Given account config includes password When validated Then validation fails.

## 7. Definition Of Done

- Account tests pass.
