import React from "react";
import { PropsWithChildren, ReactNode } from "react";

export function PageShell({ title, actions, children }: PropsWithChildren<{ title: string; actions?: ReactNode }>) {
  return (
    <main style={{ maxWidth: 1140, margin: "0 auto", padding: "1.7rem 1.1rem 2.4rem" }}>
      <header
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          gap: "1rem",
          marginBottom: "1.15rem",
          background: "var(--surface)",
          border: "1px solid var(--border)",
          borderRadius: "var(--radius-lg)",
          padding: "1rem 1.15rem",
          boxShadow: "var(--shadow)",
        }}
      >
        <h1 style={{ margin: 0, fontSize: "1.42rem", fontWeight: 800 }}>{title}</h1>
        {actions}
      </header>
      {children}
    </main>
  );
}
