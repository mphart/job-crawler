import { UserSummary } from "../model/users.types";
const MOCK_USERS:UserSummary[]=[{id:"u_2",username:"alex",totalApplied:44},{id:"u_3",username:"jamie",totalApplied:27},{id:"u_4",username:"taylor",totalApplied:16}];
// MOCK: Replace with backend user search endpoint.
export async function searchUsers(query:string):Promise<UserSummary[]>{const q=query.toLowerCase();return Promise.resolve(MOCK_USERS.filter((u)=>u.username.toLowerCase().includes(q)));}
