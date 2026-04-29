export type FeedSort="newest"|"company"|"title"|"location"|"money";
export type AppliedBy={userId:string;username:string};
export type JobPosting={id:string;company:string;title:string;location:string;compensation:string;postedAt:string;appliedAt?:string;url:string;appliedBy:AppliedBy[]};
export type FeedFilters={search:string;sortBy:FeedSort};
