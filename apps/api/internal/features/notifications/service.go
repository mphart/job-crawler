package notifications

import "job-crawler/apps/api/internal/platform/db"

type NotificationStore interface {
	UpdateNotification(userID string, payload map[string]any) (map[string]any, error)
	GetNotification(userID string) (map[string]any, error)
}

type InMemoryStore struct{ Inner *db.Store }

func (s InMemoryStore) UpdateNotification(userID string, payload map[string]any) (map[string]any, error) {
	if v, ok := payload["frequency"].(string); ok {
		s.Inner.NotificationFrequency[userID] = db.NormalizeNotificationFrequency(v)
	}
	f := s.Inner.NotificationFrequency[userID]
	if f == "" {
		f = "daily"
	}
	emailOptIn := true
	if v, ok := payload["emailOptIn"].(bool); ok {
		emailOptIn = v
	}
	return map[string]any{"emailOptIn": emailOptIn, "frequency": f}, nil
}

func (s InMemoryStore) GetNotification(userID string) (map[string]any, error) {
	frequency := db.NormalizeNotificationFrequency(s.Inner.NotificationFrequency[userID])
	if frequency == "" {
		frequency = "daily"
	}
	return map[string]any{"emailOptIn": true, "frequency": frequency}, nil
}

type MySQLStore struct{ Inner *db.MySQLAuthStore }

func (s MySQLStore) UpdateNotification(userID string, payload map[string]any) (map[string]any, error) {
	return s.Inner.UpdateNotification(userID, payload)
}

func (s MySQLStore) GetNotification(userID string) (map[string]any, error) {
	return s.Inner.GetNotification(userID)
}

type Service struct{ Store NotificationStore }

func (s Service) Update(userID string, payload map[string]any) (map[string]any, error) {
	return s.Store.UpdateNotification(userID, payload)
}

func (s Service) Get(userID string) (map[string]any, error) {
	return s.Store.GetNotification(userID)
}
