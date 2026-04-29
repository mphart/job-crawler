import React from "react";
import { useNavigate } from "react-router-dom";
import { PageShell } from "../../../shared/components/PageShell";
import { SignupForm } from "../components/SignupForm";
export function SignupPage(){const navigate=useNavigate();return <PageShell title="Sign Up"><SignupForm onSuccess={()=>navigate("/feed")}/></PageShell>;}
