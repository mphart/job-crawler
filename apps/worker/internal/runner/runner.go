package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"job-crawler/apps/worker/internal/notify"
	"job-crawler/apps/worker/internal/platform/config"
	"job-crawler/apps/worker/internal/scraper"
)

type Runner struct {
	cfg      config.Config
	client   *http.Client
	scrapeFn func(context.Context, *http.Client, []string) ([]scraper.ScrapedJob, error)
}

func New(cfg config.Config) Runner {
	return Runner{
		cfg:      cfg,
		client:   &http.Client{Timeout: cfg.RequestTimeout},
		scrapeFn: scraper.Scrape,
	}
}

func (r Runner) Run(ctx context.Context) error {
	if err := r.tick(ctx); err != nil {
		log.Printf("initial worker tick failed: %v", err)
	}

	ticker := time.NewTicker(r.cfg.RunInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := r.tick(ctx); err != nil {
				log.Printf("worker tick failed: %v", err)
			}
		}
	}
}

func (r Runner) tick(ctx context.Context) error {
	healthURL := strings.TrimRight(r.cfg.APIBaseURL, "/") + "/healthz"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
	if err != nil {
		return err
	}

	if r.cfg.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+r.cfg.BearerToken)
	}

	if err := r.doRequest(req, "health request"); err != nil {
		return err
	}

	log.Printf("worker tick ok (%s)", healthURL)
	keywords, err := r.fetchKeywords(ctx)
	if err != nil {
		return fmt.Errorf("preferences fetch failed: %w", err)
	}
	jobs, err := r.scrapeFn(ctx, r.client, keywords)
	if err != nil {
		return fmt.Errorf("scrape failed: %w", err)
	}
	if len(jobs) == 0 {
		return errors.New("scrape returned zero jobs")
	}
	if err := r.syncJobsWithRetry(ctx, jobs); err != nil {
		return err
	}
	return r.sendDigests(ctx)
}

func (r Runner) fetchKeywords(ctx context.Context) ([]string, error) {
	endpoint := strings.TrimRight(r.cfg.APIBaseURL, "/") + "/api/worker/preferences"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if r.cfg.APIToken != "" {
		req.Header.Set("X-Worker-Token", r.cfg.APIToken)
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("preferences request failed with %d: %s", res.StatusCode, string(body))
	}
	var payload struct {
		Keywords []string `json:"keywords"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if len(payload.Keywords) == 0 {
		return []string{"software engineer"}, nil
	}
	return payload.Keywords, nil
}

func (r Runner) syncJobs(ctx context.Context, jobs []scraper.ScrapedJob) error {
	syncURL := strings.TrimRight(r.cfg.APIBaseURL, "/") + "/api/worker/tick"
	payload, err := json.Marshal(jobs)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, syncURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.cfg.APIToken != "" {
		req.Header.Set("X-Worker-Token", r.cfg.APIToken)
	}

	response, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		payload, _ := io.ReadAll(response.Body)
		return fmt.Errorf("worker sync failed with %d: %s", response.StatusCode, string(payload))
	}

	log.Printf("worker sync ok (%s)", syncURL)
	return nil
}

func (r Runner) syncJobsWithRetry(ctx context.Context, jobs []scraper.ScrapedJob) error {
	attempts := r.cfg.MaxSyncRetries
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := r.syncJobs(ctx, jobs); err == nil {
			return nil
		} else {
			lastErr = err
			if attempt == attempts {
				break
			}
			log.Printf("worker sync attempt %d/%d failed: %v", attempt, attempts, err)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(r.cfg.RetryBackoff):
			}
		}
	}

	return fmt.Errorf("worker sync exhausted retries: %w", lastErr)
}

func (r Runner) doRequest(req *http.Request, op string) error {
	response, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		payload, _ := io.ReadAll(response.Body)
		if len(payload) == 0 {
			return errors.New(op + " failed with status " + response.Status)
		}
		return fmt.Errorf("%s failed with %d: %s", op, response.StatusCode, string(payload))
	}

	return nil
}

func (r Runner) sendDigests(ctx context.Context) error {
	candidates, err := r.fetchDigestCandidates(ctx)
	if err != nil {
		return err
	}
	for _, candidate := range candidates {
		if err := notify.SendDigest(r.cfg, candidate); err != nil {
			log.Printf("digest send failed for %s: %v", candidate.Email, err)
			continue
		}
		if err := r.markNotificationSent(ctx, candidate.UserID); err != nil {
			log.Printf("failed to mark notification sent for %s: %v", candidate.UserID, err)
		}
	}
	return nil
}

func (r Runner) fetchDigestCandidates(ctx context.Context) ([]notify.Candidate, error) {
	endpoint := strings.TrimRight(r.cfg.APIBaseURL, "/") + "/api/worker/digest-candidates"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if r.cfg.APIToken != "" {
		req.Header.Set("X-Worker-Token", r.cfg.APIToken)
	}
	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("digest candidates request failed with %d: %s", res.StatusCode, string(body))
	}
	var payload struct {
		Candidates []notify.Candidate `json:"candidates"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload.Candidates, nil
}

func (r Runner) markNotificationSent(ctx context.Context, userID string) error {
	endpoint := strings.TrimRight(r.cfg.APIBaseURL, "/") + "/api/worker/notifications/sent"
	body, _ := json.Marshal(map[string]string{"userId": userID})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if r.cfg.APIToken != "" {
		req.Header.Set("X-Worker-Token", r.cfg.APIToken)
	}
	return r.doRequest(req, "mark notification sent")
}
