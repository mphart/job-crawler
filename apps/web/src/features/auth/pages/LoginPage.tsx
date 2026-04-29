import { useNavigate } from "react-router-dom";
import { PageShell } from "../../../shared/components/PageShell";
import { LoginForm } from "../components/LoginForm";
export function LoginPage(){const navigate=useNavigate();return <PageShell title="Login"><LoginForm onSuccess={()=>navigate("/feed")}/></PageShell>;}
