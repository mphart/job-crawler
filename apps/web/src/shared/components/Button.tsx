import React from "react";
import { ButtonHTMLAttributes } from "react";

type Variant = "primary" | "secondary" | "danger";
type Props = ButtonHTMLAttributes<HTMLButtonElement> & { variant?: Variant };

export function Button({ variant = "primary", style, ...props }: Props) {
  const bg =
    variant === "primary"
      ? "var(--primary)"
      : variant === "danger"
        ? "var(--danger)"
        : "var(--surface-elevated)";
  const color = variant === "secondary" ? "var(--text)" : "white";

  return (
    <button
      {...props}
      style={{
        border: variant === "secondary" ? "1px solid var(--border)" : "1px solid transparent",
        background: bg,
        color,
        borderRadius: 10,
        padding: "0.5rem 0.8rem",
        fontWeight: 600,
        boxShadow: variant === "secondary" ? "none" : "0 6px 18px rgba(37,99,235,.24)",
        transition: "transform .1s ease, filter .2s ease",
        ...style,
      }}
    />
  );
}
