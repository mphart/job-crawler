import React from "react";
import { PropsWithChildren, ReactNode } from "react";

export function PageShell({ title, actions, children }: PropsWithChildren<{ title: string; actions?: ReactNode }>) {
  return (
    <main style={{ maxWidth: 1120, margin: "0 auto", padding: "1.25rem 1rem 2rem" }}>
      <header
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          gap: "0.75rem",
          marginBottom: "1rem",
          background: "var(--surface)",
          border: "1px solid var(--border)",
          borderRadius: 14,
          padding: "0.9rem 1rem",
          boxShadow: "var(--shadow)",
        }}
      >
        <h1 style={{ margin: 0, fontSize: "1.2rem" }}>{title}</h1>
        {actions}
      </header>
      {children}
    </main>
  );
}
