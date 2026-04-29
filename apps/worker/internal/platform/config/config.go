package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIBaseURL     string
	RunInterval    time.Duration
	RequestTimeout time.Duration
	BearerToken    string
	APIToken       string
	MaxSyncRetries int
	RetryBackoff   time.Duration
	SMTPHost       string
	SMTPPort       int
	SMTPUsername   string
	SMTPPassword   string
	SMTPFrom       string
}

func Load() Config {
	return Config{
		APIBaseURL:     get("WORKER_API_BASE_URL", "http://api:8080"),
		RunInterval:    duration("WORKER_RUN_INTERVAL", 30*time.Second),
		RequestTimeout: duration("WORKER_REQUEST_TIMEOUT", 5*time.Second),
		BearerToken:    os.Getenv("WORKER_BEARER_TOKEN"),
		APIToken:       get("WORKER_API_TOKEN", "worker-dev-token"),
		MaxSyncRetries: intValue("WORKER_MAX_SYNC_RETRIES", 3),
		RetryBackoff:   duration("WORKER_RETRY_BACKOFF", 1500*time.Millisecond),
		SMTPHost:       os.Getenv("WORKER_SMTP_HOST"),
		SMTPPort:       intValue("WORKER_SMTP_PORT", 587),
		SMTPUsername:   os.Getenv("WORKER_SMTP_USERNAME"),
		SMTPPassword:   os.Getenv("WORKER_SMTP_PASSWORD"),
		SMTPFrom:       get("WORKER_SMTP_FROM", "no-reply@jobcrawler.local"),
	}
}

func get(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func duration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func intValue(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	if parsed < 0 {
		return fallback
	}
	return parsed
}
