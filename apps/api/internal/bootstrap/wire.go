package bootstrap

import (
	"log"
	"net/http"

	"job-crawler/apps/api/internal/features/auth"
	"job-crawler/apps/api/internal/features/feed"
	"job-crawler/apps/api/internal/features/notifications"
	"job-crawler/apps/api/internal/features/profiles"
	"job-crawler/apps/api/internal/features/users"
	"job-crawler/apps/api/internal/features/worker"
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
	authStore := auth.AuthStore(auth.InMemoryStore{Inner: store})
	profileStore := profiles.ProfileStore(profiles.InMemoryStore{Inner: store})
	userSearchStore := users.UserSearchStore(users.InMemoryStore{Inner: store})
	feedStore := feed.FeedStore(feed.InMemoryStore{Inner: store})
	notificationStore := notifications.NotificationStore(notifications.InMemoryStore{Inner: store})
	mysqlAuthStore, err := db.NewMySQLAuthStore(cfg.MySQLDSN)
	if err != nil {
		log.Printf("mysql auth store unavailable, using in-memory auth store: %v", err)
	} else {
		authStore = auth.MySQLStore{Inner: mysqlAuthStore}
		profileStore = profiles.MySQLStore{Inner: mysqlAuthStore}
		userSearchStore = users.MySQLStore{Inner: mysqlAuthStore}
		feedStore = feed.MySQLStore{Inner: mysqlAuthStore}
		notificationStore = notifications.MySQLStore{Inner: mysqlAuthStore}
	}

	handlers := httpx.Handlers{
		AuthLogin:              auth.Handler{Service: auth.Service{Store: authStore}}.Login,
		AuthSignup:             auth.Handler{Service: auth.Service{Store: authStore}}.Signup,
		FeedList:               feed.Handler{Service: feed.Service{Store: feedStore}}.List,
		FeedAction:             feed.Handler{Service: feed.Service{Store: feedStore}}.Action,
		WorkerTick:             worker.Handler{Service: worker.Service{Store: store, MySQL: mysqlAuthStore}, WorkerToken: cfg.WorkerToken}.Tick,
		WorkerPreferences:      worker.Handler{Service: worker.Service{Store: store, MySQL: mysqlAuthStore}, WorkerToken: cfg.WorkerToken}.Preferences,
		WorkerDigestCandidates: worker.Handler{Service: worker.Service{Store: store, MySQL: mysqlAuthStore}, WorkerToken: cfg.WorkerToken}.DigestCandidates,
		WorkerNotificationSent: worker.Handler{Service: worker.Service{Store: store, MySQL: mysqlAuthStore}, WorkerToken: cfg.WorkerToken}.NotificationSent,
		ProfileUpdateMe:        profiles.Handler{Service: profiles.Service{Store: profileStore}}.UpdateMe,
		ProfileGetByID:         profiles.Handler{Service: profiles.Service{Store: profileStore}}.GetByID,
		NotificationSettings:   notifications.Handler{Service: notifications.Service{Store: notificationStore}}.Update,
		UsersSearch:            users.Handler{Service: users.Service{Store: userSearchStore}}.Search,
	}

	return App{Config: cfg, Router: httpx.NewRouter(handlers)}
}
