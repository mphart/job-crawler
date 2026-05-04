package profiles

import "job-crawler/apps/api/internal/platform/db"

type ProfileStore interface {
	Profile(requester, userID string) (db.Profile, bool, error)
	UpdateMe(userID string, patch map[string]any) (db.Profile, bool, error)
}

type InMemoryStore struct {
	Inner *db.Store
}

func (s InMemoryStore) Profile(requester, userID string) (db.Profile, bool, error) {
	p, ok := s.Inner.Profile(requester, userID)
	return p, ok, nil
}

func (s InMemoryStore) UpdateMe(userID string, patch map[string]any) (db.Profile, bool, error) {
	p, ok := s.Inner.UpdateMe(userID, patch)
	return p, ok, nil
}

type MySQLStore struct {
	Inner *db.MySQLAuthStore
}

func (s MySQLStore) Profile(requester, userID string) (db.Profile, bool, error) {
	return s.Inner.Profile(requester, userID)
}

func (s MySQLStore) UpdateMe(userID string, patch map[string]any) (db.Profile, bool, error) {
	return s.Inner.UpdateMe(userID, patch)
}

type ParseDispatcher interface {
	Enqueue(userID, resumeContentBase64 string) error
}

type Service struct {
	Store      ProfileStore
	Dispatcher ParseDispatcher
}

func (s Service) Get(requester, userID string) (db.Profile, bool, error) {
	return s.Store.Profile(requester, userID)
}
func (s Service) UpdateMe(userID string, patch map[string]any) (db.Profile, bool, error) {
	p, ok, err := s.Store.UpdateMe(userID, patch)
	if err != nil || !ok {
		return p, ok, err
	}
	if s.Dispatcher != nil {
		if raw, exists := patch["resumeContentBase64"].(string); exists && raw != "" {
			_ = s.Dispatcher.Enqueue(userID, raw)
		}
	}
	return p, ok, nil
}
