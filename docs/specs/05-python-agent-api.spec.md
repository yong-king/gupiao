# Python Agent API Spec

## 1. Background

The Python Agent must return structured analysis outputs so Go can store, audit, and display results safely.

## 2. Goals

- Define input and output schemas.
- Provide health and analyze endpoints.
- Return `missing_data` for insufficient inputs.

## 3. Non-Goals

- Do not make LLM calls in this MVP step.
- Do not generate automatic trade instructions.

## 4. Functional Scope

### Must Have

- `AnalyzeRequest`.
- `AnalyzeResult`.
- Standard-library HTTP endpoint for `/analyze`.

## 5. Testing

- Schema output tests.
- Missing data tests.
- HTTP endpoint test through function-level analysis.

## 6. Acceptance Criteria

- Given no prices When analyze runs Then signal is `data_issue`.
- Given enough prices When analyze runs Then indicators and explanation are returned.

## 7. Definition Of Done

- Python tests pass.
