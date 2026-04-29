package notifications

import "job-crawler/apps/api/internal/platform/db"

type Service struct{ Store *db.Store }

func (s Service) Update(userID string, payload map[string]any) map[string]any { if v,ok:=payload["frequency"].(string);ok{s.Store.NotificationFrequency[userID]=v}; return map[string]any{"emailOptIn": payload["emailOptIn"], "frequency": s.Store.NotificationFrequency[userID]} }
