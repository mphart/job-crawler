import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { AuthShell } from "../components/AuthShell";
import { LoginForm } from "../components/LoginForm";

export function LoginPage() {
  const navigate = useNavigate();
  return (
    <AuthShell
      title="Sign in"
      subtitle="Use the email address and password for your account. No other information is required on this page."
    >
      <LoginForm onSuccess={() => navigate("/feed")} />
      <p style={{ marginTop: "1.1rem", color: "var(--muted)", fontSize: "0.95rem" }}>
        Don&apos;t have an account? <Link to="/signup">Create an account</Link>
      </p>
    </AuthShell>
  );
}
