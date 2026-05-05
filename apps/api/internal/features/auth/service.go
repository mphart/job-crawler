package auth

import (
	"errors"
	"log"
	"strings"

	authx "job-crawler/apps/api/internal/platform/auth"
	"job-crawler/apps/api/internal/platform/mail"
)

func normalizeEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ReplaceAll(email, "\uFF20", "@") // FULLWIDTH COMMERCIAL AT
	email = strings.ReplaceAll(email, "\uFE6B", "@") // SMALL COMMERCIAL AT
	return strings.ToLower(email)
}

type AuthStore interface {
	FindUserByEmail(email string) (AuthUser, bool, error)
	CreateUser(email, username, passwordHash string, keywords []string) (AuthUser, error)
}

type AuthUser struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
}

// ProfileApplyFunc persists signup extras (resume, preferences) to the same backing store as profiles.
type ProfileApplyFunc func(userID string, patch map[string]any) error

// NotificationFrequencyApplyFunc persists digest cadence (notification_settings).
type NotificationFrequencyApplyFunc func(userID string, rawFrequency string) error

type Service struct {
	Store                      AuthStore
	Mail                       *mail.Sender
	ApplyProfile               ProfileApplyFunc
	ApplyNotificationFrequency NotificationFrequencyApplyFunc
}

func (s Service) Login(email, password string) (map[string]any, error) {
	email = normalizeEmail(email)
	u, ok, err := s.Store.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("invalid credentials")
	}
	if u.PasswordHash != "" && !authx.VerifyPassword(u.PasswordHash, password) {
		return nil, errors.New("invalid credentials")
	}
	return sessionPayload(u.ID, u.Email, u.Username), nil
}

func (s Service) Signup(req SignupRequest) (map[string]any, error) {
	email := normalizeEmail(req.Email)
	fullName := strings.TrimSpace(req.Name)
	if len(fullName) < 2 {
		return nil, errors.New("please enter your full name (at least two characters).")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least eight characters.")
	}
	freqRaw := strings.TrimSpace(strings.ToLower(req.NotificationFrequency))
	if freqRaw == "" {
		return nil, errors.New("please choose how often you want digest emails.")
	}
	switch freqRaw {
	case "daily", "twice-daily", "weekly":
	default:
		return nil, errors.New("digest frequency must be once a day, twice a day, or once a week.")
	}
	if req.Preferences == nil {
		return nil, errors.New("please provide job preferences (keywords and locations are required).")
	}
	if len(req.Preferences.Keywords) == 0 || len(req.Preferences.Locations) == 0 {
		return nil, errors.New("please provide at least one job keyword and one target location.")
	}

	hash, err := authx.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	u, err := s.Store.CreateUser(email, fullName, hash, nil)
	if err != nil {
		return nil, err
	}

	patch := signupProfilePatch(req)
	if len(patch) > 0 && s.ApplyProfile != nil {
		if err := s.ApplyProfile(u.ID, patch); err != nil {
			return nil, err
		}
	}

	if s.ApplyNotificationFrequency != nil {
		if err := s.ApplyNotificationFrequency(u.ID, freqRaw); err != nil {
			return nil, err
		}
	}

	if s.Mail != nil {
		if !s.Mail.HasHost() {
			log.Printf("signup welcome email skipped: no SMTP host (set WORKER_SMTP_HOST in .env for Docker, or API_SMTP_HOST when running the API outside Compose)")
		} else if err := s.Mail.SendWelcomeRegistration(u.Email, u.Username); err != nil {
			log.Printf("welcome registration email FAILED for %s: %v", u.Email, err)
		}
	}

	return sessionPayload(u.ID, u.Email, u.Username), nil
}

func signupProfilePatch(req SignupRequest) map[string]any {
	patch := map[string]any{}
	if fn := strings.TrimSpace(req.ResumeFileName); fn != "" {
		patch["resumeFileName"] = fn
	}
	if b := strings.TrimSpace(req.ResumeContentBase64); b != "" {
		patch["resumeContentBase64"] = b
	}
	if req.Preferences != nil {
		p := req.Preferences
		patch["preferences"] = map[string]any{
			"keywords":           p.Keywords,
			"locations":          p.Locations,
			"desiredTitles":      p.DesiredTitles,
			"preferredCompanies": p.PreferredCompanies,
			"minComp":            p.MinComp,
			"emailOptIn":         p.EmailOptIn,
			"darkMode":           p.DarkMode,
		}
	}
	return patch
}

func sessionPayload(id, email, username string) map[string]any {
	return map[string]any{
		"id":       id,
		"email":    email,
		"username": username,
		"token":    authx.IssueToken(id),
	}
}
