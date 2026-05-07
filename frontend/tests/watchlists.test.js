import test from "node:test";
import assert from "node:assert/strict";

import { validateHolding } from "../src/holdings.js";
import { validateWatchlistSymbol } from "../src/watchlists.js";

test("validates and normalizes watchlist symbols", () => {
  const result = validateWatchlistSymbol({ market: "us", symbol: " aapl " });

  assert.equal(result.valid, true);
  assert.deepEqual(result.value, { market: "US", symbol: "AAPL" });
});

test("rejects unsupported market", () => {
  const result = validateWatchlistSymbol({ market: "AUTO", symbol: "AAPL" });

  assert.equal(result.valid, false);
  assert.deepEqual(result.errors, ["unsupported_market"]);
});

test("validates holdings", () => {
  const result = validateHolding({
    market: "hk",
    symbol: "0700",
    quantity: "5",
    costBasis: "300",
  });

  assert.equal(result.valid, true);
  assert.equal(result.value.market, "HK");
  assert.equal(result.value.quantity, 5);
});

test("rejects invalid holdings", () => {
  const result = validateHolding({
    market: "US",
    symbol: "AAPL",
    quantity: "0",
    costBasis: "-1",
  });

  assert.equal(result.valid, false);
  assert.deepEqual(result.errors, ["quantity_invalid", "cost_basis_invalid"]);
});
