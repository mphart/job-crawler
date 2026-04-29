import React from "react";
import { AppliedBy } from "../model/feed.types";

export function AppliedByBadge({ users }: { users: AppliedBy[] | undefined }) {
  if (!users || users.length === 0) return null;
  return <small style={{ color: "var(--muted)" }}>Also applied: {users.map((u) => u.username).join(", ")}</small>;
}
