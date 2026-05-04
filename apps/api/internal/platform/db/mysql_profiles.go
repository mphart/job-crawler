package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

func (s *MySQLAuthStore) ensureProfileColumns() error {
	queries := []string{
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS is_private BOOLEAN NOT NULL DEFAULT FALSE",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS resume_file_name VARCHAR(255) NULL",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS resume_content_base64 LONGTEXT NULL",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS email_opt_in BOOLEAN NOT NULL DEFAULT TRUE",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS dark_mode BOOLEAN NOT NULL DEFAULT FALSE",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS min_comp INT NOT NULL DEFAULT 0",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS keywords TEXT NULL",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS locations TEXT NULL",
		"ALTER TABLE users ADD COLUMN IF NOT EXISTS desired_titles TEXT NULL",
	}
	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func (s *MySQLAuthStore) WorkerKeywords() ([]string, error) {
	rows, err := s.db.Query("SELECT keywords FROM users WHERE keywords IS NOT NULL AND keywords <> ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seen := map[string]struct{}{}
	out := make([]string, 0)
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		for _, part := range strings.Split(raw, ",") {
			keyword := strings.ToLower(strings.TrimSpace(part))
			if keyword == "" {
				continue
			}
			if _, ok := seen[keyword]; ok {
				continue
			}
			seen[keyword] = struct{}{}
			out = append(out, keyword)
		}
	}
	return out, rows.Err()
}

func splitCommaPreserve(csv string) []string {
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func joinCommaLower(parts []string) string {
	b := strings.Builder{}
	for _, p := range parts {
		p = strings.ToLower(strings.TrimSpace(p))
		if p == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString(",")
		}
		b.WriteString(p)
	}
	return b.String()
}

func joinCommaPreserve(parts []string) string {
	b := strings.Builder{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString(",")
		}
		b.WriteString(p)
	}
	return b.String()
}

func stringSliceFromAny(v any) ([]string, bool) {
	switch t := v.(type) {
	case []string:
		return t, true
	case []any:
		out := make([]string, 0, len(t))
		for _, x := range t {
			if s, ok := x.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out, true
	default:
		return nil, false
	}
}

func intFromAny(v any) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	case int64:
		return int(t), true
	case json.Number:
		i64, err := t.Int64()
		if err != nil {
			return 0, false
		}
		return int(i64), true
	default:
		return 0, false
	}
}

func (s *MySQLAuthStore) Profile(requester, userID string) (Profile, bool, error) {
	row := s.db.QueryRow(`
SELECT id, email, username, is_private, resume_file_name, email_opt_in, dark_mode, min_comp,
       IFNULL(keywords, ''), IFNULL(locations, ''), IFNULL(desired_titles, '')
FROM users WHERE id = ? LIMIT 1
`, userID)

	var id, email, username string
	var isPrivate, emailOptIn, darkMode bool
	var resume sql.NullString
	var minComp int
	var keywordCSV, locationCSV, titlesCSV string

	if err := row.Scan(&id, &email, &username, &isPrivate, &resume, &emailOptIn, &darkMode, &minComp, &keywordCSV, &locationCSV, &titlesCSV); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Profile{}, false, nil
		}
		return Profile{}, false, err
	}

	p := Profile{
		ID:             id,
		Username:       username,
		Email:          email,
		IsPrivate:      isPrivate,
		ResumeFileName: resume.String,
		Preferences: Preference{
			Keywords:      splitKeywords(keywordCSV),
			Locations:     splitCommaPreserve(locationCSV),
			DesiredTitles: splitCommaPreserve(titlesCSV),
			MinComp:       minComp,
			EmailOptIn:    emailOptIn,
			DarkMode:      darkMode,
		},
		AppliedJobs:  []JobPosting{},
		TotalApplied: 0,
	}
	appliedRows, err := s.db.Query(`
SELECT j.id, j.company, j.title, j.location, j.compensation, d.decision_at, j.url
FROM feed_decisions d
JOIN jobs j ON j.id = d.job_id
WHERE d.user_id = ? AND d.decision_type = 'APPLIED'
ORDER BY d.decision_at DESC
`, userID)
	if err != nil {
		return Profile{}, false, err
	}
	defer appliedRows.Close()
	for appliedRows.Next() {
		var job JobPosting
		var location, compensation sql.NullString
		var decisionAt time.Time
		if err := appliedRows.Scan(&job.ID, &job.Company, &job.Title, &location, &compensation, &decisionAt, &job.URL); err != nil {
			return Profile{}, false, err
		}
		job.Location = location.String
		job.Compensation = compensation.String
		job.PostedAt = decisionAt.Format(time.RFC3339)
		p.AppliedJobs = append(p.AppliedJobs, job)
		p.TotalApplied++
	}
	if err := appliedRows.Err(); err != nil {
		return Profile{}, false, err
	}
	if requester != userID && p.IsPrivate {
		p.AppliedJobs = nil
	}
	return p, true, nil
}

