import test from "node:test";
import assert from "node:assert/strict";

import { buildSparklinePath, formatDailyChange, monitorText, renderChangeCalendar, renderPriceChart, summarizeProfile } from "../src/market.js";

test("formats daily change records for display and rag review", () => {
  const text = formatDailyChange({ date: "2026-05-06", close: 105, change: 5, change_percent: 5 });
  assert.match(text, /2026-05-06/);
  assert.match(text, /\+5.00/);
});

test("builds svg price chart path", () => {
  const path = buildSparklinePath([{ price: 100 }, { price: 105 }, { price: 103 }]);
  assert.match(path, /^M /);
  assert.match(path, / L /);
  assert.match(renderPriceChart([{ price: 100 }, { price: 105 }]), /svg/);
});

test("summarizes company profile", () => {
  assert.match(summarizeProfile({ name: "APPLE INC", sector: "Consumer Electronics", recommendation: "observe" }), /APPLE/);
});

test("formats monitor price status", () => {
  assert.match(monitorText({ price: 150 }, { buy_price: 160, sell_price: 230 }), /买入关注价/);
  assert.match(monitorText({ price: 240 }, { buy_price: 160, sell_price: 230 }), /卖出关注价/);
  assert.match(monitorText({ price: 200 }, { buy_price: 160, sell_price: 230 }), /观察中/);
});

test("renders daily change calendar", () => {
  const html = renderChangeCalendar([{ date: "2026-05-06", close: 105, change_percent: 5 }]);
  assert.match(html, /change-calendar/);
  assert.match(html, /up/);
});
