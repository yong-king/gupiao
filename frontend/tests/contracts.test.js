import test from "node:test";
import assert from "node:assert/strict";

import { alertSignals, normalizeSymbol, supportedMarkets } from "../src/contracts.js";

test("supports initial market codes", () => {
  assert.deepEqual(supportedMarkets, ["US", "HK", "CN"]);
});

test("normalizes symbols to uppercase", () => {
  assert.equal(normalizeSymbol(" aapl "), "AAPL");
});

test("uses observation alert language", () => {
  assert.ok(alertSignals.includes("buy_watch"));
  assert.ok(alertSignals.includes("sell_watch"));
  assert.ok(!alertSignals.includes("buy_order"));
  assert.ok(!alertSignals.includes("sell_order"));
});
