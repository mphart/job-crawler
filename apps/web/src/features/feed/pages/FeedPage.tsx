import React from "react";
import { useMemo } from "react";
import { PageShell } from "../../../shared/components/PageShell";
import { useDebouncedValue } from "../../../shared/hooks/useDebouncedValue";
import { Button } from "../../../shared/components/Button";
import { useTheme } from "../../../shared/hooks/useTheme";
import { markApplied, rejectPosting } from "../api/feed.api";
import { FeedFilters } from "../components/FeedFilters";
import { JobList } from "../components/JobList";
import { useFeedFilters } from "../hooks/useFeedFilters";
import { useFeedQuery } from "../hooks/useFeedQuery";
import { UserSearchBar } from "../../user-discovery/components/UserSearchBar";
import { EmptyState } from "../../../shared/components/EmptyState";
import { useSession } from "../../auth/hooks/useSession";

export function FeedPage() {
  const { theme, toggleTheme } = useTheme();
  const { user, signOut } = useSession();
  const { filters, setSearch, setSortBy } = useFeedFilters();
  const debouncedSearch = useDebouncedValue(filters.search, 250);
  const resolvedFilters = useMemo(() => ({ ...filters, search: debouncedSearch }), [filters, debouncedSearch]);
  const { jobs, loading, error, setJobs } = useFeedQuery(resolvedFilters, user?.token ?? null);

  async function onApply(id: string) {
    if (!user) return;
    await markApplied(id, user.token);
    setJobs((prev) => prev.filter((j) => j.id !== id));
  }

  async function onReject(id: string) {
    if (!user) return;
    await rejectPosting(id, user.token);
    setJobs((prev) => prev.filter((j) => j.id !== id));
  }

  const actions = (
    <div style={{ display: "flex", gap: "0.5rem" }}>
      {user ? <UserSearchBar token={user.token} /> : null}
      <Button onClick={toggleTheme}>Theme: {theme}</Button>
      <Button variant="secondary" onClick={signOut}>Sign out</Button>
    </div>
  );

  return (
    <PageShell title="Welcome" actions={actions}>
      <FeedFilters search={filters.search} sortBy={filters.sortBy} onSearch={setSearch} onSort={setSortBy} />
      {loading ? <p>Loading feed...</p> : null}
      {error ? <EmptyState title="Unable to load feed" description={error} /> : null}
      {!loading && !error ? <JobList jobs={jobs} onApply={onApply} onReject={onReject} /> : null}
    </PageShell>
  );
}
