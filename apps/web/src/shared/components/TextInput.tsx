import React from "react";
import { InputHTMLAttributes } from "react";

export function TextInput(props: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      style={{
        border: "1px solid var(--border)",
        borderRadius: 10,
        padding: "0.62rem 0.75rem",
        width: "100%",
        background: "var(--surface-elevated)",
        color: "var(--text)",
      }}
    />
  );
}
