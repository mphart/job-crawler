package runner

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"job-crawler/apps/worker/internal/platform/config"
)

func testConfig(baseURL string) config.Config {
	return config.Config{
		APIBaseURL:     baseURL,
		RunInterval:    30 * time.Second,
		RequestTimeout: 2 * time.Second,
		APIToken:       "test-token",
		MaxSyncRetries: 3,
		RetryBackoff:   10 * time.Millisecond,
	}
}

func TestTick_RetriesSyncUntilSuccess(t *testing.T) {
	var syncCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.WriteHeader(http.StatusOK)
		case "/api/worker/tick":
			if r.Header.Get("X-Worker-Token") != "test-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			call := atomic.AddInt32(&syncCalls, 1)
			if call < 3 {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("transient"))
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	r := New(testConfig(server.URL))

	if err := r.tick(context.Background()); err != nil {
		t.Fatalf("expected retry to eventually succeed, got error: %v", err)
	}

	if got := atomic.LoadInt32(&syncCalls); got != 3 {
		t.Fatalf("expected 3 sync attempts, got %d", got)
	}
}

func TestTick_FailsAfterExhaustedRetries(t *testing.T) {
	var syncCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/healthz":
			w.WriteHeader(http.StatusOK)
		case "/api/worker/tick":
			atomic.AddInt32(&syncCalls, 1)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("always failing"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	cfg := testConfig(server.URL)
	cfg.MaxSyncRetries = 2
	r := New(cfg)

	if err := r.tick(context.Background()); err == nil {
		t.Fatal("expected tick to fail after retries, got nil")
	}

	if got := atomic.LoadInt32(&syncCalls); got != 2 {
		t.Fatalf("expected 2 sync attempts, got %d", got)
	}
}

func TestTick_FailsWhenHealthCheckFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("api down"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	r := New(testConfig(server.URL))
	if err := r.tick(context.Background()); err == nil {
		t.Fatal("expected health failure error, got nil")
	}
}
