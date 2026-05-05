import { requestJson } from "../../../shared/api/client";
import { Profile } from "../model/profile.types";
import { JobPosting } from "../../feed/model/feed.types";

type ApiUserPreference = {
  keywords?: string[];
  locations?: string[];
  desiredTitles?: string[];
  preferredCompanies?: string[];
  minComp?: number;
  emailOptIn?: boolean;
  darkMode?: boolean;
  Keywords?: string[];
  Locations?: string[];
  DesiredTitles?: string[];
  PreferredCompanies?: string[];
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
  appliedJobs?: ApiJobPosting[];
  ID?: string;
  Username?: string;
  Email?: string;
  IsPrivate?: boolean;
  TotalApplied?: number;
  ResumeFileName?: string;
  Preferences?: ApiUserPreference;
  AppliedJobs?: ApiJobPosting[];
};

type ApiAppliedBy = {
  userId?: string;
  username?: string;
  UserID?: string;
  Username?: string;
};

type ApiJobPosting = {
  id?: string;
  company?: string;
  title?: string;
  location?: string;
  compensation?: string;
  postedAt?: string;
  appliedAt?: string;
  url?: string;
  appliedBy?: ApiAppliedBy[];
  ID?: string;
  Company?: string;
  Title?: string;
  Location?: string;
  Compensation?: string;
  PostedAt?: string;
  AppliedAt?: string;
  URL?: string;
  AppliedBy?: ApiAppliedBy[];
};

function mapPreferences(preferences: ApiUserPreference | undefined): Profile["preferences"] {
  return {
    keywords: preferences?.keywords ?? preferences?.Keywords ?? [],
    locations: preferences?.locations ?? preferences?.Locations ?? [],
    desiredTitles: preferences?.desiredTitles ?? preferences?.DesiredTitles ?? [],
    preferredCompanies: preferences?.preferredCompanies ?? preferences?.PreferredCompanies ?? [],
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
    appliedJobs: mapJobs(profile.appliedJobs ?? profile.AppliedJobs),
  };
}

function mapJobs(jobs: ApiJobPosting[] | undefined): JobPosting[] {
  if (!jobs) return [];
  return jobs.map((job) => ({
    id: job.id ?? job.ID ?? "",
    company: job.company ?? job.Company ?? "Unknown company",
    title: job.title ?? job.Title ?? "Unknown title",
    location: job.location ?? job.Location ?? "Unknown location",
    compensation: job.compensation ?? job.Compensation ?? "",
    postedAt: job.postedAt ?? job.PostedAt ?? new Date(0).toISOString(),
    appliedAt: job.appliedAt ?? job.AppliedAt,
    url: job.url ?? job.URL ?? "",
    appliedBy: (job.appliedBy ?? job.AppliedBy ?? []).map((user) => ({
      userId: user.userId ?? user.UserID ?? "",
      username: user.username ?? user.Username ?? "unknown",
    })),
  }));
}

export async function fetchProfile(userId: string, token: string, signal?: AbortSignal): Promise<Profile> {
  const payload = await requestJson<ApiProfile>(`/api/profiles/${encodeURIComponent(userId)}`, "GET", undefined, { signal, token });
  return mapProfile(payload);
}

export async function updateProfile(update: Partial<Profile> & { resumeContentBase64?: string }, token: string): Promise<Profile> {
  const payload = await requestJson<ApiProfile>("/api/profiles/me", "PATCH", update, { token });
  return mapProfile(payload);
}

type ApiCompanySearchResult = {
  companies?: Array<{ name?: string; isVerified?: boolean }>;
};

export type CompanySearchResult = { name: string; isVerified: boolean };

export async function searchCompanies(query: string, token?: string): Promise<CompanySearchResult[]> {
  const payload = await requestJson<ApiCompanySearchResult>(`/api/companies/search?q=${encodeURIComponent(query)}`, "GET", undefined, token ? { token } : undefined);
  return (payload.companies ?? [])
    .filter((company) => (company.name ?? "").trim().length > 0)
    .map((company) => ({ name: (company.name ?? "").trim(), isVerified: company.isVerified !== false }));
}
