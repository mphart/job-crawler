import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { AuthShell } from "../components/AuthShell";
import { SignupForm } from "../components/SignupForm";

export function SignupPage() {
  const navigate = useNavigate();
  return (
    <AuthShell
      title="Create your account"
      subtitle="Tell us how you want to be matched to roles. After registration you will receive a formal confirmation by email when outbound mail is configured."
    >
      <SignupForm onSuccess={() => navigate("/feed")} />
      <p style={{ marginTop: "1.1rem", color: "var(--muted)", fontSize: "0.95rem" }}>
        Already registered? <Link to="/login">Sign in</Link>
      </p>
    </AuthShell>
  );
}
