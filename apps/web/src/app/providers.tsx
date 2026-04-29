import React from "react";
import { PropsWithChildren } from "react";
import { SessionProvider } from "../features/auth/hooks/useSession";

export function AppProviders({ children }: PropsWithChildren) {
  return <SessionProvider>{children}</SessionProvider>;
}
