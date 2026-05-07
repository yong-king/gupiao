# Database And Migration Spec

## 1. Background

Future modules need database schema changes to be tracked consistently. The foundation plan introduces migration file discovery and ordering without adding external database drivers.

## 2. Goals

- Provide a stable migrations directory.
- Provide migration discovery and deterministic ordering.
- Provide initial schema for users, settings, and audit logs.

## 3. Non-Goals

- Do not connect to a live PostgreSQL instance in this plan.
- Do not introduce third-party migration tooling yet.

## 4. Functional Scope

### Must Have

- `backend/migrations/`.
- Migration file loader.
- Initial SQL migration.

## 5. Testing

- Migration loader tests verify ordering.
- Migration loader tests reject invalid filenames.

## 6. Acceptance Criteria

- Given migration files exist When migrations are loaded Then they are sorted by version.
- Given an invalid migration filename When migrations are loaded Then an error is returned.

## 7. Definition Of Done

- Migration files and loader exist.
- Go tests pass.
