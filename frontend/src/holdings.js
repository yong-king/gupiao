import { validateWatchlistSymbol } from "./watchlists.js";

export function validateHolding(input) {
  const symbolResult = validateWatchlistSymbol(input);
  const quantity = Number(input?.quantity);
  const costBasis = Number(input?.costBasis);
  const errors = [...symbolResult.errors];

  if (!Number.isFinite(quantity) || quantity <= 0) {
    errors.push("quantity_invalid");
  }
  if (!Number.isFinite(costBasis) || costBasis < 0) {
    errors.push("cost_basis_invalid");
  }

  return {
    valid: errors.length === 0,
    errors,
    value: {
      ...symbolResult.value,
      quantity,
      costBasis,
    },
  };
}
