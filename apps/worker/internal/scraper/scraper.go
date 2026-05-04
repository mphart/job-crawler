package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type ScrapedJob struct {
	Source       string `json:"source"`
	ExternalID   string `json:"externalId"`
	Company      string `json:"company"`
	Title        string `json:"title"`
	Location     string `json:"location"`
	Compensation string `json:"compensation"`
	PostedAt     string `json:"postedAt"`
	URL          string `json:"url"`
}

func Scrape(ctx context.Context, client *http.Client, keywords []string) ([]ScrapedJob, error) {
	var out []ScrapedJob
	if len(keywords) == 0 {
		keywords = []string{"software engineer"}
	}

	// Hardcoded source priority: LinkedIn, Indeed, Glassdoor, then Greenhouse.
	for _, keyword := range keywords {
		linkedin, err := scrapeLinkedIn(ctx, client, keyword)
		if err == nil {
			out = append(out, linkedin...)
		}
		indeed, err := scrapeIndeed(ctx, client, keyword)
		if err == nil {
			out = append(out, indeed...)
		}
		glassdoor, err := scrapeGlassdoor(ctx, client, keyword)
		if err == nil {
			out = append(out, glassdoor...)
		}
		greenhouse, err := scrapeGreenhouse(ctx, client, keyword)
		if err == nil {
			out = append(out, greenhouse...)
		}
	}

	out = applyQualityGates(out)
	logMissingFieldRates(out)

	if len(out) == 0 {
		return nil, fmt.Errorf("no jobs scraped from configured sources (linkedin/indeed/glassdoor/greenhouse)")
	}
	return dedupeByURL(out), nil
}

func scrapeLinkedIn(ctx context.Context, client *http.Client, keyword string) ([]ScrapedJob, error) {
	q := url.Values{}
	q.Set("keywords", keyword)
	q.Set("location", "United States")
	q.Set("start", "0")

	endpoint := "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("linkedin status %d: %s", res.StatusCode, string(body))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	body := string(bodyBytes)
	cardRegex := regexp.MustCompile(`(?s)<li[^>]*>.*?</li>`)

	cards := cardRegex.FindAllString(body, -1)
	jobs := make([]ScrapedJob, 0, len(cards))
	for _, card := range cards {
		title := firstMatch(card,
			`(?s)base-search-card__title[^>]*>\s*(.*?)\s*<`,
			`(?s)data-tracking-control-name="public_jobs_jserp-result_search-card"[^>]*>\s*(.*?)\s*<`,
		)
		company := firstMatch(card,
			`(?s)base-search-card__subtitle[^>]*>\s*(.*?)\s*<`,
			`(?s)job-search-card__subtitle-link[^>]*>\s*(.*?)\s*<`,
			`(?s)job-search-card__subtitle-link[^>]*>(.*?)</h4>`,
			`(?s)hidden-nested-link[^>]*>\s*(.*?)\s*<`,
		)
		location := firstMatch(card,
			`(?s)job-search-card__location[^>]*>\s*(.*?)\s*<`,
			`(?s)job-search-card__listdate[^>]*>\s*(.*?)\s*<`,
		)
		compensation := firstMatch(card,
			`(?s)(?:job-search-card__salary-info|base-search-card__metadata)[^>]*>\s*(.*?)\s*<`,
			`(?s)salary[^>]*>\s*(.*?)\s*<`,
		)
		if strings.TrimSpace(compensation) == "" {
			compensation = captureMoneyWindow(card, 180)
		}
		link := firstMatch(card, `href="([^"]*linkedin\.com/jobs/view/[^"]+)"`)
		if strings.TrimSpace(title) == "" || strings.TrimSpace(link) == "" {
			continue
		}
		if !matchesKeyword(title, keyword) {
			continue
		}
		externalID := firstMatch(card, `data-entity-urn="urn:li:jobPosting:(\d+)"`)
		jobs = append(jobs, ScrapedJob{
			Source:       "linkedin",
			ExternalID:   externalID,
			Company:      normalizeCompany(html.UnescapeString(clean(stripTags(company))), "linkedin"),
			Title:        html.UnescapeString(clean(title)),
			Location:     html.UnescapeString(clean(location)),
			Compensation: normalizeCompensation(html.UnescapeString(clean(stripTags(compensation)))),
			PostedAt:     time.Now().Format(time.RFC3339),
			URL:          canonicalizeURL("linkedin", html.UnescapeString(clean(link))),
		})
	}
	return jobs, nil
}

