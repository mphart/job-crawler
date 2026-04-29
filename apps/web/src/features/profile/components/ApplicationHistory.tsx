import { EmptyState } from "../../../shared/components/EmptyState";
import { formatDate } from "../../../shared/utils/date";
import { JobPosting } from "../../feed/model/feed.types";
export function ApplicationHistory({jobs}:{jobs:JobPosting[]}){if(jobs.length===0){return <EmptyState title="No applications yet" description="Applied jobs will appear here."/>;}return <div>{jobs.map((job)=><article key={job.id} style={{borderBottom:"1px solid var(--border)",padding:"0.5rem 0"}}><strong>{job.company}</strong> - {job.title}<div><small>Applied on: {formatDate(job.appliedAt??job.postedAt)}</small></div></article>)}</div>;}
