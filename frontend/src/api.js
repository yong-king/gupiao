import { authHeader } from "./auth.js";

export const API_BASE = "http://127.0.0.1:8080";

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
    throw new Error(data?.error?.message || `HTTP ${response.status}`);
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
    throw new Error(data?.error?.message || `HTTP ${response.status}`);
  }
  return data;
}