func scrapeIndeed(ctx context.Context, client *http.Client, keyword string) ([]ScrapedJob, error) {
	endpoint := "https://www.indeed.com/jobs?q=" + url.QueryEscape(keyword) + "&l=United+States"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("indeed status %d: %s", res.StatusCode, string(body))
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	body := string(bodyBytes)
	cardRegex := regexp.MustCompile(`(?s)<div[^>]*class="[^"]*job_seen_beacon[^"]*"[^>]*>.*?</div>`)

	cards := cardRegex.FindAllString(body, -1)
	jobs := make([]ScrapedJob, 0, len(cards))
	for _, card := range cards {
		title := firstMatch(card,
			`(?s)jobTitle[^>]*>\s*(.*?)\s*<`,
			`(?s)data-testid="job-title"[^>]*>\s*(.*?)\s*<`,
		)
		link := firstMatch(card, `href="(/pagead/clk[^"]+|/rc/clk[^"]+|/viewjob\?jk=[^"]+)"`)
		if strings.TrimSpace(title) == "" || strings.TrimSpace(link) == "" {
			continue
		}
		if !matchesKeyword(title, keyword) {
			continue
		}
		fullLink := link
		if strings.HasPrefix(link, "/") {
			fullLink = "https://www.indeed.com" + link
		}
		externalID := firstMatch(fullLink, `jk=([a-zA-Z0-9]+)`)
		company := firstMatch(card,
			`(?s)data-testid="company-name"[^>]*>\s*(.*?)\s*<`,
			`(?s)class="[^"]*companyName[^"]*"[^>]*>\s*(.*?)\s*<`,
			`(?s)data-testid="company-title"[^>]*>\s*(.*?)\s*<`,
		)
		location := firstMatch(card,
			`(?s)data-testid="text-location"[^>]*>\s*(.*?)\s*<`,
			`(?s)companyLocation[^>]*>\s*(.*?)\s*<`,
		)
		compensation := firstMatch(card,
			`(?s)class="[^"]*salary-snippet[^"]*"[^>]*>\s*(.*?)\s*<`,
			`(?s)data-testid="attribute_snippet_testid"[^>]*>\s*(.*?)\s*<`,
		)
		if strings.TrimSpace(compensation) == "" {
			compensation = captureMoneyWindow(card, 180)
		}
		jobs = append(jobs, ScrapedJob{
			Source:       "indeed",
			ExternalID:   externalID,
			Company:      normalizeCompany(html.UnescapeString(clean(stripTags(company))), "indeed"),
			Title:        html.UnescapeString(clean(title)),
			Location:     html.UnescapeString(clean(stripTags(location))),
			Compensation: normalizeCompensation(html.UnescapeString(clean(stripTags(compensation)))),
			PostedAt:     time.Now().Format(time.RFC3339),
			URL:          canonicalizeURL("indeed", html.UnescapeString(clean(fullLink))),
		})
	}
	return jobs, nil
}

