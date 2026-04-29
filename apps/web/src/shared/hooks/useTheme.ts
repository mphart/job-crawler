import { useEffect, useState } from "react";
type Theme="light"|"dark";
export function useTheme(){const [theme,setTheme]=useState<Theme>("light");useEffect(()=>{document.documentElement.dataset.theme=theme;},[theme]);return {theme,toggleTheme:()=>setTheme((p)=>p==="light"?"dark":"light")};}
