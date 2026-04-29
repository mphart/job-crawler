import React from "react";
import { Profile } from "../model/profile.types";
export function ProfileHeader({profile}:{profile:Profile}){return <section className="ui-card" style={{padding:"1rem",borderRadius:"var(--radius-md)",marginBottom:"0.85rem",background:"linear-gradient(135deg, color-mix(in srgb, var(--primary) 10%, var(--surface)), var(--surface-elevated))"}}><h2 style={{marginTop:0,marginBottom:".35rem",fontSize:"1.35rem"}}>{profile.username}</h2><p className="ui-pill" style={{margin:0}}>Total Applied: {profile.totalApplied}</p></section>;}
