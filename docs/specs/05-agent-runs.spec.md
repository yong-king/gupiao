# Agent Runs Spec

## 1. Background

Go must persist Agent request/response metadata for auditability.

## 2. Goals

- Define agent run model.
- Store input, output, status, version, prompt version, and error.
- Handle invalid JSON and timeouts.

## 3. Non-Goals

- Do not build full prompt management.

## 4. Functional Scope

### Must Have

- Go Agent client.
- In-memory run repository.
- Tests for success and invalid response handling.

## 5. Testing

- Agent client success.
- Invalid JSON.
- Run repository save.

## 6. Acceptance Criteria

- Given Agent returns JSON When Go client calls it Then result and run are saved.

## 7. Definition Of Done

- Go tests pass.
