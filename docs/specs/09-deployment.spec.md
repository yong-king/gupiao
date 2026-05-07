# Deployment Spec

## 1. Background

The MVP should define local middleware and service startup assets.

## 2. Goals

- Provide Dockerfiles.
- Provide compose for PostgreSQL and Redis.
- Provide health check commands.

## 3. Non-Goals

- Do not require production high availability.

## 4. Functional Scope

### Must Have

- `deploy/docker-compose.yml`.
- `backend/Dockerfile`.
- `agent/Dockerfile`.
- `frontend/Dockerfile`.

## 5. Testing

- Static compose file exists.
- Docker availability check documented.

## 6. Acceptance Criteria

- Given Docker is available When compose starts Then PostgreSQL and Redis services are defined.

## 7. Definition Of Done

- Deployment files exist.
