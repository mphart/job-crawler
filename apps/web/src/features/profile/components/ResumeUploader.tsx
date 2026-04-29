import React from "react";
import { useState } from "react";
import { Button } from "../../../shared/components/Button";
export function ResumeUploader({current,onUpload}:{current?:string;onUpload:(fileName:string)=>void;}){const [value,setValue]=useState(current??"");return <section style={{display:"flex",gap:"0.5rem",alignItems:"center"}}><input value={value} onChange={(e)=>setValue(e.target.value)} placeholder="resume.pdf"/><Button onClick={()=>onUpload(value||"resume.pdf")}>Upload Resume</Button></section>;}
