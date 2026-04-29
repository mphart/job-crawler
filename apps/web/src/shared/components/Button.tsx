import React from "react";
import { ButtonHTMLAttributes } from "react";
type Variant="primary"|"secondary"|"danger";
type Props=ButtonHTMLAttributes<HTMLButtonElement>&{variant?:Variant};
export function Button({variant="primary",style,...props}:Props){const bg=variant==="primary"?"var(--primary)":variant==="danger"?"var(--danger)":"transparent";const color=variant==="secondary"?"var(--text)":"white";return <button {...props} style={{border:"1px solid var(--border)",background:bg,color,borderRadius:8,padding:"0.5rem 0.75rem",...style}}/>;}
