import React from "react";
import { InputHTMLAttributes } from "react";
export function TextInput(props:InputHTMLAttributes<HTMLInputElement>){return <input {...props} style={{border:"1px solid var(--border)",borderRadius:8,padding:"0.55rem 0.7rem",width:"100%"}}/>;}
