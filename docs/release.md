# Release And Regression

## Local Regression

```bash
cd backend
env GOCACHE=$PWD/.gocache go test ./...

cd ../agent
PYTHONPATH=src python3 -m unittest discover -s tests

cd ../frontend
npm test
npm run typecheck
```

## Middleware

PostgreSQL and Redis are defined in `deploy/docker-compose.yml`.

```bash
cd deploy
docker compose up --build
```

One-command deployment from the repository root:

```bash
sh deploy/scripts/deploy.sh
```

## Safety Checks

- No automatic trading APIs.
- No trading password storage.
- Account monitoring remains read-only.
- Alerts use observation language only.
