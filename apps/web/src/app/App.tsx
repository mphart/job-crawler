import React from "react";
import { BrowserRouter, useRoutes } from "react-router-dom";
import { appRoutes } from "./routes";
import { AppProviders } from "./providers";

function AppRouter() {
  const routedElement = useRoutes(appRoutes);
  if (routedElement) return routedElement;

  // DECISION: never render a blank screen when route resolution fails.
  return (
    <main style={{ padding: "1rem", fontFamily: "system-ui, sans-serif", color: "#111827", background: "#ffffff" }}>
      <h2>Route not found</h2>
      <p>Try navigating to /login.</p>
    </main>
  );
}

export function App() {
  return (
    <BrowserRouter>
      <AppProviders>
        <AppRouter />
      </AppProviders>
    </BrowserRouter>
  );
}
