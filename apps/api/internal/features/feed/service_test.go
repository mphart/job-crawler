package feed

import (
	"testing"

	"job-crawler/apps/api/internal/platform/db"
)

func TestApplyAndRejectHideJobsFromFeed(t *testing.T) {
	store := db.NewStore()
	store.Jobs["job_real_1"] = db.JobPosting{
		ID:           "job_real_1",
		Company:      "TechOne",
		Title:        "Software Engineer",
		Location:     "Remote",
		Compensation: "$180k-$210k",
		PostedAt:     "2026-01-02T00:00:00Z",
		URL:          "https://jobs.techone.com/ml-eng",
	}
	store.Jobs["job_real_2"] = db.JobPosting{
		ID:           "job_real_2",
		Company:      "TechTwo",
		Title:        "Frontend Engineer",
		Location:     "Remote",
		Compensation: "$160k-$190k",
		PostedAt:     "2026-01-01T00:00:00Z",
		URL:          "https://jobs.techtwo.com/backend",
	}
	svc := Service{Store: InMemoryStore{Inner: store}}

	initial, err := svc.List("u_1", "", "newest")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(initial) == 0 {
		t.Fatalf("expected jobs in feed")
	}

	if err := svc.Apply("u_1", initial[0].ID); err != nil {
		t.Fatalf("apply failed: %v", err)
	}
	afterApply, err := svc.List("u_1", "", "newest")
	if err != nil {
		t.Fatalf("list after apply failed: %v", err)
	}
	if len(afterApply) != len(initial)-1 {
		t.Fatalf("expected applied job removed from feed")
	}

	if len(afterApply) > 0 {
		if err := svc.Reject("u_1", afterApply[0].ID); err != nil {
			t.Fatalf("reject failed: %v", err)
		}
		afterReject, err := svc.List("u_1", "", "newest")
		if err != nil {
			t.Fatalf("list after reject failed: %v", err)
		}
		if len(afterReject) != len(afterApply)-1 {
			t.Fatalf("expected rejected job removed from feed")
		}
	}
}
