# Agent Signal Explanation Spec

## 1. Background

The Agent explains rule-triggered signals using structured data and text, not as an independent trading authority.

## 2. Goals

- Return `signal`, `confidence`, `risk_level`, `summary`, `reasoning`, `source_refs`, and `missing_data`.
- Distinguish deterministic rule triggers from observations.

## 3. Non-Goals

- Do not output automatic orders.

## 4. Functional Scope

### Must Have

- Structured JSON result.
- Missing data result.
- Deterministic explanation from input fields.

## 5. Testing

- Valid output schema.
- Missing data output.

## 6. Acceptance Criteria

- Given triggered rules When analyze runs Then triggered rules are present in output.

## 7. Definition Of Done

- Schema tests pass.
