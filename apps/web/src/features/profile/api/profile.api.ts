import { requestJson } from "../../../shared/api/client";
import { Profile } from "../model/profile.types";

type ApiUserPreference = {
  keywords?: string[];
  locations?: string[];
  desiredTitles?: string[];
  minComp?: number;
  emailOptIn?: boolean;
  darkMode?: boolean;
  Keywords?: string[];
  Locations?: string[];
  DesiredTitles?: string[];
  MinComp?: number;
  EmailOptIn?: boolean;
  DarkMode?: boolean;
};

type ApiProfile = {
  id?: string;
  username?: string;
  email?: string;
  isPrivate?: boolean;
  totalApplied?: number;
  resumeFileName?: string;
  preferences?: ApiUserPreference;
  appliedJobs?: Profile["appliedJobs"];
  ID?: string;
  Username?: string;
  Email?: string;
  IsPrivate?: boolean;
  TotalApplied?: number;
  ResumeFileName?: string;
  Preferences?: ApiUserPreference;
  AppliedJobs?: Profile["appliedJobs"];
};

function mapPreferences(preferences: ApiUserPreference | undefined): Profile["preferences"] {
  return {
    keywords: preferences?.keywords ?? preferences?.Keywords ?? [],
    locations: preferences?.locations ?? preferences?.Locations ?? [],
    desiredTitles: preferences?.desiredTitles ?? preferences?.DesiredTitles ?? [],
    minComp: preferences?.minComp ?? preferences?.MinComp ?? 0,
    emailOptIn: preferences?.emailOptIn ?? preferences?.EmailOptIn ?? false,
    darkMode: preferences?.darkMode ?? preferences?.DarkMode ?? false,
  };
}

function mapProfile(profile: ApiProfile): Profile {
  return {
    id: profile.id ?? profile.ID ?? "",
    username: profile.username ?? profile.Username ?? "unknown",
    email: profile.email ?? profile.Email ?? "",
    isPrivate: profile.isPrivate ?? profile.IsPrivate ?? false,
    totalApplied: profile.totalApplied ?? profile.TotalApplied ?? 0,
    resumeFileName: profile.resumeFileName ?? profile.ResumeFileName,
    preferences: mapPreferences(profile.preferences ?? profile.Preferences),
    appliedJobs: profile.appliedJobs ?? profile.AppliedJobs ?? [],
  };
}

export async function fetchProfile(userId: string, token: string, signal?: AbortSignal): Promise<Profile> {
  const payload = await requestJson<ApiProfile>(`/api/profiles/${encodeURIComponent(userId)}`, "GET", undefined, { signal, token });
  return mapProfile(payload);
}

export async function updateProfile(update: Partial<Profile>, token: string): Promise<Profile> {
  const payload = await requestJson<ApiProfile>("/api/profiles/me", "PATCH", update, { token });
  return mapProfile(payload);
}
