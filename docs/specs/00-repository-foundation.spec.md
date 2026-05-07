# Repository Foundation Spec

## 1. Background

This module establishes the monorepo structure for the Go backend, Python Agent, and Vue frontend. It provides the minimum files needed for future modules to implement and test code consistently.

## 2. Goals

- Create the repository directory layout.
- Provide minimal runnable or testable entry points for backend, agent, and frontend.
- Document local development commands.

## 3. Non-Goals

- Do not implement business APIs.
- Do not implement market data refresh.
- Do not implement alerts.
- Do not implement automatic buy, sell, or order submission.

## 4. Functional Scope

### Must Have

- `backend/`, `agent/`, `frontend/`, `deploy/`, and `docs/` directories.
- Minimal backend health package and test.
- Minimal Python Agent health package and test.
- Minimal frontend package with runnable test script.
- Shared contract documentation.

### Out Of Scope

- Real database migrations.
- Real Vue application pages.
- Real Python LLM integration.

## 5. Testing

- Backend must pass `go test ./...`.
- Agent must pass `python3 -m unittest discover -s tests`.
- Frontend must pass `npm test`.

## 6. Acceptance Criteria

- Given a fresh checkout When backend tests run Then the backend health test passes.
- Given a fresh checkout When Agent tests run Then the Agent health test passes.
- Given a fresh checkout When frontend tests run Then the frontend contract test passes.

## 7. Definition Of Done

- Must Have scope is implemented.
- Relevant tests pass.
- Test commands and results are recorded in the plan delivery notes.
