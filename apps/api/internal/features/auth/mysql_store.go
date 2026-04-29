package auth

import "job-crawler/apps/api/internal/platform/db"

type MySQLStore struct {
	Inner *db.MySQLAuthStore
}

func (s MySQLStore) FindUserByEmail(email string) (AuthUser, bool, error) {
	user, ok, err := s.Inner.FindUserByEmail(email)
	if err != nil {
		return AuthUser{}, false, err
	}
	return AuthUser{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}, ok, nil
}

func (s MySQLStore) CreateUser(email, username, passwordHash string, keywords []string) (AuthUser, error) {
	user, err := s.Inner.CreateUser(email, username, passwordHash, keywords)
	if err != nil {
		return AuthUser{}, err
	}
	return AuthUser{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}, nil
}
