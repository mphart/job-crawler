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

    let activeController: AbortController | null = null;
    let pollTimeout: ReturnType<typeof setTimeout> | null = null;
    let firstLoad = true;

    const load = async () => {
      activeController?.abort();
      const controller = new AbortController();
      activeController = controller;

      if (firstLoad) {
        setLoading(true);
      }
      setError(null);

      try {
        const data = await fetchFeed(filters, token, controller.signal);
        if (!controller.signal.aborted) {
          setJobs(data);
        }
      } catch (err: unknown) {
        if (!controller.signal.aborted) {
          setError(err instanceof Error ? err.message : "Failed to load feed.");
        }
      } finally {
        if (!controller.signal.aborted && firstLoad) {
          setLoading(false);
        }
        firstLoad = false;
      }

      if (!controller.signal.aborted) {
        pollTimeout = setTimeout(load, 10000);
      }
    };

    void load();

    return () => {
      activeController?.abort();
      if (pollTimeout) {
        clearTimeout(pollTimeout);
      }
    };
  }, [filters, token]);

  return { jobs, loading, error, setJobs };
}
