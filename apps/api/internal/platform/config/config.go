package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr    string
	JWTSecret   string
	WorkerToken string
	MySQLDSN    string
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
}

func Load() Config {
	return Config{
		HTTPAddr:    get("API_HTTP_ADDR", ":8080"),
		JWTSecret:   get("API_JWT_SECRET", "dev-secret"),
		WorkerToken: get("API_WORKER_TOKEN", "worker-dev-token"),
		MySQLDSN:    get("API_MYSQL_DSN", "jobcrawler:jobcrawler@tcp(mysql:3306)/jobcrawler?parseTime=true"),
		SMTPHost:    get("API_SMTP_HOST", ""),
		SMTPPort:    intEnv("API_SMTP_PORT", 587),
		SMTPUser:    get("API_SMTP_USERNAME", ""),
		SMTPPass:    get("API_SMTP_PASSWORD", ""),
		SMTPFrom:    get("API_SMTP_FROM", "no-reply@jobcrawler.local"),
	}
}

func get(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

func intEnv(k string, fallback int) int {
	v := os.Getenv(k)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
