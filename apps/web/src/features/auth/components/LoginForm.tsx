import React from "react";
import { FormEvent, useState } from "react";
import { Button } from "../../../shared/components/Button";
import { TextInput } from "../../../shared/components/TextInput";
import { ApiError } from "../../../shared/api/client";
import { isEmail, normalizeEmailInput } from "../../../shared/utils/validation";
import { login } from "../api/auth.api";
import { useSession } from "../hooks/useSession";

export function validateLoginInput(email: string, password: string): string | null {
  if (!isEmail(normalizeEmailInput(email))) return "Please enter a valid email address.";
  if (password.length < 8) return "Password must have at least 8 characters.";
  return null;
}

export function LoginForm({ onSuccess }: { onSuccess: () => void }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { setUser } = useSession();

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    const validationError = validateLoginInput(email, password);
    if (validationError) {
      setError(validationError);
      return;
    }

    try {
      setIsSubmitting(true);
      setError(null);
      const session = await login({ email: normalizeEmailInput(email), password });
      setUser(session);
      onSuccess();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Unable to sign in right now. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="ui-card auth-card" style={{ display: "grid", gap: "0.85rem", maxWidth: 440, padding: "1.25rem 1.2rem 1.35rem" }}>
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
        placeholder="Password"
        type="password"
        name="password"
        autoComplete="current-password"
        aria-label="Password"
      />
      {error ? <small style={{ color: "var(--danger)", fontWeight: 600 }}>{error}</small> : null}
      <Button type="submit" disabled={isSubmitting}>{isSubmitting ? "Signing in..." : "Sign in"}</Button>
    </form>
  );
}
