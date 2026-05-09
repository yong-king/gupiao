import test from "node:test";
import assert from "node:assert/strict";

import { buildSparklinePath, formatDailyChange, monitorText, renderCandlestickChart, renderChangeCalendar, renderPriceChart, renderRealtimeQuote, summarizeMarketNumbers, summarizeProfile } from "../src/market.js";

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
  assert.match(renderPriceChart([{ price: 100 }]), /chart-point/);
});

test("renders realtime quote board and kline candles", () => {
  const quote = { price: 10.71, open: 10.43, high: 10.78, low: 10.42, previous_close: 10.54, change_percent: 1.61, volume: 12381200, source: "tencent" };
  assert.match(renderRealtimeQuote(quote), /最新价/);
  assert.match(renderRealtimeQuote(quote), /tencent/);
  const html = renderCandlestickChart([
    { date: "2026-05-07", open: 10.31, close: 10.54, high: 10.67, low: 10.31, volume: 13836200 },
    { date: "2026-05-08", open: 10.43, close: 10.71, high: 10.78, low: 10.42, volume: 12381200 },
  ]);
  assert.match(html, /kline-chart/);
  assert.match(html, /candle/);
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

test("summarizes market numbers", () => {
  const text = summarizeMarketNumbers({ price: 10.54, open: 10.31, high: 10.67, low: 10.31 }, { change_percent: 2.4 });
  assert.match(text, /现价 10.54/);
  assert.match(text, /\+2.40%/);
});
