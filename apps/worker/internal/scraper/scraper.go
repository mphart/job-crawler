package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
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

	// Hardcoded source priority: LinkedIn and Indeed first, then Greenhouse.
	for _, keyword := range keywords {
		linkedin, err := scrapeLinkedIn(ctx, client, keyword)
		if err == nil {
			out = append(out, linkedin...)
		}
		indeed, err := scrapeIndeed(ctx, client, keyword)
		if err == nil {
			out = append(out, indeed...)
		}
		greenhouse, err := scrapeGreenhouse(ctx, client, keyword)
		if err == nil {
			out = append(out, greenhouse...)
		}
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no jobs scraped from configured sources (linkedin/indeed/greenhouse)")
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
	titleRegex := regexp.MustCompile(`(?s)base-search-card__title[^>]*>\s*(.*?)\s*<`)
	companyRegex := regexp.MustCompile(`(?s)base-search-card__subtitle[^>]*>\s*(.*?)\s*<`)
	companyAltRegex := regexp.MustCompile(`(?s)job-search-card__subtitle-link[^>]*>\s*(.*?)\s*<`)
	locationRegex := regexp.MustCompile(`(?s)job-search-card__location[^>]*>\s*(.*?)\s*<`)
	compensationRegex := regexp.MustCompile(`(?s)(job-search-card__salary-info|base-search-card__metadata)[^>]*>\s*(.*?)\s*<`)
	linkRegex := regexp.MustCompile(`href="([^"]*linkedin\.com/jobs/view/[^"]+)"`)
	idRegex := regexp.MustCompile(`data-entity-urn="urn:li:jobPosting:(\d+)"`)

	cards := cardRegex.FindAllString(body, -1)
	jobs := make([]ScrapedJob, 0, len(cards))
	for _, card := range cards {
		title := capture(titleRegex, card)
		company := capture(companyRegex, card)
		if strings.TrimSpace(company) == "" {
			company = capture(companyAltRegex, card)
		}
		location := capture(locationRegex, card)
		compensation := capture(compensationRegex, card)
		if strings.TrimSpace(compensation) == "" {
			compensation = captureMoneyFromText(card)
		}
		link := capture(linkRegex, card)
		if strings.TrimSpace(title) == "" || strings.TrimSpace(link) == "" {
			continue
		}
		if !matchesKeyword(title, keyword) {
			continue
		}
		externalID := capture(idRegex, card)
		jobs = append(jobs, ScrapedJob{
			Source:       "linkedin",
			ExternalID:   externalID,
			Company:      normalizeCompany(html.UnescapeString(clean(stripTags(company))), "linkedin"),
			Title:        html.UnescapeString(clean(title)),
			Location:     html.UnescapeString(clean(location)),
			Compensation: normalizeCompensation(html.UnescapeString(clean(stripTags(compensation)))),
			PostedAt:     time.Now().Format(time.RFC3339),
			URL:          html.UnescapeString(clean(link)),
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
	titleRegex := regexp.MustCompile(`(?s)jobTitle[^>]*>\s*(.*?)\s*<`)
	companyRegex := regexp.MustCompile(`(?s)data-testid="company-name"[^>]*>\s*(.*?)\s*<`)
	companyAltRegex := regexp.MustCompile(`(?s)class="[^"]*companyName[^"]*"[^>]*>\s*(.*?)\s*<`)
	locationRegex := regexp.MustCompile(`(?s)data-testid="text-location"[^>]*>\s*(.*?)\s*<`)
	salaryRegex := regexp.MustCompile(`(?s)class="[^"]*salary-snippet[^"]*"[^>]*>\s*(.*?)\s*<`)
	linkRegex := regexp.MustCompile(`href="(/pagead/clk[^"]+|/rc/clk[^"]+|/viewjob\?jk=[^"]+)"`)
	idRegex := regexp.MustCompile(`jk=([a-zA-Z0-9]+)`)

	cards := cardRegex.FindAllString(body, -1)
	jobs := make([]ScrapedJob, 0, len(cards))
	for _, card := range cards {
		title := capture(titleRegex, card)
		link := capture(linkRegex, card)
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
		externalID := capture(idRegex, fullLink)
		company := capture(companyRegex, card)
		if strings.TrimSpace(company) == "" {
			company = capture(companyAltRegex, card)
		}
		jobs = append(jobs, ScrapedJob{
			Source:       "indeed",
			ExternalID:   externalID,
			Company:      normalizeCompany(html.UnescapeString(clean(stripTags(company))), "indeed"),
			Title:        html.UnescapeString(clean(title)),
			Location:     html.UnescapeString(clean(capture(locationRegex, card))),
			Compensation: html.UnescapeString(clean(capture(salaryRegex, card))),
			PostedAt:     time.Now().Format(time.RFC3339),
			URL:          html.UnescapeString(clean(fullLink)),
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
			URL:          j.AbsoluteURL,
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
	default:
		return "Unknown Company"
	}
}

func normalizeCompensation(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	return trimmed
}

func captureMoneyFromText(value string) string {
	moneyRegex := regexp.MustCompile(`(?i)\$[\d,]+(?:\s*-\s*\$[\d,]+)?(?:\s*(?:/year|yr|hour|hr|k))?`)
	match := moneyRegex.FindString(value)
	return strings.TrimSpace(match)
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