func scrapeGlassdoor(ctx context.Context, client *http.Client, keyword string) ([]ScrapedJob, error) {
	endpoint := "https://www.glassdoor.com/Job/jobs.htm?sc.keyword=" + url.QueryEscape(keyword)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("glassdoor status %d: %s", res.StatusCode, string(body))
	}
	bodyBytes, _ := io.ReadAll(res.Body)
	body := string(bodyBytes)

	cardRegex := regexp.MustCompile(`(?s)<li[^>]*data-test="jobListing"[^>]*>.*?</li>`)
	cards := cardRegex.FindAllString(body, -1)
	jobs := make([]ScrapedJob, 0, len(cards))
	for _, card := range cards {
		title := firstMatch(card,
			`(?s)data-test="job-title"[^>]*>\s*(.*?)\s*<`,
			`(?s)jobTitle[^>]*>\s*(.*?)\s*<`,
		)
		if strings.TrimSpace(title) == "" || !matchesKeyword(title, keyword) {
			continue
		}
		company := firstMatch(card,
			`(?s)data-test="employer-name"[^>]*>\s*(.*?)\s*<`,
			`(?s)employerName[^>]*>\s*(.*?)\s*<`,
		)
		location := firstMatch(card,
			`(?s)data-test="emp-location"[^>]*>\s*(.*?)\s*<`,
			`(?s)location[^>]*>\s*(.*?)\s*<`,
		)
		compensation := firstMatch(card,
			`(?s)data-test="detailSalary"[^>]*>\s*(.*?)\s*<`,
			`(?s)salaryEstimate[^>]*>\s*(.*?)\s*<`,
		)
		if strings.TrimSpace(compensation) == "" {
			compensation = captureMoneyWindow(card, 200)
		}
		link := firstMatch(card, `href="(/job-listing/[^"]+)"`)
		if strings.TrimSpace(link) == "" {
			continue
		}
		fullLink := link
		if strings.HasPrefix(link, "/") {
			fullLink = "https://www.glassdoor.com" + link
		}
		externalID := firstMatch(fullLink, `jobListingId=([a-zA-Z0-9]+)`)
		jobs = append(jobs, ScrapedJob{
			Source:       "glassdoor",
			ExternalID:   externalID,
			Company:      normalizeCompany(html.UnescapeString(clean(stripTags(company))), "glassdoor"),
			Title:        html.UnescapeString(clean(stripTags(title))),
			Location:     html.UnescapeString(clean(stripTags(location))),
			Compensation: normalizeCompensation(html.UnescapeString(clean(stripTags(compensation)))),
			PostedAt:     time.Now().Format(time.RFC3339),
			URL:          canonicalizeURL("glassdoor", html.UnescapeString(clean(fullLink))),
		})
	}
	return jobs, nil
}

func scrapeGreenhouse(ctx context.Context, client *http.Client, keyword string) ([]ScrapedJob, error) {
	endpoint := "https://boards-api.greenhouse.io/v1/boards/stripe/jobs?content=true"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "job-crawler-worker/1.0")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("greenhouse status %d: %s", res.StatusCode, string(body))
	}

	type greenhouseResp struct {
		Jobs []struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			AbsoluteURL string `json:"absolute_url"`
			Location    struct {
				Name string `json:"name"`
			} `json:"location"`
		} `json:"jobs"`
	}
	var payload greenhouseResp
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	jobs := make([]ScrapedJob, 0, len(payload.Jobs))
	for _, j := range payload.Jobs {
		if strings.TrimSpace(j.AbsoluteURL) == "" || strings.TrimSpace(j.Title) == "" {
			continue
		}
		if !matchesKeyword(j.Title, keyword) {
			continue
		}
		jobs = append(jobs, ScrapedJob{
			Source:       "greenhouse",
			ExternalID:   fmt.Sprintf("%d", j.ID),
			Company:      normalizeCompany("Stripe", "greenhouse"),
			Title:        j.Title,
			Location:     j.Location.Name,
			Compensation: "",
			PostedAt:     time.Now().Format(time.RFC3339),
			URL:          canonicalizeURL("greenhouse", j.AbsoluteURL),
		})
	}
	return jobs, nil
}

func matchesKeyword(title, keyword string) bool {
	return strings.Contains(strings.ToLower(title), strings.ToLower(strings.TrimSpace(keyword)))
}

