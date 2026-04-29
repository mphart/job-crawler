import { FormEvent, useState } from "react";
import { Button } from "../../../shared/components/Button";
import { TextInput } from "../../../shared/components/TextInput";
import { isEmail } from "../../../shared/utils/validation";
import { login } from "../api/auth.api";
export function LoginForm({onSuccess}:{onSuccess:()=>void}){const [email,setEmail]=useState("");const [password,setPassword]=useState("");const [error,setError]=useState<string|null>(null);async function handleSubmit(event:FormEvent){event.preventDefault();if(!isEmail(email)){setError("Please enter a valid email.");return;}if(password.length<8){setError("Password must have at least 8 characters.");return;}await login({email,password});onSuccess();}return <form onSubmit={handleSubmit} style={{display:"grid",gap:"0.75rem",maxWidth:420}}><TextInput value={email} onChange={(e)=>setEmail(e.target.value)} placeholder="Email" type="email"/><TextInput value={password} onChange={(e)=>setPassword(e.target.value)} placeholder="Password" type="password"/>{error?<small style={{color:"var(--danger)"}}>{error}</small>:null}<Button type="submit">Sign in</Button></form>;}
