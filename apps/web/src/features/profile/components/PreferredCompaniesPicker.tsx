import React, { useState } from "react";
import { Button } from "../../../shared/components/Button";
import { TextInput } from "../../../shared/components/TextInput";
import { CompanySearchResult, searchCompanies } from "../api/profile.api";

type Props = {
  token?: string;
  selected: string[];
  onChange: (next: string[]) => void;
};

export function PreferredCompaniesPicker({ token, selected, onChange }: Props) {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<CompanySearchResult[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSearch() {
    if (query.trim().length < 2) {
      setError("Enter at least 2 characters to search companies.");
      setResults([]);
      return;
    }
    setError(null);
    setIsSearching(true);
    try {
      setResults(await searchCompanies(query.trim(), token));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to search companies right now.");
      setResults([]);
    } finally {
      setIsSearching(false);
    }
  }

  function addCompany(company: string) {
    if (selected.some((current) => current.toLowerCase() === company.toLowerCase())) {
      return;
    }
    onChange([...selected, company]);
  }

  function removeCompany(company: string) {
    onChange(selected.filter((current) => current.toLowerCase() !== company.toLowerCase()));
  }

  return (
    <section className="ui-card" style={{ display: "grid", gap: "0.75rem", padding: "0.9rem" }}>
      <h3 style={{ margin: 0 }}>Preferred Companies</h3>
      <p style={{ margin: 0, color: "var(--muted)" }}>
        Search and add verified companies. Roles from these companies are prioritized near the top of your feed.
      </p>
      <div style={{ display: "flex", gap: "0.5rem", alignItems: "center" }}>
        <TextInput value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Search companies..." aria-label="Search companies" />
        <Button type="button" variant="secondary" onClick={onSearch} disabled={isSearching}>{isSearching ? "Searching..." : "Search"}</Button>
      </div>
      {error ? <small style={{ color: "var(--danger)" }}>{error}</small> : null}
      {results.length > 0 ? (
        <div style={{ display: "grid", gap: "0.4rem" }}>
          {results.map((company) => (
            <div key={company.name} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", gap: "0.5rem" }}>
              <span>{company.name} {company.isVerified ? "✓" : ""}</span>
              <Button type="button" onClick={() => addCompany(company.name)} disabled={selected.some((item) => item.toLowerCase() === company.name.toLowerCase())}>
                Add
              </Button>
            </div>
          ))}
        </div>
      ) : query.trim().length >= 2 && !isSearching && !error ? (
        <small style={{ color: "var(--muted)" }}>No verified companies found for that search yet.</small>
      ) : null}
      {selected.length > 0 ? (
        <div style={{ display: "flex", flexWrap: "wrap", gap: "0.45rem" }}>
          {selected.map((company) => (
            <button
              key={company}
              type="button"
              onClick={() => removeCompany(company)}
              style={{
                border: "1px solid var(--border)",
                borderRadius: "999px",
                background: "var(--surface-elevated)",
                color: "var(--text)",
                padding: "0.2rem 0.55rem",
                cursor: "pointer",
              }}
              aria-label={`Remove ${company}`}
            >
              {company} ×
            </button>
          ))}
        </div>
      ) : (
        <small style={{ color: "var(--muted)" }}>No companies selected yet.</small>
      )}
    </section>
  );
}
