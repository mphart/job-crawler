package auth

import (
    "errors"

    authx "job-crawler/apps/api/internal/platform/auth"
    "job-crawler/apps/api/internal/platform/db"
)

type Service struct{ Store *db.Store }

func (s Service) Login(email, password string) (map[string]any, error) {
    u, ok := s.Store.FindUserByEmail(email)
    if !ok { return nil, errors.New("invalid credentials") }
    if u.PasswordHash != "" && !authx.VerifyPassword(u.PasswordHash, password) { return nil, errors.New("invalid credentials") }
    return map[string]any{"id":u.ID,"email":u.Email,"username":u.Username}, nil
}

func (s Service) Signup(email, username, password string, keywords []string) (map[string]any, error) {
    hash, err := authx.HashPassword(password)
    if err != nil { return nil, err }
    u := s.Store.CreateUser(email, username, hash, keywords)
    return map[string]any{"id":u.ID,"email":u.Email,"username":u.Username}, nil
}
