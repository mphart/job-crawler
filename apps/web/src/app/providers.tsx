import React from "react";
import { PropsWithChildren } from "react";
import { SessionProvider } from "../features/auth/hooks/useSession";
import { ThemeProvider } from "../shared/theme/ThemeProvider";

export function AppProviders({ children }: PropsWithChildren) {
  return (
    <ThemeProvider>
      <SessionProvider>{children}</SessionProvider>
    </ThemeProvider>
  );
}
