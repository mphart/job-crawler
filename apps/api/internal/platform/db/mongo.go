package db

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

type User struct {
	ID, Email, Username, PasswordHash string
	IsPrivate                         bool
	ResumeFileName                    string
}
type Preference struct {
	Keywords, Locations, DesiredTitles, PreferredCompanies []string
	MinComp                                                int
	EmailOptIn, DarkMode                                   bool
}
type AppliedBy struct{ UserID, Username string }
type JobPosting struct {
	ID, Company, Title, Location, Compensation, PostedAt, URL string
	AppliedBy                                                 []AppliedBy
}
type FeedDecision struct {
	UserID, JobID, DecisionType string
	DecisionAt                  string
}

type Profile struct {
	ID, Username, Email string
	IsPrivate           bool
	TotalApplied        int
	ResumeFileName      string
	Preferences         Preference
	AppliedJobs         []JobPosting
}

type Store struct {
	mu                    sync.RWMutex
	Users                 map[string]User
	Preferences           map[string]Preference
	Jobs                  map[string]JobPosting
	Decisions             []FeedDecision
	NotificationFrequency map[string]string
}

type ScrapedJob struct {
	Source       string `json:"source"`
	ExternalID   string `json:"externalId"`
	Company      string `json:"company"`
	Title        string `json:"title"`
	Location     string `json:"location"`
	Compensation string `json:"compensation"`
	PostedAt     string `json:"postedAt"`
	URL          string `json:"url"`
}

func NewStore() *Store {
	s := &Store{Users: map[string]User{}, Preferences: map[string]Preference{}, Jobs: map[string]JobPosting{}, NotificationFrequency: map[string]string{}}
	s.Users["u_1"] = User{ID: "u_1", Email: "mason@example.com", Username: "mason", PasswordHash: "", IsPrivate: false}
	s.Users["u_2"] = User{ID: "u_2", Email: "alex@example.com", Username: "alex", PasswordHash: "", IsPrivate: false}
	s.Preferences["u_1"] = Preference{Keywords: []string{"software engineer", "frontend"}, Locations: []string{"Remote", "Austin"}, DesiredTitles: []string{"Software Engineer", "Frontend Engineer"}, MinComp: 120000, EmailOptIn: true}
	return s
}

func (s *Store) CreateUser(email, username, passwordHash string, keywords []string) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	var id string
	for {
		var buf [10]byte
		_, _ = rand.Read(buf[:])
		id = "u_" + hex.EncodeToString(buf[:])
		if _, taken := s.Users[id]; !taken {
			break
		}
	}
	u := User{ID: id, Email: email, Username: username, PasswordHash: passwordHash}
	s.Users[id] = u
	s.Preferences[id] = Preference{Keywords: keywords, EmailOptIn: true}
	return u
}

func (s *Store) FindUserByEmail(email string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.Users {
		if strings.EqualFold(u.Email, email) {
			return u, true
		}
	}
	return User{}, false
}

func (s *Store) Feed(userID, search, sortBy string) []JobPosting {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hidden := map[string]bool{}
	for _, d := range s.Decisions {
		if d.UserID == userID {
			hidden[d.JobID] = true
		}
	}
	q := strings.ToLower(strings.TrimSpace(search))
	out := make([]JobPosting, 0, len(s.Jobs))
	for _, j := range s.Jobs {
		if hidden[j.ID] {
			continue
		}
		hay := strings.ToLower(j.Company + " " + j.Title + " " + j.Location)
		if q != "" && !strings.Contains(hay, q) {
			continue
		}
		out = append(out, j)
	}
	if q == "" {
		pref := s.Preferences[userID]
		roleSignals := make([]string, 0, len(pref.Keywords)+len(pref.DesiredTitles))
		roleSignals = append(roleSignals, pref.Keywords...)
		roleSignals = append(roleSignals, pref.DesiredTitles...)
		filtered := make([]JobPosting, 0, len(out))
		for _, job := range out {
			if inMemoryIsRelevant(job, roleSignals, pref.PreferredCompanies) {
				filtered = append(filtered, job)
			}
		}
		out = filtered
	}
	sort.Slice(out, func(i, j int) bool {
		switch sortBy {
		case "company":
			return out[i].Company < out[j].Company
		case "title":
			return out[i].Title < out[j].Title
		case "location":
			return out[i].Location < out[j].Location
		case "money":
			return out[i].Compensation > out[j].Compensation
		default:
			return out[i].PostedAt > out[j].PostedAt
		}
	})
	if q == "" {
		pref := s.Preferences[userID]
		if len(pref.PreferredCompanies) > 0 {
			sort.SliceStable(out, func(i, j int) bool {
				return inMemoryPreferredCompanyScore(out[i], pref.PreferredCompanies) > inMemoryPreferredCompanyScore(out[j], pref.PreferredCompanies)
			})
		}
	}
	return out
}

