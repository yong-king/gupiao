import test from "node:test";
import assert from "node:assert/strict";

import { formatAlertMessage, formatRefreshStatus, getViewCopy, hasTradingControls, isKnownView, layoutClassForAuthState, navItems, refreshModes } from "../src/app.js";

test("nav contains operations console views", () => {
  assert.ok(navItems.some((item) => item.id === "watchlists"));
  assert.ok(navItems.some((item) => item.id === "alerts"));
  assert.equal(hasTradingControls(navItems), false);
});

test("every nav item has usable view copy", () => {
  for (const item of navItems) {
    const view = getViewCopy(item.id);
    assert.equal(view.title, item.label);
    assert.ok(view.description.length > 0);
    assert.ok(view.empty.length > 0);
  }
  assert.equal(isKnownView("missing"), false);
  assert.equal(getViewCopy("missing").title, "股票池");
});

test("layout class separates login and app shell", () => {
  assert.equal(layoutClassForAuthState(false), "auth-layout");
  assert.equal(layoutClassForAuthState(true), "app-shell");
});

test("refresh modes are safe modes", () => {
  assert.deepEqual(refreshModes.map((item) => item.id), ["manual", "conservative", "standard"]);
});

test("formats refresh statuses", () => {
  assert.match(formatRefreshStatus({ status: "rate_limited", error: "retry later" }), /冷却/);
  assert.match(formatRefreshStatus({ status: "failed", error: "down" }), /失败/);
  assert.match(formatRefreshStatus({ status: "succeeded" }), /成功/);
});

test("formats alert messages", () => {
  const message = formatAlertMessage({ signal: "risk_warning", market: "US", symbol: "AAPL", riskLevel: "high", dataTime: "2026-05-05T00:00:00Z" });
  assert.match(message.title, /risk_warning/);
  assert.match(message.meta, /high/);
});
