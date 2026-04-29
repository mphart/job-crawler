package main

import (
    "context"
    "log"
    "os/signal"
    "syscall"

    "job-crawler/apps/worker/internal/platform/config"
    "job-crawler/apps/worker/internal/runner"
)

func main() {
    cfg := config.Load()
    worker := runner.New(cfg)

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    log.Printf("worker starting; interval=%s api=%s", cfg.RunInterval, cfg.APIBaseURL)
    if err := worker.Run(ctx); err != nil {
        log.Fatal(err)
    }
    log.Println("worker stopped")
}
