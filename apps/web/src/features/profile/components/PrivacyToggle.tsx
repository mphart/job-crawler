import { ChangeEvent } from "react";
export function PrivacyToggle({value,onChange}:{value:boolean;onChange:(next:boolean)=>void;}){return <label><input type="checkbox" checked={!value} onChange={(e:ChangeEvent<HTMLInputElement>)=>onChange(!e.target.checked)}/> Public Profile</label>;}
