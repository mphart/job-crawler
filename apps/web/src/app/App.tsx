import { BrowserRouter, useRoutes } from "react-router-dom";
import { appRoutes } from "./routes";
import { AppProviders } from "./providers";
function AppRouter(){return useRoutes(appRoutes);}
export function App(){return <BrowserRouter><AppProviders><AppRouter/></AppProviders></BrowserRouter>;}
