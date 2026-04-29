package config

import "testing"

func TestLoadDefaults(t *testing.T) {
    cfg := Load()
    if cfg.HTTPAddr == "" || cfg.JWTSecret == "" {
        t.Fatal("expected defaults")
    }
}
