import React from "react";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { TextInput } from "../../../shared/components/TextInput";
import { searchUsers } from "../api/users.api";
import { UserSummary } from "../model/users.types";

export function UserSearchBar({ token }: { token: string }) {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<UserSummary[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!query.trim()) {
      setResults([]);
      setError(null);
      return;
    }

    const controller = new AbortController();
    setIsLoading(true);
    setError(null);

    searchUsers(query, token, controller.signal)
      .then((users) => {
        if (!controller.signal.aborted) {
          setResults(users);
        }
      })
      .catch((err: unknown) => {
        if (!controller.signal.aborted) {
          setError(err instanceof Error ? err.message : "Failed to search users.");
          setResults([]);
        }
      })
      .finally(() => {
        if (!controller.signal.aborted) {
          setIsLoading(false);
        }
      });

    return () => controller.abort();
  }, [query, token]);

  return (
    <div style={{ position: "relative", minWidth: 220 }}>
      <TextInput value={query} onChange={(e) => setQuery(e.target.value)} placeholder="Find users" />
      {isLoading ? <small>Searching...</small> : null}
      {error ? <small style={{ color: "var(--danger)" }}>{error}</small> : null}
      {results.length > 0 ? (
        <div
          style={{
            position: "absolute",
            top: "105%",
            left: 0,
            right: 0,
            border: "1px solid var(--border)",
            borderRadius: 8,
            background: "var(--surface)",
            zIndex: 10,
          }}
        >
          {results.map((user) => (
            <Link key={user.id} to={`/profile/${user.id}`} style={{ display: "block", padding: "0.5rem" }}>
              {user.username} ({user.totalApplied} applied)
            </Link>
          ))}
        </div>
      ) : null}
    </div>
  );
}
