package runner

import (
    "context"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
    "time"

    "job-crawler/apps/worker/internal/platform/config"
)

type Runner struct {
    cfg config.Config
    client *http.Client
}

func New(cfg config.Config) Runner {
    return Runner{
        cfg: cfg,
        client: &http.Client{Timeout: cfg.RequestTimeout},
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

    response, err := r.client.Do(req)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    if response.StatusCode < 200 || response.StatusCode > 299 {
        payload, _ := io.ReadAll(response.Body)
        return fmt.Errorf("health request failed with %d: %s", response.StatusCode, string(payload))
    }

    log.Printf("worker tick ok (%s)", healthURL)
    return r.syncJobs(ctx)
}

func (r Runner) syncJobs(ctx context.Context) error {
    syncURL := strings.TrimRight(r.cfg.APIBaseURL, "/") + "/api/worker/tick"
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, syncURL, nil)
    if err != nil {
        return err
    }
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
