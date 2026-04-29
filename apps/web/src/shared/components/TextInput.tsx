import React from "react";
import { InputHTMLAttributes } from "react";

export function TextInput(props: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      style={{
        border: "1px solid var(--border)",
        borderRadius: "var(--radius-sm)",
        padding: "0.68rem 0.82rem",
        width: "100%",
        background: "color-mix(in srgb, var(--surface-elevated) 94%, white 6%)",
        color: "var(--text)",
        boxShadow: "inset 0 1px 0 rgba(255,255,255,.45)",
        transition: "border-color .2s ease, box-shadow .2s ease, transform .15s ease",
      }}
      onFocus={(e) => {
        e.currentTarget.style.borderColor = "var(--primary)";
        e.currentTarget.style.boxShadow = "0 0 0 3px color-mix(in srgb, var(--accent) 35%, transparent)";
        props.onFocus?.(e);
      }}
      onBlur={(e) => {
        e.currentTarget.style.borderColor = "var(--border)";
        e.currentTarget.style.boxShadow = "inset 0 1px 0 rgba(255,255,255,.45)";
        props.onBlur?.(e);
      }}
    />
  );
}
