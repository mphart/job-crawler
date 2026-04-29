import { EmptyState } from "../../../shared/components/EmptyState";
import { JobPosting } from "../model/feed.types";
import { JobListItem } from "./JobListItem";
export function JobList({jobs,onApply,onReject}:{jobs:JobPosting[];onApply:(id:string)=>void;onReject:(id:string)=>void;}){if(jobs.length===0){return <EmptyState title="No jobs found" description="Try broadening your search or changing sort options."/>;}return <section>{jobs.map((job)=><JobListItem key={job.id} job={job} onApply={onApply} onReject={onReject}/>)}</section>;}
