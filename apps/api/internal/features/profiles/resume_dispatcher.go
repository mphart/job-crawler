package profiles

import (
	"log"

	"job-crawler/apps/api/internal/platform/db"
	"job-crawler/apps/api/internal/platform/resume"
)

type ResumeParseDispatcher struct {
	store *db.MySQLAuthStore
	queue chan parseJob
}

type parseJob struct {
	userID  string
	content string
}

func NewResumeParseDispatcher(store *db.MySQLAuthStore) *ResumeParseDispatcher {
	if store == nil {
		return nil
	}
	d := &ResumeParseDispatcher{
		store: store,
		queue: make(chan parseJob, 128),
	}
	go d.loop()
	return d
}

func (d *ResumeParseDispatcher) Enqueue(userID, resumeContentBase64 string) error {
	d.queue <- parseJob{userID: userID, content: resumeContentBase64}
	_ = d.store.UpsertResumeSignals(db.ResumeSignals{
		UserID:      userID,
		ParseStatus: "pending",
	})
	return nil
}

func (d *ResumeParseDispatcher) loop() {
	for job := range d.queue {
		signals, err := resume.ParseBase64(job.content)
		if err != nil {
			_ = d.store.UpsertResumeSignals(db.ResumeSignals{
				UserID:      job.userID,
				ParseStatus: "failed",
				ParseError:  err.Error(),
			})
			continue
		}
		if err := d.store.UpsertResumeSignals(db.ResumeSignals{
			UserID:       job.userID,
			Keywords:     signals.Keywords,
			RoleFamilies: signals.RoleFamilies,
			Locations:    signals.Locations,
			ParseStatus:  "ready",
			ParseError:   "",
		}); err != nil {
			log.Printf("resume signal upsert failed for %s: %v", job.userID, err)
		}
	}
}
