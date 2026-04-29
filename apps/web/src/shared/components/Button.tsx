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
        borderRadius: "var(--radius-sm)",
        padding: "0.62rem 0.9rem",
        fontWeight: 600,
        letterSpacing: ".01em",
        boxShadow: variant === "secondary" ? "var(--shadow-sm)" : "var(--shadow)",
        transition: "transform .15s ease, filter .2s ease, box-shadow .25s ease",
        outline: "none",
        ...style,
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = "translateY(-1px)";
        e.currentTarget.style.filter = "brightness(1.04)";
        props.onMouseEnter?.(e);
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = "translateY(0)";
        e.currentTarget.style.filter = "brightness(1)";
        props.onMouseLeave?.(e);
      }}
      onFocus={(e) => {
        e.currentTarget.style.boxShadow = `${variant === "secondary" ? "var(--shadow-sm)" : "var(--shadow)"}, 0 0 0 3px color-mix(in srgb, var(--accent) 35%, transparent)`;
        props.onFocus?.(e);
      }}
      onBlur={(e) => {
        e.currentTarget.style.boxShadow = variant === "secondary" ? "var(--shadow-sm)" : "var(--shadow)";
        props.onBlur?.(e);
      }}
    />
  );
}
