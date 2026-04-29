package feed

import "job-crawler/apps/api/internal/platform/db"

type FeedStore interface {
	Feed(userID, search, sortBy string) ([]db.JobPosting, error)
	Decide(userID, jobID, decision string) error
}

type InMemoryStore struct{ Inner *db.Store }

func (s InMemoryStore) Feed(userID, search, sortBy string) ([]db.JobPosting, error) {
	return s.Inner.Feed(userID, search, sortBy), nil
}
func (s InMemoryStore) Decide(userID, jobID, decision string) error {
	return s.Inner.Decide(userID, jobID, decision)
}

type MySQLStore struct{ Inner *db.MySQLAuthStore }

func (s MySQLStore) Feed(userID, search, sortBy string) ([]db.JobPosting, error) {
	return s.Inner.Feed(userID, search, sortBy)
}
func (s MySQLStore) Decide(userID, jobID, decision string) error {
	return s.Inner.Decide(userID, jobID, decision)
}

type Service struct{ Store FeedStore }

func (s Service) List(userID, search, sortBy string) ([]db.JobPosting, error) {
	return s.Store.Feed(userID, search, sortBy)
}
func (s Service) Apply(userID, jobID string) error  { return s.Store.Decide(userID, jobID, "APPLIED") }
func (s Service) Reject(userID, jobID string) error { return s.Store.Decide(userID, jobID, "REJECTED") }
