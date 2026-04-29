import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { TextInput } from "../../../shared/components/TextInput";
import { searchUsers } from "../api/users.api";
import { UserSummary } from "../model/users.types";
export function UserSearchBar(){const [query,setQuery]=useState("");const [results,setResults]=useState<UserSummary[]>([]);useEffect(()=>{if(!query.trim()){setResults([]);return;}searchUsers(query).then(setResults);},[query]);return <div style={{position:"relative",minWidth:220}}><TextInput value={query} onChange={(e)=>setQuery(e.target.value)} placeholder="Find users"/>{results.length>0?<div style={{position:"absolute",top:"105%",left:0,right:0,border:"1px solid var(--border)",borderRadius:8,background:"var(--surface)",zIndex:10}}>{results.map((user)=><Link key={user.id} to={`/profile/${user.id}`} style={{display:"block",padding:"0.5rem"}}>{user.username} ({user.totalApplied} applied)</Link>)}</div>:null}</div>;}
