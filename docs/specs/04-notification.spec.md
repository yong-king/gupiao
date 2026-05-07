# Notification Spec

## 1. Background

The MVP needs a notification abstraction so frontend alerts and future channels can share the same message structure.

## 2. Goals

- Provide in-memory notification center.
- Convert alert events to messages.
- Support read/unread status.

## 3. Non-Goals

- Do not integrate email or chat providers yet.
- Do not implement WebSocket in this plan.

## 4. Functional Scope

### Must Have

- Notification message model.
- Publish and list operations.
- Mark read operation.

## 5. Testing

- Publish and list.
- Mark read.

## 6. Acceptance Criteria

- Given an alert event is created When notification is published Then frontend can list it.

## 7. Definition Of Done

- Notification tests pass.
