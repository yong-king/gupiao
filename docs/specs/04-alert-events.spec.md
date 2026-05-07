# Alert Events Spec

## 1. Background

Triggered alerts must be stored structurally with source data and data time.

## 2. Goals

- Store alert event fields.
- Deduplicate within cooldown windows.
- Preserve triggered rule ids and source refs.

## 3. Non-Goals

- Do not send external notifications in this module.

## 4. Functional Scope

### Must Have

- Alert event model.
- Event repository.
- Cooldown/deduplication check.

## 5. Testing

- Event creation.
- Duplicate suppressed inside cooldown.
- Event allowed after cooldown.

## 6. Acceptance Criteria

- Given a rule was triggered recently When it triggers again inside cooldown Then no duplicate event is created.

## 7. Definition Of Done

- Event tests pass.
