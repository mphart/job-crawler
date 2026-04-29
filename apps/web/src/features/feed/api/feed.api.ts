import { requestJson, requestVoid } from "../../../shared/api/client";
import { FeedFilters, JobPosting } from "../model/feed.types";

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
  url?: string;
  appliedBy?: ApiAppliedBy[];
  ID?: string;
  Company?: string;
  Title?: string;
  Location?: string;
  Compensation?: string;
  PostedAt?: string;
  URL?: string;
  AppliedBy?: ApiAppliedBy[];
};

function mapAppliedBy(users?: ApiAppliedBy[]): JobPosting["appliedBy"] {
  if (!users) return [];
  return users.map((user) => ({
    userId: user.userId ?? user.UserID ?? "",
    username: user.username ?? user.Username ?? "unknown",
  }));
}

function mapJob(job: ApiJobPosting): JobPosting {
  return {
    id: job.id ?? job.ID ?? "",
    company: job.company ?? job.Company ?? "Unknown company",
    title: job.title ?? job.Title ?? "Unknown title",
    location: job.location ?? job.Location ?? "Unknown location",
    compensation: job.compensation ?? job.Compensation ?? "",
    postedAt: job.postedAt ?? job.PostedAt ?? new Date(0).toISOString(),
    url: job.url ?? job.URL ?? "",
    appliedBy: mapAppliedBy(job.appliedBy ?? job.AppliedBy),
  };
}

export function filterAndSortFeedJobs(jobs: JobPosting[], filters: FeedFilters): JobPosting[] {
  const q = filters.search.toLowerCase();
  const filtered = jobs.filter((job) => `${job.company} ${job.title} ${job.location}`.toLowerCase().includes(q));

  return [...filtered].sort((a, b) => {
    if (filters.sortBy === "newest") return b.postedAt.localeCompare(a.postedAt);
    if (filters.sortBy === "company") return a.company.localeCompare(b.company);
    if (filters.sortBy === "title") return a.title.localeCompare(b.title);
    if (filters.sortBy === "location") return a.location.localeCompare(b.location);
    return b.compensation.localeCompare(a.compensation);
  });
}

function feedPath(filters: FeedFilters): string {
  const params = new URLSearchParams();
  if (filters.search.trim()) params.set("search", filters.search.trim());
  params.set("sortBy", filters.sortBy);
  return `/api/feed?${params.toString()}`;
}

export async function fetchFeed(filters: FeedFilters, token: string, signal?: AbortSignal): Promise<JobPosting[]> {
  const payload = await requestJson<ApiJobPosting[]>(feedPath(filters), "GET", undefined, { signal, token });
  return payload.map(mapJob);
}

export async function markApplied(jobId: string, token: string): Promise<void> {
  return requestVoid(`/api/feed/${jobId}/apply`, "POST", undefined, { token });
}

export async function rejectPosting(jobId: string, token: string): Promise<void> {
  return requestVoid(`/api/feed/${jobId}/reject`, "POST", undefined, { token });
}