func capture(re *regexp.Regexp, source string) string {
	match := re.FindStringSubmatch(source)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func firstMatch(source string, patterns ...string) string {
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		value := capture(re, source)
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func clean(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\t", " ")
	value = strings.TrimSpace(value)
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(value, " ")
}

func stripTags(value string) string {
	tagRegex := regexp.MustCompile(`(?s)<[^>]+>`)
	return tagRegex.ReplaceAllString(value, " ")
}

func normalizeCompany(company, source string) string {
	trimmed := strings.TrimSpace(company)
	if trimmed != "" {
		return trimmed
	}
	switch source {
	case "linkedin":
		return "Unknown LinkedIn Company"
	case "indeed":
		return "Unknown Indeed Company"
	case "greenhouse":
		return "Unknown Greenhouse Company"
	case "glassdoor":
		return "Unknown Glassdoor Company"
	default:
		return "Unknown Company"
	}
}

func normalizeCompensation(value string) string {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		return ""
	}
	trimmed = strings.ReplaceAll(trimmed, "per year", "/year")
	trimmed = strings.ReplaceAll(trimmed, "per hour", "/hour")
	trimmed = strings.ReplaceAll(trimmed, "yr", "/year")
	trimmed = strings.ReplaceAll(trimmed, "hr", "/hour")
	trimmed = strings.ReplaceAll(trimmed, "k", "000")
	money := captureMoneyFromText(trimmed)
	if money != "" {
		return strings.ToUpper(strings.ReplaceAll(money, "/year", "/YEAR"))
	}
	return strings.TrimSpace(value)
}

func captureMoneyFromText(value string) string {
	moneyRegex := regexp.MustCompile(`(?i)\$[\d,]+(?:\s*-\s*\$[\d,]+)?(?:\s*(?:/year|/hour|year|hour|yr|hr|k))?`)
	match := moneyRegex.FindString(value)
	return strings.TrimSpace(match)
}

func captureMoneyWindow(value string, around int) string {
	lowered := strings.ToLower(value)
	i := strings.Index(lowered, "$")
	if i == -1 {
		return captureMoneyFromText(value)
	}
	start := i - around
	if start < 0 {
		start = 0
	}
	end := i + around
	if end > len(value) {
		end = len(value)
	}
	return captureMoneyFromText(value[start:end])
}

func canonicalizeURL(source, raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	q := u.Query()
	keep := url.Values{}
	switch source {
	case "linkedin":
		if v := q.Get("currentJobId"); v != "" {
			keep.Set("currentJobId", v)
		}
	case "indeed":
		if v := q.Get("jk"); v != "" {
			keep.Set("jk", v)
		}
	case "glassdoor":
		if v := q.Get("jobListingId"); v != "" {
			keep.Set("jobListingId", v)
		}
	}
	u.RawQuery = keep.Encode()
	u.Fragment = ""
	return u.String()
}

func qualityScore(job ScrapedJob) int {
	score := 0
	if !strings.Contains(strings.ToLower(job.Company), "unknown") && strings.TrimSpace(job.Company) != "" {
		score += 3
	}
	if strings.TrimSpace(job.Compensation) != "" {
		score += 3
	}
	if strings.TrimSpace(job.Title) != "" {
		score += 2
	}
	if strings.TrimSpace(job.Location) != "" {
		score += 1
	}
	if strings.TrimSpace(job.URL) != "" {
		score += 1
	}
	return score
}

func applyQualityGates(jobs []ScrapedJob) []ScrapedJob {
	out := make([]ScrapedJob, 0, len(jobs))
	for _, job := range jobs {
		if qualityScore(job) < 6 {
			continue
		}
		out = append(out, job)
	}
	return out
}

func logMissingFieldRates(jobs []ScrapedJob) {
	if len(jobs) == 0 {
		return
	}
	var missingCompany, missingComp, missingTitle, missingLocation, missingURL int
	for _, j := range jobs {
		if strings.TrimSpace(j.Company) == "" || strings.Contains(strings.ToLower(j.Company), "unknown") {
			missingCompany++
		}
		if strings.TrimSpace(j.Compensation) == "" {
			missingComp++
		}
		if strings.TrimSpace(j.Title) == "" {
			missingTitle++
		}
		if strings.TrimSpace(j.Location) == "" {
			missingLocation++
		}
		if strings.TrimSpace(j.URL) == "" {
			missingURL++
		}
	}
	total := float64(len(jobs))
	log.Printf("scraper quality missing rates: company=%.2f comp=%.2f title=%.2f location=%.2f url=%.2f",
		float64(missingCompany)/total,
		float64(missingComp)/total,
		float64(missingTitle)/total,
		float64(missingLocation)/total,
		float64(missingURL)/total,
	)
}

func dedupeByURL(jobs []ScrapedJob) []ScrapedJob {
	seen := map[string]struct{}{}
	out := make([]ScrapedJob, 0, len(jobs))
	for _, job := range jobs {
		url := strings.TrimSpace(job.URL)
		if url == "" {
			continue
		}
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		out = append(out, job)
	}
	return out
}
