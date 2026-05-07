export function validateRegistration(input) {
  const errors = [];
  if (!input?.email || !String(input.email).includes("@")) errors.push("email_invalid");
  if (!input?.password || String(input.password).length < 8) errors.push("password_too_short");
  return { valid: errors.length === 0, errors };
}

export function authHeader(token) {
  return token ? { Authorization: `Bearer ${token}` } : {};
}
