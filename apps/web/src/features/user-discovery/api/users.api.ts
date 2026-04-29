import { requestJson } from "../../../shared/api/client";
import { UserSummary } from "../model/users.types";

type ApiUserSummary = {
  id?: string;
  username?: string;
  totalApplied?: number;
  ID?: string;
  Username?: string;
  TotalApplied?: number;
};

function mapUser(user: ApiUserSummary): UserSummary {
  return {
    id: user.id ?? user.ID ?? "",
    username: user.username ?? user.Username ?? "unknown",
    totalApplied: user.totalApplied ?? user.TotalApplied ?? 0,
  };
}

export async function searchUsers(query: string, token: string, signal?: AbortSignal): Promise<UserSummary[]> {
  const params = new URLSearchParams({ q: query });
  const payload = await requestJson<ApiUserSummary[]>(`/api/users/search?${params.toString()}`, "GET", undefined, { signal, token });
  return payload.map(mapUser);
}
