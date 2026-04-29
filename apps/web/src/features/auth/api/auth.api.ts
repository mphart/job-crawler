import { LoginRequest, SessionUser, SignupRequest } from "../model/auth.types";
// MOCK: Replace with backend API integration when auth endpoints are available.
export async function login(request:LoginRequest):Promise<SessionUser>{return Promise.resolve({id:"u_1",email:request.email,username:"demo-user"});}
export async function signup(request:SignupRequest):Promise<SessionUser>{return Promise.resolve({id:"u_new",email:request.email,username:request.username});}
