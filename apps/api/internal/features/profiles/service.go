package profiles

import "job-crawler/apps/api/internal/platform/db"

type Service struct{ Store *db.Store }

func (s Service) Get(requester, userID string) (db.Profile, bool) { return s.Store.Profile(requester, userID) }
func (s Service) UpdateMe(userID string, patch map[string]any) (db.Profile, bool) { return s.Store.UpdateMe(userID, patch) }
