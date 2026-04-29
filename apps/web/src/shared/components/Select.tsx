import React from "react";
import { SelectHTMLAttributes } from "react";

export function Select(props: SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      {...props}
      style={{
        border: "1px solid var(--border)",
        borderRadius: 10,
        padding: "0.58rem 0.68rem",
        background: "var(--surface-elevated)",
        color: "var(--text)",
      }}
    />
  );
}
