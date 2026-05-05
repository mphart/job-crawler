import { JobPosting } from "../../feed/model/feed.types";
export type UserPreference={keywords:string[];locations:string[];desiredTitles:string[];preferredCompanies:string[];minComp:number;emailOptIn:boolean;darkMode:boolean};
export type Profile={id:string;username:string;email:string;isPrivate:boolean;totalApplied:number;resumeFileName?:string;preferences:UserPreference;appliedJobs:JobPosting[]};
