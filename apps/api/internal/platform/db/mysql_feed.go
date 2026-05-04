package db

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"
)

func (s *MySQLAuthStore) ensureJobSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS jobs (
  id VARCHAR(128) PRIMARY KEY,
  source VARCHAR(64) NOT NULL,
  external_id VARCHAR(128) NULL,
  company VARCHAR(255) NOT NULL,
  title VARCHAR(255) NOT NULL,
  location VARCHAR(255) NULL,
  compensation VARCHAR(128) NULL,
  posted_at DATETIME NOT NULL,
  url TEXT NOT NULL,
  UNIQUE KEY uniq_job_url (url(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS feed_decisions (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id VARCHAR(64) NOT NULL,
  job_id VARCHAR(128) NOT NULL,
  decision_type VARCHAR(16) NOT NULL,
  decision_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uniq_user_job_decision (user_id, job_id),
  KEY idx_decisions_user (user_id),
  KEY idx_decisions_job (job_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS notification_settings (
  user_id VARCHAR(64) PRIMARY KEY,
  frequency VARCHAR(32) NOT NULL DEFAULT 'daily',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS notification_deliveries (
  user_id VARCHAR(64) PRIMARY KEY,
  last_sent_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
	}
	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

type DigestCandidate struct {
	UserID    string       `json:"userId"`
	Email     string       `json:"email"`
	Username  string       `json:"username"`
	Frequency string       `json:"frequency"`
	Jobs      []JobPosting `json:"jobs"`
}

func (s *MySQLAuthStore) IngestScrapedJobs(jobs []ScrapedJob) (int, error) {
	inserted := 0
	for _, job := range jobs {
		if strings.TrimSpace(job.URL) == "" || strings.TrimSpace(job.Title) == "" {
			continue
		}
		id := "job_" + strings.ReplaceAll(strings.ToLower(job.Source), " ", "_") + "_" + strings.ReplaceAll(strings.ToLower(job.ExternalID), " ", "_")
		if strings.TrimSpace(job.ExternalID) == "" {
			id = "job_" + strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%d", time.Now().UnixNano())), " ", "_")
		}

		postedAt := time.Now()
		if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(job.PostedAt)); err == nil {
			postedAt = parsed
		}

		result, err := s.db.Exec(`
INSERT INTO jobs (id, source, external_id, company, title, location, compensation, posted_at, url)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  company = VALUES(company),
  title = VALUES(title),
  location = VALUES(location),
  compensation = VALUES(compensation),
  posted_at = VALUES(posted_at)
`, id, job.Source, job.ExternalID, job.Company, job.Title, job.Location, job.Compensation, postedAt, job.URL)
		if err != nil {
			return inserted, err
		}
		affected, _ := result.RowsAffected()
		if affected > 0 {
			inserted++
		}
	}
	return inserted, nil
}

func (s *MySQLAuthStore) Feed(userID, search, sortBy string) ([]JobPosting, error) {
	sortSQL := "j.posted_at DESC"
	switch sortBy {
	case "company":
		sortSQL = "j.company ASC"
	case "title":
		sortSQL = "j.title ASC"
	case "location":
		sortSQL = "j.location ASC"
	case "money":
		sortSQL = "j.compensation DESC"
	}

	query := `
SELECT j.id, j.company, j.title, j.location, j.compensation, j.posted_at, j.url
FROM jobs j
LEFT JOIN feed_decisions d
  ON d.job_id = j.id AND d.user_id = ?
WHERE d.id IS NULL
  AND (? = '' OR LOWER(CONCAT(j.company, ' ', j.title, ' ', j.location)) LIKE CONCAT('%', LOWER(?), '%'))
ORDER BY ` + sortSQL + `
LIMIT 200`

	rows, err := s.db.Query(query, userID, strings.TrimSpace(search), strings.TrimSpace(search))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]JobPosting, 0)
	for rows.Next() {
		var job JobPosting
		var postedAt time.Time
		var location, compensation sql.NullString
		if err := rows.Scan(&job.ID, &job.Company, &job.Title, &location, &compensation, &postedAt, &job.URL); err != nil {
			return nil, err
		}
		job.Location = location.String
		job.Compensation = compensation.String
		job.PostedAt = postedAt.Format(time.RFC3339)
		job.AppliedBy = []AppliedBy{}
		out = append(out, job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	signals, err := s.userSignals(userID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(search) == "" && len(signals) > 0 {
		sort.SliceStable(out, func(i, j int) bool {
			return recommendationScore(out[i], signals) > recommendationScore(out[j], signals)
		})
	}
	return out, nil
}

func (s *MySQLAuthStore) Decide(userID, jobID, decision string) error {
	_, err := s.db.Exec(`
INSERT INTO feed_decisions (user_id, job_id, decision_type, decision_at)
VALUES (?, ?, ?, NOW())
ON DUPLICATE KEY UPDATE decision_type = VALUES(decision_type), decision_at = VALUES(decision_at)
`, userID, jobID, decision)
	return err
}

func (s *MySQLAuthStore) UpdateNotification(userID string, payload map[string]any) (map[string]any, error) {
	frequency := "daily"
	if v, ok := payload["frequency"].(string); ok && strings.TrimSpace(v) != "" {
		frequency = NormalizeNotificationFrequency(v)
	}
	emailOptIn := true
	if v, ok := payload["emailOptIn"].(bool); ok {
		emailOptIn = v
	}

	if _, err := s.db.Exec("UPDATE users SET email_opt_in = ? WHERE id = ?", emailOptIn, userID); err != nil {
		return nil, err
	}
	if _, err := s.db.Exec(`
INSERT INTO notification_settings (user_id, frequency) VALUES (?, ?)
ON DUPLICATE KEY UPDATE frequency = VALUES(frequency)
`, userID, frequency); err != nil {
		return nil, err
	}
	return map[string]any{"emailOptIn": emailOptIn, "frequency": frequency}, nil
}

func (s *MySQLAuthStore) GetNotification(userID string) (map[string]any, error) {
	var emailOptIn bool
	if err := s.db.QueryRow("SELECT email_opt_in FROM users WHERE id = ? LIMIT 1", userID).Scan(&emailOptIn); err != nil {
		if err == sql.ErrNoRows {
			return map[string]any{"emailOptIn": true, "frequency": "daily"}, nil
		}
		return nil, err
	}

	frequency := "daily"
	_ = s.db.QueryRow("SELECT frequency FROM notification_settings WHERE user_id = ? LIMIT 1", userID).Scan(&frequency)
	if strings.TrimSpace(frequency) == "" {
		frequency = "daily"
	}
	frequency = NormalizeNotificationFrequency(frequency)
	return map[string]any{"emailOptIn": emailOptIn, "frequency": frequency}, nil
}

func (s *MySQLAuthStore) DigestCandidates(now time.Time) ([]DigestCandidate, error) {
	rows, err := s.db.Query(`
SELECT u.id, u.email, u.username, u.keywords, IFNULL(ns.frequency, 'daily')
FROM users u
LEFT JOIN notification_settings ns ON ns.user_id = u.id
WHERE u.email_opt_in = TRUE
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	candidates := make([]DigestCandidate, 0)
	for rows.Next() {
		var userID, email, username, keywordCSV, frequency string
		if err := rows.Scan(&userID, &email, &username, &keywordCSV, &frequency); err != nil {
			return nil, err
		}
		lastSent := time.Time{}
		_ = s.db.QueryRow("SELECT last_sent_at FROM notification_deliveries WHERE user_id = ? LIMIT 1", userID).Scan(&lastSent)
		if !isDueForFrequency(lastSent, now, frequency) {
			continue
		}

		jobs, err := s.matchJobsForKeywords(userID, keywordCSV)
		if err != nil {
			return nil, err
		}
		if len(jobs) == 0 {
			continue
		}

		candidates = append(candidates, DigestCandidate{
			UserID: userID, Email: email, Username: username, Frequency: NormalizeNotificationFrequency(frequency), Jobs: jobs,
		})
	}
	return candidates, rows.Err()
}

func (s *MySQLAuthStore) MarkNotificationSent(userID string, sentAt time.Time) error {
	_, err := s.db.Exec(`
INSERT INTO notification_deliveries (user_id, last_sent_at)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE last_sent_at = VALUES(last_sent_at)
`, userID, sentAt)
	return err
}

func (s *MySQLAuthStore) matchJobsForKeywords(userID, keywordCSV string) ([]JobPosting, error) {
	keywords, err := s.userSignals(userID)
	if err != nil {
		return nil, err
	}
	if len(keywords) == 0 {
		keywords = splitKeywords(keywordCSV)
	}
	rows, err := s.db.Query(`
SELECT id, company, title, location, compensation, posted_at, url
FROM jobs
ORDER BY posted_at DESC
LIMIT 250`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]JobPosting, 0, 12)
	for rows.Next() {
		var job JobPosting
		var postedAt time.Time
		var location, compensation sql.NullString
		if err := rows.Scan(&job.ID, &job.Company, &job.Title, &location, &compensation, &postedAt, &job.URL); err != nil {
			return nil, err
		}
		job.Location = location.String
		job.Compensation = compensation.String
		job.PostedAt = postedAt.Format(time.RFC3339)
		if len(keywords) == 0 || jobMatchesAny(job, keywords) {
			matches = append(matches, job)
		}
		if len(matches) >= 15 {
			break
		}
	}
	return matches, rows.Err()
}

func (s *MySQLAuthStore) userSignals(userID string) ([]string, error) {
	merged := map[string]struct{}{}
	var keywordCSV string
	_ = s.db.QueryRow("SELECT IFNULL(keywords, '') FROM users WHERE id = ? LIMIT 1", userID).Scan(&keywordCSV)
	for _, k := range splitKeywords(keywordCSV) {
		merged[k] = struct{}{}
	}
	var resumeKeywords, resumeRoles, resumeLocations string
	_ = s.db.QueryRow(`
SELECT IFNULL(keywords, ''), IFNULL(role_families, ''), IFNULL(locations, '')
FROM resume_signals WHERE user_id = ? AND parse_status = 'ready' LIMIT 1
`, userID).Scan(&resumeKeywords, &resumeRoles, &resumeLocations)
	for _, set := range []string{resumeKeywords, resumeRoles, resumeLocations} {
		for _, k := range splitKeywords(set) {
			merged[k] = struct{}{}
		}
	}
	out := make([]string, 0, len(merged))
	for k := range merged {
		out = append(out, k)
	}
	sort.Strings(out)
	return out, nil
}

func splitKeywords(csv string) []string {
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		k := strings.ToLower(strings.TrimSpace(p))
		if k != "" {
			out = append(out, k)
		}
	}
	return out
}

func jobMatchesAny(job JobPosting, keywords []string) bool {
	hay := strings.ToLower(job.Company + " " + job.Title + " " + job.Location)
	for _, k := range keywords {
		if strings.Contains(hay, k) {
			return true
		}
	}
	return false
}

// NormalizeNotificationFrequency maps legacy values and defaults unknowns to "daily".
func NormalizeNotificationFrequency(raw string) string {
	f := strings.TrimSpace(strings.ToLower(raw))
	switch f {
	case "daily", "twice-daily", "weekly":
		return f
	case "instant":
		return "daily"
	case "every-2-weeks":
		return "weekly"
	default:
		return "daily"
	}
}

func (s *MySQLAuthStore) SetNotificationFrequency(userID, rawFrequency string) error {
	f := NormalizeNotificationFrequency(rawFrequency)
	_, err := s.db.Exec(`
INSERT INTO notification_settings (user_id, frequency) VALUES (?, ?)
ON DUPLICATE KEY UPDATE frequency = VALUES(frequency)
`, userID, f)
	return err
}

func recommendationScore(job JobPosting, keywords []string) int {
	score := 0
	hay := strings.ToLower(job.Company + " " + job.Title + " " + job.Location)
	for _, k := range keywords {
		if strings.Contains(hay, k) {
			score += 5
		}
	}
	if strings.TrimSpace(job.Compensation) != "" {
		score++
	}
	return score
}

func isDueForFrequency(lastSent, now time.Time, frequency string) bool {
	if lastSent.IsZero() {
		return true
	}
	f := NormalizeNotificationFrequency(frequency)
	var interval time.Duration
	switch f {
	case "twice-daily":
		interval = 12 * time.Hour
	case "weekly":
		interval = 7 * 24 * time.Hour
	default:
		interval = 24 * time.Hour
	}
	return now.Sub(lastSent) >= interval
}
