# Shared Contract Spec

## 1. Background

The project needs stable cross-service conventions before modules start exchanging API payloads and Agent outputs.

## 2. Goals

- Define time format.
- Define initial market codes.
- Define symbol normalization.
- Define error response shape.
- Define audit field expectations.
- Define allowed alert language.

## 3. Non-Goals

- Do not finalize all future API schemas.
- Do not define broker-specific account contracts.

## 4. Functional Scope

### Must Have

- Shared contract document under `docs/`.
- Backend contract constants for initial market codes and error codes.
- Python contract constants for health response.
- Frontend testable contract fixture.

## 5. Testing

- Backend tests verify market and error code constants.
- Python tests verify health payload shape.
- Frontend tests verify contract data.

## 6. Acceptance Criteria

- Given future modules need error codes When they read shared contracts Then initial codes are documented.
- Given future modules need alert copy When they read shared contracts Then automatic trading language is prohibited.

## 7. Definition Of Done

- Shared contracts are documented.
- Tests validate foundation contract constants.
