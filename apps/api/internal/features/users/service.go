package users

import "job-crawler/apps/api/internal/platform/db"

type Service struct{ Store *db.Store }

func (s Service) Search(q string) []map[string]any { return s.Store.SearchUsers(q) }
