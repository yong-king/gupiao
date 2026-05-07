# Refresh Job Spec

## 1. Background

Users need manual and scheduled refreshes that are traceable and rate limited.

## 2. Goals

- Create refresh jobs.
- Execute jobs against monitored symbols.
- Track status transitions.

## 3. Non-Goals

- Do not trigger alert rules in this plan.
- Do not call Python Agent in this plan.

## 4. Functional Scope

### Must Have

- `refresh_job` model.
- `pending`, `running`, `succeeded`, `failed`, `rate_limited` statuses.
- Manual refresh service.
- Minimal scheduled refresh runner.

## 5. Testing

- Job status success.
- Provider failure transitions to failed.
- Manual refresh API returns job status and snapshots.

## 6. Acceptance Criteria

- Given a watchlist When manual refresh is requested Then a refresh job is created and snapshots are saved.
- Given provider fails When refresh executes Then the job status is failed and error is recorded.

## 7. Definition Of Done

- Refresh tests pass.
