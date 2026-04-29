package users

import "job-crawler/apps/api/internal/platform/db"

type UserSearchStore interface {
	SearchUsers(q string) ([]map[string]any, error)
}

type InMemoryStore struct{ Inner *db.Store }

func (s InMemoryStore) SearchUsers(q string) ([]map[string]any, error) {
	return s.Inner.SearchUsers(q), nil
}

type MySQLStore struct{ Inner *db.MySQLAuthStore }

func (s MySQLStore) SearchUsers(q string) ([]map[string]any, error) {
	return s.Inner.SearchUsers(q)
}

type Service struct{ Store UserSearchStore }

func (s Service) Search(q string) ([]map[string]any, error) { return s.Store.SearchUsers(q) }
