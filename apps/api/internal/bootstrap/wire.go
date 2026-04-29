package bootstrap

import (
    "net/http"

    "job-crawler/apps/api/internal/features/auth"
    "job-crawler/apps/api/internal/features/feed"
    "job-crawler/apps/api/internal/features/notifications"
    "job-crawler/apps/api/internal/features/profiles"
    "job-crawler/apps/api/internal/features/users"
    "job-crawler/apps/api/internal/platform/config"
    "job-crawler/apps/api/internal/platform/db"
    httpx "job-crawler/apps/api/internal/platform/http"
)

type App struct {
    Config config.Config
    Router http.Handler
}

func NewApp() App {
    cfg := config.Load()
    store := db.NewStore()

    handlers := httpx.Handlers{
        Auth: &auth.Handler{Service: auth.Service{Store: store}},
        Feed: &feed.Handler{Service: feed.Service{Store: store}},
        Profiles: &profiles.Handler{Service: profiles.Service{Store: store}},
        Notifications: &notifications.Handler{Service: notifications.Service{Store: store}},
        Users: &users.Handler{Service: users.Service{Store: store}},
    }

    return App{Config: cfg, Router: httpx.NewRouter(handlers)}
}
