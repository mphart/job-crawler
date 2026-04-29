import { requestJson } from "../../../shared/api/client";
import { LoginRequest, SessionUser, SignupRequest } from "../model/auth.types";

export async function login(request: LoginRequest): Promise<SessionUser> {
  return requestJson<SessionUser>("/api/auth/login", "POST", request);
}

export async function signup(request: SignupRequest): Promise<SessionUser> {
  return requestJson<SessionUser>("/api/auth/signup", "POST", request);
}