func inMemoryIsRelevant(job JobPosting, roleSignals []string, preferredCompanies []string) bool {
	roleMatch := len(roleSignals) == 0
	hay := strings.ToLower(job.Company + " " + job.Title)
	for _, signal := range roleSignals {
		s := strings.ToLower(strings.TrimSpace(signal))
		if s != "" && strings.Contains(hay, s) {
			roleMatch = true
			break
		}
	}
	preferred := false
	company := strings.ToLower(strings.TrimSpace(job.Company))
	for _, c := range preferredCompanies {
		p := strings.ToLower(strings.TrimSpace(c))
		if p != "" && strings.Contains(company, p) {
			preferred = true
			break
		}
	}
	if preferred {
		return roleMatch
	}
	return roleMatch || len(roleSignals) == 0
}

func inMemoryPreferredCompanyScore(job JobPosting, preferredCompanies []string) int {
	company := strings.ToLower(strings.TrimSpace(job.Company))
	for _, preferred := range preferredCompanies {
		p := strings.ToLower(strings.TrimSpace(preferred))
		if p != "" && strings.Contains(company, p) {
			return 1
		}
	}
	return 0
}

func (s *Store) Decide(userID, jobID, decision string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Jobs[jobID]; !ok {
		return errors.New("job not found")
	}
	s.Decisions = append(s.Decisions, FeedDecision{UserID: userID, JobID: jobID, DecisionType: decision, DecisionAt: time.Now().Format(time.RFC3339)})
	return nil
}

func (s *Store) Profile(requester, userID string) (Profile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.Users[userID]
	if !ok {
		return Profile{}, false
	}
	p := Profile{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		IsPrivate:      u.IsPrivate,
		ResumeFileName: u.ResumeFileName,
		Preferences:    s.Preferences[userID],
	}
	for _, d := range s.Decisions {
		if d.UserID == userID && d.DecisionType == "APPLIED" {
			if j, ok := s.Jobs[d.JobID]; ok {
				j.PostedAt = d.DecisionAt
				p.AppliedJobs = append(p.AppliedJobs, j)
				p.TotalApplied++
			}
		}
	}
	if requester != userID && p.IsPrivate {
		p.AppliedJobs = nil
	}
	return p, true
}

func mongoStringSlice(v any) []string {
	switch t := v.(type) {
	case []string:
		out := make([]string, 0, len(t))
		for _, s := range t {
			s = strings.TrimSpace(s)
			if s != "" {
				out = append(out, s)
			}
		}
		return out
	case []any:
		out := make([]string, 0, len(t))
		for _, x := range t {
			if s, ok := x.(string); ok {
				s = strings.TrimSpace(s)
				if s != "" {
					out = append(out, s)
				}
			}
		}
		return out
	default:
		return nil
	}
}

