import { Profile } from "../model/profile.types";
const MOCK_PROFILE:Profile={id:"u_1",username:"mason",email:"mason@example.com",isPrivate:false,totalApplied:12,resumeFileName:"mason-resume.pdf",preferences:{keywords:["software engineer","frontend"],locations:["Remote","Austin"],desiredTitles:["Software Engineer","Frontend Engineer"],minComp:120000,emailOptIn:true,darkMode:false},appliedJobs:[]};
// MOCK: Replace with API once profile endpoint exists.
export async function fetchProfile(userId:string):Promise<Profile>{return Promise.resolve({...MOCK_PROFILE,id:userId});}
export async function updateProfile(update:Partial<Profile>):Promise<Profile>{return Promise.resolve({...MOCK_PROFILE,...update});}
