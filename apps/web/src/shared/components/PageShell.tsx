import React from "react";
import { PropsWithChildren, ReactNode } from "react";
export function PageShell({title,actions,children}:PropsWithChildren<{title:string;actions?:ReactNode}>){return <main style={{maxWidth:1100,margin:"0 auto",padding:"1rem"}}><header style={{display:"flex",alignItems:"center",justifyContent:"space-between",marginBottom:"1rem"}}><h1 style={{margin:0}}>{title}</h1>{actions}</header>{children}</main>;}
