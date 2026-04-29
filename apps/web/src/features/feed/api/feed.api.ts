import { requestJson, requestVoid } from "../../../shared/api/client";
import { FeedFilters, JobPosting } from "../model/feed.types";

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
  return requestJson<JobPosting[]>(feedPath(filters), "GET", undefined, { signal, token });
}

export async function markApplied(jobId: string, token: string): Promise<void> {
  return requestVoid(`/api/feed/${jobId}/apply`, "POST", undefined, { token });
}

export async function rejectPosting(jobId: string, token: string): Promise<void> {
  return requestVoid(`/api/feed/${jobId}/reject`, "POST", undefined, { token });
}