func (s *Store) UpdateMe(userID string, patch map[string]any) (Profile, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.Users[userID]
	if !ok {
		return Profile{}, false
	}
	if v, ok := patch["isPrivate"].(bool); ok {
		u.IsPrivate = v
	}
	s.Users[userID] = u
	if v, ok := patch["resumeFileName"].(string); ok {
		u.ResumeFileName = v
		s.Users[userID] = u
	}
	if v, ok := patch["resumeContentBase64"].(string); ok {
		_ = v
	}
	if raw, ok := patch["preferences"].(map[string]any); ok {
		pref := s.Preferences[userID]
		if v, ok := raw["emailOptIn"].(bool); ok {
			pref.EmailOptIn = v
		}
		if v, ok := raw["darkMode"].(bool); ok {
			pref.DarkMode = v
		}
		switch v := raw["minComp"].(type) {
		case float64:
			pref.MinComp = int(v)
		case int:
			pref.MinComp = v
		case int64:
			pref.MinComp = int(v)
		}
		if v, ok := raw["keywords"]; ok {
			pref.Keywords = mongoStringSlice(v)
		}
		if v, ok := raw["locations"]; ok {
			pref.Locations = mongoStringSlice(v)
		}
		if v, ok := raw["desiredTitles"]; ok {
			pref.DesiredTitles = mongoStringSlice(v)
		}
		if v, ok := raw["preferredCompanies"]; ok {
			pref.PreferredCompanies = mongoStringSlice(v)
		}
		s.Preferences[userID] = pref
	}
	u = s.Users[userID]
	p := Profile{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		IsPrivate:      u.IsPrivate,
		ResumeFileName: u.ResumeFileName,
		Preferences:    s.Preferences[userID],
	}
	for _, d := range s.Decisions {
		if d.UserID == userID && d.DecisionType == "APPLIED" {
			if j, ok := s.Jobs[d.JobID]; ok {
				j.PostedAt = d.DecisionAt
				p.AppliedJobs = append(p.AppliedJobs, j)
				p.TotalApplied++
			}
		}
	}
	return p, true
}

func (s *Store) SearchUsers(q string) []map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := strings.ToLower(strings.TrimSpace(q))
	out := []map[string]any{}
	for _, u := range s.Users {
		if n == "" || strings.Contains(strings.ToLower(u.Username), n) {
			total := 0
			for _, d := range s.Decisions {
				if d.UserID == u.ID && d.DecisionType == "APPLIED" {
					total++
				}
			}
			out = append(out, map[string]any{"id": u.ID, "username": u.Username, "totalApplied": total})
		}
	}
	return out
}

func (s *Store) SearchCompanies(q string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	needle := strings.ToLower(strings.TrimSpace(q))
	seen := map[string]struct{}{}
	out := make([]string, 0, 25)
	for _, job := range s.Jobs {
		company := strings.TrimSpace(job.Company)
		if company == "" {
			continue
		}
		key := strings.ToLower(company)
		if needle != "" && !strings.Contains(key, needle) {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, company)
		if len(out) >= 25 {
			break
		}
	}
	sort.Strings(out)
	return out
}

func (s *Store) IngestScrapedJobs(jobs []ScrapedJob) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	inserted := 0
	for _, job := range jobs {
		if strings.TrimSpace(job.URL) == "" || strings.TrimSpace(job.Title) == "" {
			continue
		}
		exists := false
		for _, existing := range s.Jobs {
			if existing.URL == job.URL {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		id := "job_scraped_" + strings.ReplaceAll(strings.ToLower(time.Now().Format("20060102150405.000000000")), ".", "") + "_" + strings.ReplaceAll(strings.ToLower(job.Source), " ", "_")
		if _, ok := s.Jobs[id]; ok {
			id = id + "_" + strings.ReplaceAll(strings.ToLower(strings.TrimSpace(job.ExternalID)), " ", "_")
		}

		postedAt := job.PostedAt
		if strings.TrimSpace(postedAt) == "" {
			postedAt = time.Now().Format(time.RFC3339)
		}

		s.Jobs[id] = JobPosting{
			ID:           id,
			Company:      job.Company,
			Title:        job.Title,
			Location:     job.Location,
			Compensation: job.Compensation,
			PostedAt:     postedAt,
			URL:          job.URL,
			AppliedBy:    nil,
		}
		inserted++
	}
	return inserted
}

func (s *Store) WorkerKeywords() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	seen := map[string]struct{}{}
	out := make([]string, 0)
	for _, pref := range s.Preferences {
		for _, keyword := range pref.Keywords {
			normalized := strings.ToLower(strings.TrimSpace(keyword))
			if normalized == "" {
				continue
			}
			if _, ok := seen[normalized]; ok {
				continue
			}
			seen[normalized] = struct{}{}
			out = append(out, normalized)
		}
		for _, company := range pref.PreferredCompanies {
			normalized := strings.ToLower(strings.TrimSpace(company))
			if normalized == "" {
				continue
			}
			if _, ok := seen[normalized]; ok {
				continue
			}
			seen[normalized] = struct{}{}
			out = append(out, normalized)
		}
	}
	return out
}
