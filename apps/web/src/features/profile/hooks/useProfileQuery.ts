import { useEffect, useState } from "react";
import { fetchProfile } from "../api/profile.api";
import { Profile } from "../model/profile.types";
export function useProfileQuery(userId:string){const [profile,setProfile]=useState<Profile|null>(null);const [loading,setLoading]=useState(true);useEffect(()=>{setLoading(true);fetchProfile(userId).then((data)=>{setProfile(data);setLoading(false);});},[userId]);return {profile,loading,setProfile};}
