package auth

import (
	"errors"
	"testing"

	"job-crawler/apps/api/internal/platform/db"
)

func TestSignupAndLoginIssueToken(t *testing.T) {
	store := db.NewStore()
	svc := Service{
		Store: InMemoryStore{Inner: store},
		ApplyProfile: func(userID string, patch map[string]any) error {
			_, ok := store.UpdateMe(userID, patch)
			if !ok {
				return errors.New("profile apply failed")
			}
			return nil
		},
		ApplyNotificationFrequency: func(userID string, raw string) error {
			store.NotificationFrequency[userID] = raw
			return nil
		},
	}

	created, err := svc.Signup(SignupRequest{
		Email:                   "new@example.com",
		Name:                    "New User",
		Password:                "password123",
		NotificationFrequency:   "daily",
		Preferences: &SignupPreferences{
			Keywords:      []string{"frontend"},
			Locations:     []string{"Remote"},
			DesiredTitles: []string{"Software Engineer"},
			MinComp:       120000,
			EmailOptIn:    true,
		},
	})
	if err != nil {
		t.Fatalf("signup failed: %v", err)
	}
	if created["id"] == "" || created["token"] == "" {
		t.Fatalf("expected id and token in signup response")
	}

	session, err := svc.Login("new@example.com", "password123")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if session["email"] != "new@example.com" {
		t.Fatalf("expected matching login email")
	}
	if session["token"] == "" {
		t.Fatalf("expected token in login response")
	}
}
