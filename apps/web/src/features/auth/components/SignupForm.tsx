import React from "react";
import { FormEvent, useState } from "react";
import { Button } from "../../../shared/components/Button";
import { TextInput } from "../../../shared/components/TextInput";
import { ApiError } from "../../../shared/api/client";
import { isEmail, normalizeEmailInput } from "../../../shared/utils/validation";
import { fileToBase64 } from "../../../shared/utils/file";
import { signup } from "../api/auth.api";
import { useSession } from "../hooks/useSession";
import { useTheme } from "../../../shared/hooks/useTheme";
import type { NotificationFrequency } from "../../notifications/model/notifications.types";
import { PreferredCompaniesPicker } from "../../profile/components/PreferredCompaniesPicker";

const MAX_RESUME_BYTES = 4 * 1024 * 1024;

export function splitCommaPrefs(line: string): string[] {
  return line.split(",").map((s) => s.trim()).filter(Boolean);
}

export function validateSignupInput(
  email: string,
  fullName: string,
  password: string,
  keywordsLine: string,
  locationsLine: string,
): string | null {
  if (!isEmail(normalizeEmailInput(email))) return "Please enter a valid email address.";
  if (fullName.trim().length < 2) return "Please enter your full name (at least 2 characters).";
  if (password.length < 8) return "Password must be at least 8 characters.";
  if (splitCommaPrefs(keywordsLine).length === 0) return "Add at least one job keyword (comma-separated).";
  if (splitCommaPrefs(locationsLine).length === 0) return "Add at least one target location (comma-separated).";
  return null;
}

function readStoredDarkMode(): boolean {
  try {
    return localStorage.getItem("job-crawler-theme") === "dark";
  } catch {
    return false;
  }
}

