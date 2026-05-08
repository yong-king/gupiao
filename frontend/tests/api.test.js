import test from "node:test";
import assert from "node:assert/strict";

import { shouldInvalidateSession } from "../src/api.js";

test("detects expired auth token responses", () => {
  assert.equal(shouldInvalidateSession({ error: { message: "Invalid or expired token." } }, 400), true);
  assert.equal(shouldInvalidateSession({ error: { message: "Authentication required." } }, 401), true);
  assert.equal(shouldInvalidateSession({ error: { message: "market and symbol are required." } }, 400), false);
});
