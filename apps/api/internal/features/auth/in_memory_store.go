package auth

import "job-crawler/apps/api/internal/platform/db"

type InMemoryStore struct {
	Inner *db.Store
}

func (s InMemoryStore) FindUserByEmail(email string) (AuthUser, bool, error) {
	user, ok := s.Inner.FindUserByEmail(email)
	if !ok {
		return AuthUser{}, false, nil
	}
	return AuthUser{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}, true, nil
}

func (s InMemoryStore) CreateUser(email, username, passwordHash string, keywords []string) (AuthUser, error) {
	user := s.Inner.CreateUser(email, username, passwordHash, keywords)
	return AuthUser{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}, nil
}
