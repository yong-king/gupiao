# Alert Rules Spec

## 1. Background

Users need deterministic observation rules that produce traceable alert events after refreshes.

## 2. Goals

- Support price and percentage threshold rules.
- Support enable/disable and cooldown fields.
- Keep rule output as observation language.

## 3. Non-Goals

- Do not place orders.
- Do not generate deterministic trading instructions.
- Do not call LLM in this plan.

## 4. Functional Scope

### Must Have

- Alert rule model.
- In-memory rule repository.
- Rule evaluator for price below, price above, change percent above, change percent below.

## 5. Testing

- Rule validation.
- Threshold matching.
- Disabled rule does not trigger.

## 6. Acceptance Criteria

- Given price is below a configured threshold When rules are evaluated Then a `buy_watch` or configured observation signal can be emitted.

## 7. Definition Of Done

- Rule tests pass.
