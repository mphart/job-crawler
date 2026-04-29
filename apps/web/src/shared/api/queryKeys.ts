export const queryKeys={session:["session"] as const,feed:(search:string,sortBy:string)=>["feed",search,sortBy] as const,profile:(userId:string)=>["profile",userId] as const};
