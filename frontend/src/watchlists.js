import { normalizeSymbol, supportedMarkets } from "./contracts.js";

export function normalizeMarket(market) {
  return String(market || "").trim().toUpperCase();
}

export function validateWatchlistSymbol(input) {
  const market = normalizeMarket(input?.market);
  const symbol = normalizeSymbol(input?.symbol);
  const errors = [];

  if (!supportedMarkets.includes(market)) {
    errors.push("unsupported_market");
  }
  if (!symbol) {
    errors.push("symbol_required");
  }

  return {
    valid: errors.length === 0,
    errors,
    value: { market, symbol },
  };
}
