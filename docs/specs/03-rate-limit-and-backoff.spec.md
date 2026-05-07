# Rate Limit And Backoff Spec

## 1. Background

The system must avoid frequent polling of data sources or accounts.

## 2. Goals

- Enforce cooldown by scope.
- Track consecutive failures.
- Expose rate-limited status.

## 3. Non-Goals

- Do not implement distributed Redis rate limiting yet.
- Do not implement account login cooldown in this plan.

## 4. Functional Scope

### Must Have

- In-memory cooldown limiter.
- Failure tracker with exponential backoff duration.
- Tests for cooldown hit and recovery.

## 5. Testing

- First request allowed.
- Second request inside cooldown rejected.
- Request after cooldown allowed.

## 6. Acceptance Criteria

- Given a scope was refreshed recently When refresh is requested again Then the request is rejected with `rate_limited`.

## 7. Definition Of Done

- Rate limit tests pass.