export function SignupForm({ onSuccess }: { onSuccess: () => void }) {
  const { setTheme } = useTheme();
  const [email, setEmail] = useState("");
  const [fullName, setFullName] = useState("");
  const [password, setPassword] = useState("");
  const [digestFrequency, setDigestFrequency] = useState<NotificationFrequency>("daily");
  const [keywordsLine, setKeywordsLine] = useState("software engineer, golang");
  const [locationsLine, setLocationsLine] = useState("Remote, United States");
  const [titlesLine, setTitlesLine] = useState("Software Engineer, Backend Engineer");
  const [preferredCompanies, setPreferredCompanies] = useState<string[]>([]);
  const [minComp, setMinComp] = useState("120000");
  const [emailOptIn, setEmailOptIn] = useState(true);
  const [darkMode, setDarkMode] = useState(readStoredDarkMode);
  const [resumeName, setResumeName] = useState("");
  const [resumeBase64, setResumeBase64] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { setUser } = useSession();

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    const validationError = validateSignupInput(email, fullName, password, keywordsLine, locationsLine);
    if (validationError) {
      setError(validationError);
      return;
    }
    const minParsed = Math.max(0, parseInt(minComp.replace(/[^0-9]/g, ""), 10) || 0);

    try {
      setError(null);
      setIsSubmitting(true);
      const session = await signup({
        email: normalizeEmailInput(email),
        name: fullName.trim(),
        password,
        notificationFrequency: digestFrequency,
        preferences: {
          keywords: splitCommaPrefs(keywordsLine),
          locations: splitCommaPrefs(locationsLine),
          desiredTitles: splitCommaPrefs(titlesLine),
          preferredCompanies,
          minComp: minParsed,
          emailOptIn,
          darkMode,
        },
        resumeFileName: resumeName.trim() || undefined,
        resumeContentBase64: resumeBase64.trim() || undefined,
      });
      setUser(session);
      onSuccess();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Unable to create account right now. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="ui-card auth-card" style={{ display: "grid", gap: "1.1rem", padding: "1.35rem 1.25rem 1.45rem" }}>
      <fieldset className="auth-fieldset">
        <legend>Account</legend>
        <div style={{ display: "grid", gap: "0.75rem" }}>
          <TextInput
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
            placeholder="Full legal name"
            name="name"
            autoComplete="name"
            aria-label="Full name"
          />
          <TextInput
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Email address"
            type="email"
            name="email"
            autoComplete="email"
            aria-label="Email address"
          />
          <TextInput
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Password (8+ characters)"
            type="password"
            name="password"
            autoComplete="new-password"
            aria-label="Password"
          />
        </div>
      </fieldset>

      <fieldset className="auth-fieldset">
        <legend>Job digest emails</legend>
        <p className="auth-hint">How often should we email you matching roles? (Requires email alerts below.)</p>
        <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}>
          <label style={{ display: "flex", gap: "0.5rem", alignItems: "center", fontSize: "0.95rem" }}>
            <input type="radio" name="digestFrequency" value="daily" checked={digestFrequency === "daily"} onChange={() => setDigestFrequency("daily")} />
            Once a day
          </label>
          <label style={{ display: "flex", gap: "0.5rem", alignItems: "center", fontSize: "0.95rem" }}>
            <input type="radio" name="digestFrequency" value="twice-daily" checked={digestFrequency === "twice-daily"} onChange={() => setDigestFrequency("twice-daily")} />
            Twice a day
          </label>
          <label style={{ display: "flex", gap: "0.5rem", alignItems: "center", fontSize: "0.95rem" }}>
            <input type="radio" name="digestFrequency" value="weekly" checked={digestFrequency === "weekly"} onChange={() => setDigestFrequency("weekly")} />
            Once a week
          </label>
        </div>
      </fieldset>

      <fieldset className="auth-fieldset">
        <legend>Job preferences</legend>
        <div style={{ display: "grid", gap: "0.75rem" }}>
          <TextInput
            value={keywordsLine}
            onChange={(e) => setKeywordsLine(e.target.value)}
            placeholder="Keywords (comma-separated)"
            aria-label="Job keywords"
          />
          <p className="auth-hint">Example: software engineer, kubernetes, typescript</p>
          <TextInput
            value={locationsLine}
            onChange={(e) => setLocationsLine(e.target.value)}
            placeholder="Locations (comma-separated)"
            aria-label="Target locations"
          />
          <TextInput
            value={titlesLine}
            onChange={(e) => setTitlesLine(e.target.value)}
            placeholder="Desired titles (comma-separated, optional)"
            aria-label="Desired job titles"
          />
          <TextInput
            value={minComp}
            onChange={(e) => setMinComp(e.target.value)}
            placeholder="Minimum compensation (annual, USD)"
            inputMode="numeric"
            aria-label="Minimum compensation"
          />
          <div style={{ display: "flex", flexWrap: "wrap", gap: "1rem", marginTop: "0.25rem" }}>
            <label style={{ display: "flex", gap: "0.45rem", alignItems: "center", fontSize: "0.92rem" }}>
              <input type="checkbox" checked={emailOptIn} onChange={(e) => setEmailOptIn(e.target.checked)} /> Job alert emails
            </label>
            <label style={{ display: "flex", gap: "0.45rem", alignItems: "center", fontSize: "0.92rem" }}>
              <input
                type="checkbox"
                checked={darkMode}
                onChange={(e) => {
                  const next = e.target.checked;
                  setDarkMode(next);
                  setTheme(next ? "dark" : "light");
                }}
              />{" "}
              Dark theme
            </label>
          </div>
        </div>
      </fieldset>

      <fieldset className="auth-fieldset">
        <legend>Resume (PDF)</legend>
        <p className="auth-hint">Strongly recommended. Maximum size {Math.round(MAX_RESUME_BYTES / (1024 * 1024))} MB.</p>
        <input
          type="file"
          accept="application/pdf,.pdf"
          aria-label="Resume PDF file"
          onChange={async (event) => {
            const file = event.target.files?.[0];
            if (!file) return;
            if (file.type !== "application/pdf" && !file.name.toLowerCase().endsWith(".pdf")) {
              setError("Please upload a PDF resume.");
              return;
            }
            if (file.size > MAX_RESUME_BYTES) {
              setError(`Resume must be under ${MAX_RESUME_BYTES / (1024 * 1024)} MB.`);
              return;
            }
            setResumeName(file.name);
            setResumeBase64(await fileToBase64(file));
            setError(null);
          }}
        />
        {resumeName ? <small style={{ color: "var(--muted)" }}>Selected: {resumeName}</small> : <small style={{ color: "var(--muted)" }}>Optional — you can upload later from your profile.</small>}
      </fieldset>
      <PreferredCompaniesPicker selected={preferredCompanies} onChange={setPreferredCompanies} />

      {error ? <small style={{ color: "var(--danger)", fontWeight: 600 }}>{error}</small> : null}
      <Button type="submit" disabled={isSubmitting}>
        {isSubmitting ? "Creating account…" : "Create account"}
      </Button>
    </form>
  );
}
