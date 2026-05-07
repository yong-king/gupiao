# Jijin Stock Monitoring Agent

Jijin is a stock monitoring, analysis, and alerting system. It monitors user-provided symbols, watchlists, manually imported holdings, or read-only account data and produces explainable alerts for human review.

The system does not place trades. It must not implement automatic buy, sell, or order submission features.

## Repository Layout

```text
backend/   Go API, jobs, rules, storage, notifications, and audit logs.
agent/     Python analysis and Agent service.
frontend/  Vue operations console.
deploy/    Local deployment assets.
docs/      Specs, contracts, and development notes.
plans/     Executable module plans.
```

## Quick Checks

Run the current foundation tests:

```bash
cd backend && env GOCACHE=$PWD/.gocache go test ./...
cd agent && PYTHONPATH=src python3 -m unittest discover -s tests
cd frontend && npm test
```

Full local regression:

```bash
cd backend && env GOCACHE=$PWD/.gocache go test ./...
cd ../agent && PYTHONPATH=src python3 -m unittest discover -s tests
cd ../frontend && npm test && npm run typecheck
cd .. && docker compose -f deploy/docker-compose.yml config
```

## Local App

The current MVP frontend can be opened directly from:

```text
/Users/youngking/Documents/jijin/frontend/index.html
```

Middleware and services are defined in:

```bash
cd deploy
docker compose up --build
```

One-command deployment:

```bash
sh deploy/scripts/deploy.sh
```

## Project Rules

- Global agent instructions: `agent.md`
- Spec standard: `spec.md`
- Plan index: `plan.md`
