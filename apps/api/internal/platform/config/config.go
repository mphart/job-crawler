package config

import "os"

type Config struct {
	HTTPAddr    string
	JWTSecret   string
	WorkerToken string
	MySQLDSN    string
}

func Load() Config {
	return Config{
		HTTPAddr:    get("API_HTTP_ADDR", ":8080"),
		JWTSecret:   get("API_JWT_SECRET", "dev-secret"),
		WorkerToken: get("API_WORKER_TOKEN", "worker-dev-token"),
		MySQLDSN:    get("API_MYSQL_DSN", "jobcrawler:jobcrawler@tcp(mysql:3306)/jobcrawler?parseTime=true"),
	}
}

func get(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
