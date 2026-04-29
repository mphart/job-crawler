import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { PageShell } from "../../../shared/components/PageShell";
import { LoginForm } from "../components/LoginForm";
export function LoginPage(){
    const navigate=useNavigate();
    return <PageShell title="Login">
        <LoginForm onSuccess={()=>navigate("/feed")}/>
            <p style={{marginTop:"0.75rem"}}>Don't have an account? 
                <Link to="/signup">Sign up</Link>
            </p>
    </PageShell>;
}
