import { ChangeEvent } from "react";
export function EmailOptInToggle({value,onChange}:{value:boolean;onChange:(next:boolean)=>void;}){return <label><input type="checkbox" checked={value} onChange={(e:ChangeEvent<HTMLInputElement>)=>onChange(e.target.checked)}/> Get Emails</label>;}
