import { Profile } from "../model/profile.types";
export function ProfileHeader({profile}:{profile:Profile}){return <section style={{border:"1px solid var(--border)",padding:"0.8rem",borderRadius:10,marginBottom:"0.8rem",background:"var(--surface)"}}><h2 style={{marginTop:0}}>{profile.username}</h2><p style={{margin:0}}>Total Applied: {profile.totalApplied}</p></section>;}
