import React from "react";
import { Button } from "../../../shared/components/Button";
import { formatDate } from "../../../shared/utils/date";
import { formatCompensation } from "../../../shared/utils/money";
import { AppliedByBadge } from "./AppliedByBadge";
import { JobPosting } from "../model/feed.types";

export function JobListItem({
  job,
  onApply,
  onReject,
}: {
  job: JobPosting;
  onApply: (id: string) => void;
  onReject: (id: string) => void;
}) {
  const hasPostingUrl = job.url.trim().length > 0;

  return (
    <article
      style={{
        border: "1px solid var(--border)",
        borderRadius: "var(--radius-md)",
        padding: "0.95rem 1rem",
        marginBottom: "0.85rem",
        background: "var(--surface)",
        boxShadow: "var(--shadow-sm)",
      }}
    >
      <div style={{ display: "flex", justifyContent: "space-between", gap: "0.75rem" }}>
        <div>
          <strong style={{ fontSize: "1rem" }}>{job.company}</strong> | {job.title} | {formatCompensation(job.compensation)} | {job.location}
          <div>
            <small>Posted: {formatDate(job.postedAt)}</small>
          </div>
          {job.appliedAt ? (
            <div>
              <small>Applied on: {formatDate(job.appliedAt)}</small>
            </div>
          ) : null}
          <AppliedByBadge users={job.appliedBy} />
          <div style={{ marginTop: "0.35rem" }}>
            {hasPostingUrl ? (
              <a href={job.url} target="_blank" rel="noreferrer" style={{ color: "var(--primary)" }}>
                View posting
              </a>
            ) : (
              <small style={{ color: "var(--muted)" }}>Posting link unavailable</small>
            )}
          </div>
        </div>
        <div style={{ display: "flex", gap: "0.5rem" }}>
          <Button onClick={() => onApply(job.id)}>Apply</Button>
          <Button variant="danger" onClick={() => onReject(job.id)}>
            Reject
          </Button>
        </div>
      </div>
    </article>
  );
}
