export function validateAccountConfig(input) {
  const errors = [];
  const metadata = input?.metadata || {};
  if (!input?.id) errors.push("account_id_required");
  if (!input?.alias) errors.push("alias_required");
  if (input?.readOnly !== true) errors.push("read_only_required");
  for (const key of Object.keys(metadata)) {
    const normalized = key.toLowerCase();
    if (normalized.includes("password") || normalized.includes("secret") || normalized.includes("token")) {
      errors.push("sensitive_metadata_not_allowed");
    }
  }
  return { valid: errors.length === 0, errors };
}
