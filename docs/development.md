# Development Guide

## Prerequisites

- Go
- Python 3
- Node.js and npm

## Backend

```bash
cd backend
env GOCACHE=$PWD/.gocache go test ./...
env GOCACHE=$PWD/.gocache go run ./cmd/server
```

The backend defaults to `:8080`. Override with `GO_BACKEND_ADDR`. Use a workspace-local `GOCACHE` when running tests in sandboxed environments.

## Python Agent

```bash
cd agent
PYTHONPATH=src python3 -m unittest discover -s tests
PYTHONPATH=src python3 -m agent_core.server
```

The Agent defaults to `127.0.0.1:8090`. Override with `AGENT_HOST` and `AGENT_PORT`.

## Frontend

```bash
cd frontend
npm test
npm run typecheck
```

The frontend is intentionally minimal in the repository foundation plan. Vue runtime dependencies are introduced by the frontend implementation plan.

## Testing Rule

Every code generation or modification must run the relevant tests before a plan can be marked completed. If tests cannot run, record the blocker and leave the plan unfinished.
