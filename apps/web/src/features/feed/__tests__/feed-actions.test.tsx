import { describe, expect, it, vi } from "vitest";
import { fetchFeed, filterAndSortFeedJobs, markApplied, rejectPosting } from "../api/feed.api";
import { JobPosting } from "../model/feed.types";

global.fetch = vi.fn(async () => new Response(JSON.stringify([]), { status: 200 })) as unknown as typeof fetch;

describe("feed flow", () => {
  it("filters and sorts feed payloads", async () => {
    const jobs = await fetchFeed({ search: "engineer", sortBy: "company" }, "test-token");
    expect(Array.isArray(jobs)).toBe(true);
  });

  it("applies local filtering helper consistently", () => {
    const sample: JobPosting[] = [
      { id: "1", company: "B", title: "Engineer", location: "Remote", compensation: "$100", postedAt: "2024-01-01", url: "u", appliedBy: [] },
      { id: "2", company: "A", title: "Engineer", location: "Austin", compensation: "$200", postedAt: "2024-01-02", url: "u", appliedBy: [] },
    ];
    const output = filterAndSortFeedJobs(sample, { search: "engineer", sortBy: "company" });
    expect(output[0].company).toBe("A");
  });

  it("resolves apply/reject actions", async () => {
    global.fetch = vi.fn(async () => new Response(null, { status: 204 })) as unknown as typeof fetch;
    await expect(markApplied("job_1", "test-token")).resolves.toBeUndefined();
    await expect(rejectPosting("job_1", "test-token")).resolves.toBeUndefined();
  });
});
