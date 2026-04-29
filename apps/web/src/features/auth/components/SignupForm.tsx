import React from "react";
import { FormEvent, useState } from "react";
import { Button } from "../../../shared/components/Button";
import { TextInput } from "../../../shared/components/TextInput";
import { isEmail } from "../../../shared/utils/validation";
import { signup } from "../api/auth.api";
import { useSession } from "../hooks/useSession";

export function validateSignupInput(email: string, username: string, password: string): string | null {
  if (!isEmail(email)) return "Please enter a valid email.";
  if (username.trim().length < 2) return "Username must have at least 2 characters.";
  if (password.length < 8) return "Password must have at least 8 characters.";
  return null;
}

export function SignupForm({ onSuccess }: { onSuccess: () => void }) {
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [keywords, setKeywords] = useState("software engineer, frontend");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { setUser } = useSession();

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    const validationError = validateSignupInput(email, username, password);
    if (validationError) {
      setError(validationError);
      return;
    }

    try {
      setError(null);
      setIsSubmitting(true);
      const session = await signup({
        email,
        username,
        password,
        keywords: keywords.split(",").map((keyword) => keyword.trim()).filter(Boolean),
      });
      setUser(session);
      onSuccess();
    } catch {
      setError("Unable to create account right now. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} style={{ display: "grid", gap: "0.75rem", maxWidth: 520 }}>
      <TextInput value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Email" />
      <TextInput value={username} onChange={(e) => setUsername(e.target.value)} placeholder="Username" />
      <TextInput value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Password" type="password" />
      <TextInput value={keywords} onChange={(e) => setKeywords(e.target.value)} placeholder="Keywords (comma separated)" />
      {error ? <small style={{ color: "var(--danger)" }}>{error}</small> : null}
      <Button type="submit" disabled={isSubmitting}>{isSubmitting ? "Creating..." : "Create account"}</Button>
    </form>
  );
}
