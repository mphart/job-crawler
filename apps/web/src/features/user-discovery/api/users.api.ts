import { requestJson } from "../../../shared/api/client";
import { UserSummary } from "../model/users.types";

export async function searchUsers(query: string, token: string, signal?: AbortSignal): Promise<UserSummary[]> {
  const params = new URLSearchParams({ q: query });
  return requestJson<UserSummary[]>(`/api/users/search?${params.toString()}`, "GET", undefined, { signal, token });
}
