import { useEffect, useState } from "react";
import { fetchFeed } from "../api/feed.api";
import { FeedFilters, JobPosting } from "../model/feed.types";

export function useFeedQuery(filters: FeedFilters, token: string | null) {
  const [jobs, setJobs] = useState<JobPosting[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!token) {
      setJobs([]);
      setLoading(false);
      setError("Missing session token.");
      return;
    }

    const controller = new AbortController();
    setLoading(true);
    setError(null);

    fetchFeed(filters, token, controller.signal)
      .then((data) => {
        if (!controller.signal.aborted) {
          setJobs(data);
        }
      })
      .catch((err: unknown) => {
        if (!controller.signal.aborted) {
          setError(err instanceof Error ? err.message : "Failed to load feed.");
        }
      })
      .finally(() => {
        if (!controller.signal.aborted) {
          setLoading(false);
        }
      });

    return () => controller.abort();
  }, [filters, token]);

  return { jobs, loading, error, setJobs };
}
