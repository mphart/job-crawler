package worker

import (
	"time"

	"job-crawler/apps/api/internal/platform/db"
)

type Service struct {
	Store *db.Store
	MySQL *db.MySQLAuthStore
}

func (s Service) Tick() db.JobPosting {
	return s.Store.AddGeneratedJob()
}

func (s Service) Ingest(jobs []db.ScrapedJob) int {
	if s.MySQL != nil {
		if inserted, err := s.MySQL.IngestScrapedJobs(jobs); err == nil {
			return inserted
		}
	}
	return s.Store.IngestScrapedJobs(jobs)
}

func (s Service) Keywords() []string {
	if s.MySQL != nil {
		if keywords, err := s.MySQL.WorkerKeywords(); err == nil && len(keywords) > 0 {
			return keywords
		}
	}
	return s.Store.WorkerKeywords()
}

func (s Service) DigestCandidates() ([]db.DigestCandidate, error) {
	if s.MySQL == nil {
		return []db.DigestCandidate{}, nil
	}
	return s.MySQL.DigestCandidates(nowUTC())
}

func (s Service) MarkNotificationSent(userID string) error {
	if s.MySQL == nil {
		return nil
	}
	return s.MySQL.MarkNotificationSent(userID, nowUTC())
}

var nowUTC = func() time.Time { return time.Now().UTC() }
