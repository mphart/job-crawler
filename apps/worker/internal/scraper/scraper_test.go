package scraper

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestCanonicalizeURL(t *testing.T) {
	got := canonicalizeURL("indeed", "https://www.indeed.com/viewjob?jk=abc123&from=search")
	if got != "https://www.indeed.com/viewjob?jk=abc123" {
		t.Fatalf("unexpected canonical url: %s", got)
	}
}

func TestNormalizeCompensationVariants(t *testing.T) {
	input := "$120k - $140k per year"
	got := normalizeCompensation(input)
	if got == "" {
		t.Fatalf("expected normalized compensation")
	}
	if !strings.Contains(got, "/YEAR") {
		t.Fatalf("expected yearly suffix, got %q", got)
	}
}

func TestParsePostedAtRelativeDate(t *testing.T) {
	got := parsePostedAt("3 days ago")
	parsed, err := time.Parse(time.RFC3339, got)
	if err != nil {
		t.Fatalf("expected RFC3339 date, got %q", got)
	}
	if parsed.After(time.Now().UTC()) {
		t.Fatalf("expected past timestamp, got %q", got)
	}
}

func TestExtractLinkedInPostedAtPrefersDatetime(t *testing.T) {
	card := `<li><time class="job-search-card__listdate" datetime="2026-05-01T12:00:00Z">4 days ago</time></li>`
	got := extractLinkedInPostedAt(card)
	if got != "2026-05-01T12:00:00Z" {
		t.Fatalf("expected datetime to be used, got %q", got)
	}
}

func TestApplyQualityGatesRejectsWeakRows(t *testing.T) {
	rows := []ScrapedJob{
		{Source: "linkedin", Company: "Unknown LinkedIn Company", Title: "Engineer", Location: "Remote", URL: "https://example.com/a"},
		{Source: "linkedin", Company: "Acme", Title: "Engineer", Location: "Remote", Compensation: "$120,000", URL: "https://example.com/b"},
	}
	out := applyQualityGates(rows)
	if len(out) != 1 {
		t.Fatalf("expected one high-quality row, got %d", len(out))
	}
	if out[0].URL != "https://example.com/b" {
		t.Fatalf("unexpected row survived gate: %+v", out[0])
	}
}

func TestLinkedInTieredExtractionFromFixtureLikeHTML(t *testing.T) {
	data, err := os.ReadFile("testdata/linkedin_sample.html")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}
	card := string(data)
	title := firstMatch(card, `(?s)base-search-card__title[^>]*>\s*(.*?)\s*<`)
	company := firstMatch(card, `(?s)base-search-card__subtitle[^>]*>\s*(.*?)\s*<`, `(?s)job-search-card__subtitle-link[^>]*>\s*(.*?)\s*<`, `(?s)job-search-card__subtitle-link[^>]*>(.*?)</h4>`)
	comp := firstMatch(card, `(?s)(?:job-search-card__salary-info|base-search-card__metadata)[^>]*>\s*(.*?)\s*<`)
	if title == "" || company == "" || comp == "" {
		t.Fatalf("expected title/company/comp extraction from linkedin fixture, got %q %q %q", title, company, comp)
	}
}

func TestGlassdoorTieredExtractionFromFixtureLikeHTML(t *testing.T) {
	data, err := os.ReadFile("testdata/glassdoor_sample.html")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}
	card := string(data)
	title := firstMatch(card, `(?s)data-test="job-title"[^>]*>\s*(.*?)\s*<`)
	company := firstMatch(card, `(?s)data-test="employer-name"[^>]*>\s*(.*?)\s*<`)
	comp := firstMatch(card, `(?s)data-test="detailSalary"[^>]*>\s*(.*?)\s*<`)
	if title == "" || company == "" || comp == "" {
		t.Fatalf("expected title/company/comp extraction from glassdoor fixture, got %q %q %q", title, company, comp)
	}
}

func TestIndeedTieredExtractionFromFixtureLikeHTML(t *testing.T) {
	data, err := os.ReadFile("testdata/indeed_sample.html")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}
	card := string(data)
	title := firstMatch(card, `(?s)jobTitle[^>]*>\s*(.*?)\s*<`)
	company := firstMatch(card, `(?s)data-testid="company-name"[^>]*>\s*(.*?)\s*<`)
	comp := firstMatch(card, `(?s)class="[^"]*salary-snippet[^"]*"[^>]*>\s*(.*?)\s*<`)
	if title == "" || company == "" || comp == "" {
		t.Fatalf("expected title/company/comp extraction from indeed fixture, got %q %q %q", title, company, comp)
	}
}
