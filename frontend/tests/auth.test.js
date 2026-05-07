import test from "node:test";
import assert from "node:assert/strict";

import { authHeader, validateRegistration } from "../src/auth.js";

test("validates registration fields", () => {
  assert.equal(validateRegistration({ email: "bad", password: "short" }).valid, false);
  assert.equal(validateRegistration({ email: "test@example.com", password: "password123" }).valid, true);
});

test("builds bearer auth header", () => {
  assert.deepEqual(authHeader("abc"), { Authorization: "Bearer abc" });
});
