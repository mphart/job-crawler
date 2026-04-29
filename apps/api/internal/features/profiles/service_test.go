package profiles

import (
	"testing"

	"job-crawler/apps/api/internal/platform/db"
)

func TestProfilePrivacyAndAppliedHistory(t *testing.T) {
	store := db.NewStore()
	feedDecision := db.FeedDecision{UserID: "u_1", JobID: "job_1", DecisionType: "APPLIED", DecisionAt: "2026-01-01T00:00:00Z"}
	store.Decisions = append(store.Decisions, feedDecision)
	store.Users["u_1"] = db.User{ID: "u_1", Email: "mason@example.com", Username: "mason", IsPrivate: true}

	svc := Service{Store: InMemoryStore{Inner: store}}

	own, ok, err := svc.Get("u_1", "u_1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatalf("expected own profile")
	}
	if own.TotalApplied < 1 {
		t.Fatalf("expected applied count from decisions")
	}

	other, ok, err := svc.Get("u_2", "u_1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatalf("expected profile lookup to succeed")
	}
	if !other.IsPrivate {
		t.Fatalf("expected profile to remain private")
	}
}
