package db

import "testing"

func TestRecommendationScore_PrefersKeywordMatches(t *testing.T) {
	jobA := JobPosting{Company: "Acme", Title: "Mechanical Engineer", Location: "Austin"}
	jobB := JobPosting{Company: "Acme", Title: "Software Engineer", Location: "Austin"}
	keywords := []string{"mechanical engineer"}

	scoreA := recommendationScore(jobA, keywords)
	scoreB := recommendationScore(jobB, keywords)
	if scoreA <= scoreB {
		t.Fatalf("expected mechanical role to score higher (%d <= %d)", scoreA, scoreB)
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
