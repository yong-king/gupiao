# Market Data Source Spec

## 1. Background

Refresh jobs need a consistent provider interface before real market data integrations are added.

## 2. Goals

- Define a provider interface for quote snapshots.
- Provide a deterministic mock provider.
- Preserve source and data time.

## 3. Non-Goals

- Do not implement realtime market data.
- Do not call external APIs in this plan.

## 4. Functional Scope

### Must Have

- Quote request and snapshot models.
- Mock provider with configurable quotes and failures.
- Snapshot repository.

## 5. Testing

- Mock provider success and failure.
- Snapshot repository saves and lists by symbol.

## 6. Acceptance Criteria

- Given symbols When refresh runs Then snapshots are returned with data time and source.
- Given provider failure When refresh runs Then the job records failure without generating unsupported conclusions.

## 7. Definition Of Done

- Provider and repository tests pass.
