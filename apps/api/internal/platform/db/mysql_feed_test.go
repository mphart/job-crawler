package db

import "testing"

func TestRecommendationScore_PrefersKeywordMatches(t *testing.T) {
	jobA := JobPosting{Company: "Acme", Title: "Mechanical Engineer", Location: "Austin"}
	jobB := JobPosting{Company: "Acme", Title: "Software Engineer", Location: "Austin"}
	keywords := []string{"mechanical engineer"}

	scoreA := recommendationScore(jobA, keywords, nil)
	scoreB := recommendationScore(jobB, keywords, nil)
	if scoreA <= scoreB {
		t.Fatalf("expected mechanical role to score higher (%d <= %d)", scoreA, scoreB)
	}
}

func TestRecommendationScore_PrefersPreferredCompany(t *testing.T) {
	jobA := JobPosting{Company: "Stripe", Title: "Backend Engineer", Location: "Remote"}
	jobB := JobPosting{Company: "OtherCo", Title: "Backend Engineer", Location: "Remote"}
	preferredCompanies := []string{"stripe"}

	scoreA := recommendationScore(jobA, nil, preferredCompanies)
	scoreB := recommendationScore(jobB, nil, preferredCompanies)
	if scoreA <= scoreB {
		t.Fatalf("expected preferred company to score higher (%d <= %d)", scoreA, scoreB)
	}
}

func TestSplitKeywords_TrimsAndNormalizes(t *testing.T) {
	got := splitKeywords("  Mechanical Engineer, CAD , ,Austin ")
	if len(got) != 3 {
		t.Fatalf("expected 3 keywords, got %v", got)
	}
	if got[0] != "mechanical engineer" {
		t.Fatalf("expected normalized lowercase keyword, got %s", got[0])
	}
}
