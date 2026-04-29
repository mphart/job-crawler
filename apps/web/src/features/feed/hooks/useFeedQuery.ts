import { useEffect, useState } from "react";
import { fetchFeed } from "../api/feed.api";
import { FeedFilters, JobPosting } from "../model/feed.types";
export function useFeedQuery(filters:FeedFilters){const [jobs,setJobs]=useState<JobPosting[]>([]);const [loading,setLoading]=useState(true);useEffect(()=>{setLoading(true);fetchFeed(filters).then((data)=>{setJobs(data);setLoading(false);});},[filters]);return {jobs,loading,setJobs};}