func (s *MySQLAuthStore) UpdateMe(userID string, patch map[string]any) (Profile, bool, error) {
	assignments := make([]string, 0, 12)
	args := make([]any, 0, 14)

	if v, ok := patch["isPrivate"].(bool); ok {
		assignments = append(assignments, "is_private = ?")
		args = append(args, v)
	}
	if v, ok := patch["resumeFileName"].(string); ok {
		assignments = append(assignments, "resume_file_name = ?")
		args = append(args, v)
	}
	if v, ok := patch["resumeContentBase64"].(string); ok {
		assignments = append(assignments, "resume_content_base64 = ?")
		args = append(args, v)
	}
	if pref, ok := patch["preferences"].(map[string]any); ok {
		if v, ok := pref["emailOptIn"].(bool); ok {
			assignments = append(assignments, "email_opt_in = ?")
			args = append(args, v)
		}
		if v, ok := pref["darkMode"].(bool); ok {
			assignments = append(assignments, "dark_mode = ?")
			args = append(args, v)
		}
		if v, ok := pref["minComp"]; ok {
			if n, ok2 := intFromAny(v); ok2 {
				assignments = append(assignments, "min_comp = ?")
				args = append(args, n)
			}
		}
		if v, ok := pref["keywords"]; ok {
			if parts, ok2 := stringSliceFromAny(v); ok2 {
				csv := joinCommaLower(parts)
				assignments = append(assignments, "keywords = ?")
				args = append(args, csv)
			}
		}
		if v, ok := pref["locations"]; ok {
			if parts, ok2 := stringSliceFromAny(v); ok2 {
				csv := joinCommaPreserve(parts)
				assignments = append(assignments, "locations = ?")
				args = append(args, csv)
			}
		}
		if v, ok := pref["desiredTitles"]; ok {
			if parts, ok2 := stringSliceFromAny(v); ok2 {
				csv := joinCommaPreserve(parts)
				assignments = append(assignments, "desired_titles = ?")
				args = append(args, csv)
			}
		}
	}

	if len(assignments) > 0 {
		query := "UPDATE users SET " + strings.Join(assignments, ", ") + " WHERE id = ?"
		args = append(args, userID)
		if _, err := s.db.Exec(query, args...); err != nil {
			return Profile{}, false, err
		}
	}

	p, ok, err := s.Profile(userID, userID)
	return p, ok, err
}

func (s *MySQLAuthStore) SearchUsers(q string) ([]map[string]any, error) {
	query := `
SELECT id, username
FROM users
WHERE ? = '' OR LOWER(username) LIKE CONCAT('%', LOWER(?), '%')
ORDER BY username ASC
LIMIT 25`
	rows, err := s.db.Query(query, strings.TrimSpace(q), strings.TrimSpace(q))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]map[string]any, 0)
	for rows.Next() {
		var id, username string
		if err := rows.Scan(&id, &username); err != nil {
			return nil, err
		}
		var appliedCount int
		_ = s.db.QueryRow(`
SELECT COUNT(*) FROM feed_decisions WHERE user_id = ? AND decision_type = 'APPLIED'
`, id).Scan(&appliedCount)
		result = append(result, map[string]any{
			"id":           id,
			"username":     username,
			"totalApplied": appliedCount,
		})
	}
	return result, rows.Err()
}
