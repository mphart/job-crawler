package worker

import "job-crawler/apps/api/internal/platform/db"

type Service struct {
    Store *db.Store
}

func (s Service) Tick() db.JobPosting {
    return s.Store.AddGeneratedJob()
}
