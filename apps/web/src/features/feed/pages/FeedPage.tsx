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
export function FeedPage(){const {theme,toggleTheme}=useTheme();const {filters,setSearch,setSortBy}=useFeedFilters();const debouncedSearch=useDebouncedValue(filters.search,250);const resolvedFilters=useMemo(()=>({...filters,search:debouncedSearch}),[filters,debouncedSearch]);const {jobs,loading,setJobs}=useFeedQuery(resolvedFilters);async function onApply(id:string){await markApplied(id);setJobs((p)=>p.filter((j)=>j.id!==id));}async function onReject(id:string){await rejectPosting(id);setJobs((p)=>p.filter((j)=>j.id!==id));}return <PageShell title="Welcome" actions={<div style={{display:"flex",gap:"0.5rem"}}><UserSearchBar/><Button onClick={toggleTheme}>Theme: {theme}</Button></div>}><FeedFilters search={filters.search} sortBy={filters.sortBy} onSearch={setSearch} onSort={setSortBy}/>{loading?<p>Loading feed...</p>:<JobList jobs={jobs} onApply={onApply} onReject={onReject}/>}</PageShell>;}
