package config

import "os"

type Config struct {
    HTTPAddr string
    JWTSecret string
}

func Load() Config {
    return Config{
        HTTPAddr: get("API_HTTP_ADDR", ":8080"),
        JWTSecret: get("API_JWT_SECRET", "dev-secret"),
    }
}

func get(k, fallback string) string {
    if v := os.Getenv(k); v != "" { return v }
    return fallback
}
