/** Normalizes common user-input issues (whitespace, fullwidth @). */
export function normalizeEmailInput(value: string): string {
  return value
    .trim()
    .replace(/\uFF20/g, "@")
    .replace(/\uFE6B/g, "@");
}

/** Practical check for typical mailbox addresses (local@domain.tld). */
export function isEmail(value: string): boolean {
  const s = normalizeEmailInput(value);
  const at = s.indexOf("@");
  if (at <= 0 || at !== s.lastIndexOf("@")) {
    return false;
  }
  const local = s.slice(0, at);
  const domain = s.slice(at + 1);
  if (!local || !domain || !domain.includes(".")) {
    return false;
  }
  if (domain.startsWith(".") || domain.endsWith(".") || domain.includes("..")) {
    return false;
  }
  return true;
}

export function hasMinLength(value: string, min: number): boolean {
  return value.trim().length >= min;
}
