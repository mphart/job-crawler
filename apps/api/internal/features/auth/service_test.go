package auth

import (
	"testing"

	"job-crawler/apps/api/internal/platform/db"
)

func TestSignupAndLoginIssueToken(t *testing.T) {
	svc := Service{Store: InMemoryStore{Inner: db.NewStore()}}

	created, err := svc.Signup("new@example.com", "new-user", "password123", []string{"frontend"})
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
