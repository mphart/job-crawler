import { FormEvent, useState } from "react";
import { Button } from "../../../shared/components/Button";
import { TextInput } from "../../../shared/components/TextInput";
import { isEmail } from "../../../shared/utils/validation";
import { login } from "../api/auth.api";
import { useSession } from "../hooks/useSession";

export function validateLoginInput(email: string, password: string): string | null {
  if (!isEmail(email)) return "Please enter a valid email.";
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
      const session = await login({ email, password });
      setUser(session);
      onSuccess();
    } catch {
      setError("Unable to sign in right now. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} style={{ display: "grid", gap: "0.75rem", maxWidth: 420 }}>
      <TextInput value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Email" type="email" />
      <TextInput value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Password" type="password" />
      {error ? <small style={{ color: "var(--danger)" }}>{error}</small> : null}
      <Button type="submit" disabled={isSubmitting}>{isSubmitting ? "Signing in..." : "Sign in"}</Button>
    </form>
  );
}
