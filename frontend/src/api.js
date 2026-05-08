import { authHeader } from "./auth.js";

export const API_BASE = "http://127.0.0.1:8080";

export function shouldInvalidateSession(data, status) {
  const message = data?.error?.message || "";
  return status === 401 || /invalid or expired token/i.test(message);
}

function handleAPIError(response, data) {
  if (shouldInvalidateSession(data, response.status)) {
    const browser = globalThis.window;
    browser?.localStorage?.removeItem("jijin_token");
    browser?.localStorage?.removeItem("jijin_email");
    browser?.localStorage?.removeItem("jijin_user_id");
    if (browser?.dispatchEvent && typeof CustomEvent !== "undefined") {
      browser.dispatchEvent(new CustomEvent("jijin-auth-expired"));
    }
    throw new Error("登录已过期，请重新登录。");
  }
  throw new Error(data?.error?.message || `HTTP ${response.status}`);
}

export async function postJSON(path, body, token = "") {
  const response = await fetch(`${API_BASE}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Request-ID": `web-${Date.now()}`,
      ...authHeader(token),
    },
    body: JSON.stringify(body),
  });
  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    handleAPIError(response, data);
  }
  return data;
}

export async function getJSON(path, token = "") {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      "X-Request-ID": `web-${Date.now()}`,
      ...authHeader(token),
    },
  });
  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    handleAPIError(response, data);
  }
  return data;
}

export async function deleteJSON(path, token = "") {
  const response = await fetch(`${API_BASE}${path}`, {
    method: "DELETE",
    headers: {
      "X-Request-ID": `web-${Date.now()}`,
      ...authHeader(token),
    },
  });
  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    handleAPIError(response, data);
  }
  return data;
}
