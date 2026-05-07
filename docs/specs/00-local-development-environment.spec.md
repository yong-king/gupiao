# Local Development Environment Spec

## 1. Background

Developers need consistent local commands to validate each service boundary before business features are implemented.

## 2. Goals

- Provide `.env.example`.
- Provide development documentation.
- Define default ports and service names.

## 3. Non-Goals

- Do not require Docker in the foundation plan.
- Do not require external network downloads.

## 4. Functional Scope

### Must Have

- Environment variable template.
- Backend run and test commands.
- Agent run and test commands.
- Frontend test command.

## 5. Testing

- Commands in the development guide must match existing scripts or entry points.

## 6. Acceptance Criteria

- Given a developer opens `docs/development.md` When they follow the current test commands Then each command exists and can run.

## 7. Definition Of Done

- `.env.example` exists.
- `docs/development.md` exists.
- Foundation tests pass.
