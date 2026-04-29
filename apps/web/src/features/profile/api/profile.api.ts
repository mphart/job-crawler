import { requestJson } from "../../../shared/api/client";
import { Profile } from "../model/profile.types";

export async function fetchProfile(userId: string, token: string, signal?: AbortSignal): Promise<Profile> {
  return requestJson<Profile>(`/api/profiles/${encodeURIComponent(userId)}`, "GET", undefined, { signal, token });
}

export async function updateProfile(update: Partial<Profile>, token: string): Promise<Profile> {
  return requestJson<Profile>("/api/profiles/me", "PATCH", update, { token });
}
