import test from "node:test";
import assert from "node:assert/strict";

import { validateAccountConfig } from "../src/accounts.js";

test("validates read-only account config", () => {
  const result = validateAccountConfig({ id: "acct-1", alias: "Main", readOnly: true, metadata: {} });
  assert.equal(result.valid, true);
});

test("rejects sensitive account metadata", () => {
  const result = validateAccountConfig({ id: "acct-1", alias: "Main", readOnly: true, metadata: { token: "x" } });
  assert.equal(result.valid, false);
  assert.ok(result.errors.includes("sensitive_metadata_not_allowed"));
});
