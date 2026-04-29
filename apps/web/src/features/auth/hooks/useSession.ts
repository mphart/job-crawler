import { createContext, PropsWithChildren, ReactNode, createElement, useContext, useEffect, useMemo, useState } from "react";
import { SessionUser } from "../model/auth.types";

const SESSION_KEY = "job-crawler.session";

type SessionContextType = {
  user: SessionUser | null;
  isHydrating: boolean;
  setUser: (user: SessionUser | null) => void;
  signOut: () => void;
};

const SessionContext = createContext<SessionContextType | null>(null);

export function SessionProvider({ children }: PropsWithChildren) {
  const [user, setUserState] = useState<SessionUser | null>(null);
  const [isHydrating, setIsHydrating] = useState(true);

  useEffect(() => {
    const raw = window.localStorage.getItem(SESSION_KEY);
    if (raw) {
      try {
        setUserState(JSON.parse(raw) as SessionUser);
      } catch {
        window.localStorage.removeItem(SESSION_KEY);
      }
    }
    setIsHydrating(false);
  }, []);

  const setUser = (next: SessionUser | null) => {
    setUserState(next);
    if (next) {
      window.localStorage.setItem(SESSION_KEY, JSON.stringify(next));
    } else {
      window.localStorage.removeItem(SESSION_KEY);
    }
  };

  const value = useMemo<SessionContextType>(
    () => ({ user, isHydrating, setUser, signOut: () => setUser(null) }),
    [user, isHydrating]
  );

  return createElement(SessionContext.Provider, { value }, children as ReactNode);
}

export function useSession() {
  const context = useContext(SessionContext);
  if (!context) {
    throw new Error("useSession must be used within SessionProvider");
  }
  return context;
}
