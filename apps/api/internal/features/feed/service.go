package feed

import "job-crawler/apps/api/internal/platform/db"

type Service struct{ Store *db.Store }

func (s Service) List(userID, search, sortBy string) []db.JobPosting { return s.Store.Feed(userID, search, sortBy) }
func (s Service) Apply(userID, jobID string) error { return s.Store.Decide(userID, jobID, "APPLIED") }
func (s Service) Reject(userID, jobID string) error { return s.Store.Decide(userID, jobID, "REJECTED") }
