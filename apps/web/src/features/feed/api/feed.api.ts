import { mockRequestJson } from "../../../shared/api/client";
import { FeedFilters, JobPosting } from "../model/feed.types";

const MOCK_JOBS: JobPosting[] = [
  {
    id: "job_1",
    company: "Cisco",
    title: "Software Engineer",
    location: "Austin, TX",
    compensation: "$130k-$155k",
    postedAt: new Date(Date.now() - 86400000).toISOString(),
    url: "https://example.com/jobs/1",
    appliedBy: [{ userId: "u_2", username: "alex" }],
  },
  {
    id: "job_2",
    company: "365 Retail Markets",
    title: "Frontend Engineer",
    location: "Remote",
    compensation: "$120k-$140k",
    postedAt: new Date(Date.now() - 172800000).toISOString(),
    url: "https://example.com/jobs/2",
    appliedBy: [],
  },
];

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

export async function fetchFeed(filters: FeedFilters, signal?: AbortSignal): Promise<JobPosting[]> {
  // MOCK: Replace with backend feed endpoint.
  return mockRequestJson(() => filterAndSortFeedJobs(MOCK_JOBS, filters), { signal });
}

export async function markApplied(jobId: string): Promise<void> {
  return mockRequestJson(() => {
    void jobId;
  });
}

export async function rejectPosting(jobId: string): Promise<void> {
  return mockRequestJson(() => {
    void jobId;
  });
}
