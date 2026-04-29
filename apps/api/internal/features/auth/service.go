package auth

import (
	"errors"

	authx "job-crawler/apps/api/internal/platform/auth"
)

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

type Service struct{ Store AuthStore }

func (s Service) Login(email, password string) (map[string]any, error) {
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

func (s Service) Signup(email, username, password string, keywords []string) (map[string]any, error) {
	hash, err := authx.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u, err := s.Store.CreateUser(email, username, hash, keywords)
	if err != nil {
		return nil, err
	}
	return sessionPayload(u.ID, u.Email, u.Username), nil
}

func sessionPayload(id, email, username string) map[string]any {
	return map[string]any{
		"id":       id,
		"email":    email,
		"username": username,
		"token":    authx.IssueToken(id),
	}
}
