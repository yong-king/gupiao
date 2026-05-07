# Audit Log Spec

## 1. Background

The system must be auditable because refreshes, alerts, and Agent outputs influence investment monitoring decisions.

## 2. Goals

- Define an audit log entry model.
- Provide an append-only repository interface.
- Provide a testable in-memory implementation.

## 3. Non-Goals

- Do not implement a full audit query UI.
- Do not store sensitive secrets in audit logs.

## 4. Functional Scope

### Must Have

- Entry fields for actor, action, target, request id, source, data time, and metadata.
- Append and list operations.
- Metadata copy protection.

## 5. Testing

- Append and list tests.
- Metadata mutation safety test.

## 6. Acceptance Criteria

- Given an action occurs When it is audited Then the entry can be listed later.
- Given metadata is mutated after append When audit entries are listed Then stored audit metadata is unchanged.

## 7. Definition Of Done

- Audit model and tests exist.
- Go tests pass.
