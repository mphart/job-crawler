import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { PageShell } from "../../../shared/components/PageShell";
import { SignupForm } from "../components/SignupForm";
export function SignupPage(){const navigate=useNavigate();return <PageShell title="Sign Up"><SignupForm onSuccess={()=>navigate("/feed")}/><p style={{marginTop:"0.85rem",color:"var(--muted)"}}>Already have an account? <Link to="/login">Log in</Link></p></PageShell>;}
