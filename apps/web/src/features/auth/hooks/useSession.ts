import { createContext, PropsWithChildren, useContext, useMemo, useState } from "react";
import { SessionUser } from "../model/auth.types";
const SessionContext=createContext<{user:SessionUser|null;setUser:(user:SessionUser|null)=>void}|null>(null);
export function SessionProvider({children}:PropsWithChildren){const [user,setUser]=useState<SessionUser|null>(null);const value=useMemo(()=>({user,setUser}),[user]);return <SessionContext.Provider value={value}>{children}</SessionContext.Provider>;}
export function useSession(){const context=useContext(SessionContext);if(!context){throw new Error("useSession must be used within SessionProvider");}return context;}
