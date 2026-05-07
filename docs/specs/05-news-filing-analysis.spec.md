# News Filing Analysis Spec

## 1. Background

News and filing text should be summarized as risk observations with source references.

## 2. Goals

- Extract simple risk markers.
- Preserve source refs.
- Avoid unsupported investment conclusions.

## 3. Non-Goals

- Do not claim sentiment accuracy.
- Do not browse external news sources.

## 4. Functional Scope

### Must Have

- Text risk keyword extraction.
- Source ref passthrough.

## 5. Testing

- Empty text.
- Text with risk markers.

## 6. Acceptance Criteria

- Given text includes loss or investigation When analyzed Then risk summary mentions those markers.

## 7. Definition Of Done

- Text analysis tests pass.
