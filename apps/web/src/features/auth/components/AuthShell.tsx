import React, { PropsWithChildren } from "react";
import "./auth-shell.css";

type AuthShellProps = PropsWithChildren<{
  title: string;
  subtitle?: string;
}>;

export function AuthShell({ title, subtitle, children }: AuthShellProps) {
  return (
    <div className="auth-scene">
      <div className="auth-scene__aurora" aria-hidden />
      <div className="auth-scene__mesh" aria-hidden />
      <div className="auth-scene__noise" aria-hidden />
      <main className="auth-scene__main">
        <header className="auth-scene__header">
          <p className="auth-scene__eyebrow">Job Crawler</p>
          <h1 className="auth-scene__title">{title}</h1>
          {subtitle ? <p className="auth-scene__subtitle">{subtitle}</p> : null}
        </header>
        {children}
      </main>
    </div>
  );
}
