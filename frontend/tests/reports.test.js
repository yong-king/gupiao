import test from "node:test";
import assert from "node:assert/strict";

import { formatDailyReport } from "../src/reports.js";

test("formats populated daily report", () => {
  const result = formatDailyReport({
    date: "2026-05-05",
    summary: "Watch alerts.",
    riskPoints: ["US:AAPL risk"],
    needsConfirmation: ["Check AAPL"],
    dataTime: "2026-05-05T08:00:00Z",
  });

  assert.equal(result.empty, false);
  assert.equal(result.riskPoints.length, 1);
  assert.equal(result.dataTime, "2026-05-05T08:00:00Z");
});

test("formats empty daily report", () => {
  const result = formatDailyReport({ summary: "No alerts." });

  assert.equal(result.empty, true);
  assert.equal(result.summary, "No alerts.");
});
