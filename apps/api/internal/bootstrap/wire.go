package bootstrap

import (
	"errors"
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
	"job-crawler/apps/api/internal/platform/mail"
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
	var resumeDispatcher *profiles.ResumeParseDispatcher
	if err != nil {
		log.Printf("mysql auth store unavailable, using in-memory auth store: %v", err)
	} else {
		authStore = auth.MySQLStore{Inner: mysqlAuthStore}
		profileStore = profiles.MySQLStore{Inner: mysqlAuthStore}
		userSearchStore = users.MySQLStore{Inner: mysqlAuthStore}
		feedStore = feed.MySQLStore{Inner: mysqlAuthStore}
		notificationStore = notifications.MySQLStore{Inner: mysqlAuthStore}
		resumeDispatcher = profiles.NewResumeParseDispatcher(mysqlAuthStore)
	}

	mailSender := mail.NewSender(mail.Config{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUser,
		Password: cfg.SMTPPass,
		From:     cfg.SMTPFrom,
	})
	if mailSender.HasHost() {
		if err := mailSender.ValidateSMTP(); err != nil {
			log.Printf("api: SMTP incomplete — welcome emails will fail until fixed: %v", err)
		} else {
			log.Printf("api: transactional SMTP configured (host=%q from=%q)", cfg.SMTPHost, cfg.SMTPFrom)
		}
	} else {
		log.Printf("api: transactional SMTP disabled — set WORKER_SMTP_HOST (and username, password, from) in .env; Compose maps them to the API container")
	}

	var applySignupProfile auth.ProfileApplyFunc
	if mysqlAuthStore != nil {
		applySignupProfile = func(userID string, patch map[string]any) error {
			_, ok, err := mysqlAuthStore.UpdateMe(userID, patch)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("could not save profile after signup")
			}
			return nil
		}
	} else {
		applySignupProfile = func(userID string, patch map[string]any) error {
			_, ok := store.UpdateMe(userID, patch)
			if !ok {
				return errors.New("could not save profile after signup")
			}
			return nil
		}
	}

	var applyNotificationFrequency auth.NotificationFrequencyApplyFunc
	if mysqlAuthStore != nil {
		applyNotificationFrequency = func(userID string, raw string) error {
			return mysqlAuthStore.SetNotificationFrequency(userID, raw)
		}
	} else {
		applyNotificationFrequency = func(userID string, raw string) error {
			store.NotificationFrequency[userID] = db.NormalizeNotificationFrequency(raw)
			return nil
		}
	}

	authSvc := auth.Service{
		Store:                      authStore,
		Mail:                       mailSender,
		ApplyProfile:               applySignupProfile,
		ApplyNotificationFrequency: applyNotificationFrequency,
	}

	handlers := httpx.Handlers{
		AuthLogin:              auth.Handler{Service: authSvc}.Login,
		AuthSignup:             auth.Handler{Service: authSvc}.Signup,
		FeedList:               feed.Handler{Service: feed.Service{Store: feedStore}}.List,
		FeedAction:             feed.Handler{Service: feed.Service{Store: feedStore}}.Action,
		WorkerTick:             worker.Handler{Service: worker.Service{Store: store, MySQL: mysqlAuthStore}, WorkerToken: cfg.WorkerToken}.Tick,
		WorkerPreferences:      worker.Handler{Service: worker.Service{Store: store, MySQL: mysqlAuthStore}, WorkerToken: cfg.WorkerToken}.Preferences,
		WorkerDigestCandidates: worker.Handler{Service: worker.Service{Store: store, MySQL: mysqlAuthStore}, WorkerToken: cfg.WorkerToken}.DigestCandidates,
		WorkerNotificationSent: worker.Handler{Service: worker.Service{Store: store, MySQL: mysqlAuthStore}, WorkerToken: cfg.WorkerToken}.NotificationSent,
		ProfileUpdateMe:        profiles.Handler{Service: profiles.Service{Store: profileStore, Dispatcher: resumeDispatcher}}.UpdateMe,
		ProfileGetByID:         profiles.Handler{Service: profiles.Service{Store: profileStore, Dispatcher: resumeDispatcher}}.GetByID,
		NotificationSettings:   notifications.Handler{Service: notifications.Service{Store: notificationStore}}.Update,
		UsersSearch:            users.Handler{Service: users.Service{Store: userSearchStore}}.Search,
	}

	return App{Config: cfg, Router: httpx.NewRouter(handlers)}
}
