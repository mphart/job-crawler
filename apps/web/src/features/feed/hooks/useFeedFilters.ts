import { useState } from "react";
import { FeedFilters, FeedSort } from "../model/feed.types";
export function useFeedFilters(initial:FeedFilters={search:"",sortBy:"newest"}){const [filters,setFilters]=useState<FeedFilters>(initial);return {filters,setSearch:(search:string)=>setFilters((p)=>({...p,search})),setSortBy:(sortBy:FeedSort)=>setFilters((p)=>({...p,sortBy}))};}
