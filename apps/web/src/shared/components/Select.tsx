import { SelectHTMLAttributes } from "react";
export function Select(props:SelectHTMLAttributes<HTMLSelectElement>){return <select {...props} style={{border:"1px solid var(--border)",borderRadius:8,padding:"0.5rem 0.65rem"}}/>;}
