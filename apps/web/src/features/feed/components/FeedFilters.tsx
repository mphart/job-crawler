import { ChangeEvent } from "react";
import { Select } from "../../../shared/components/Select";
import { TextInput } from "../../../shared/components/TextInput";
import { FeedSort } from "../model/feed.types";
export function FeedFilters({search,sortBy,onSearch,onSort}:{search:string;sortBy:FeedSort;onSearch:(value:string)=>void;onSort:(value:FeedSort)=>void;}){return <section style={{display:"grid",gap:"0.5rem",marginBottom:"1rem"}}><TextInput value={search} onChange={(e)=>onSearch(e.target.value)} placeholder="Search jobs" aria-label="search jobs"/><Select value={sortBy} onChange={(event:ChangeEvent<HTMLSelectElement>)=>onSort(event.target.value as FeedSort)} aria-label="sort jobs"><option value="newest">Newest</option><option value="company">Company</option><option value="title">Title</option><option value="location">Location</option><option value="money">Compensation</option></Select></section>;}
