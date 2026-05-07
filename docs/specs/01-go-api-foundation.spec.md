# Go API Foundation Spec

## 1. Background

The backend needs a stable API foundation before business modules add watchlists, refresh jobs, alerts, and Agent calls.

## 2. Goals

- Provide a central router.
- Provide a health endpoint.
- Provide a structured JSON error response.
- Keep API behavior testable with the Go standard library.

## 3. Non-Goals

- Do not implement authentication.
- Do not implement business APIs.
- Do not implement automatic buy, sell, or order submission.

## 4. Functional Scope

### Must Have

- `GET /healthz`.
- JSON error writer with shared error codes.
- Router constructor used by the server entry point.

## 5. Testing

- Handler tests verify health output.
- Error response tests verify status code, error code, and request id.

## 6. Acceptance Criteria

- Given the server starts When `/healthz` is requested Then it returns `{"status":"ok","service":"backend"}`.
- Given an API error is written When the response is decoded Then it includes `error.code`, `error.message`, and `error.request_id`.

## 7. Definition Of Done

- API foundation is implemented.
- Go tests pass.
