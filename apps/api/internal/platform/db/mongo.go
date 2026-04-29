package db

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

type User struct {
	ID, Email, Username, PasswordHash string
	IsPrivate                         bool
}
type Preference struct {
	Keywords, Locations, DesiredTitles []string
	MinComp                            int
	EmailOptIn, DarkMode               bool
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
	GeneratedJobCount     int
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
	now := time.Now()
	s.Jobs["job_1"] = JobPosting{ID: "job_1", Company: "Cisco", Title: "Software Engineer", Location: "Austin, TX", Compensation: "$130k-$155k", PostedAt: now.Add(-24 * time.Hour).Format(time.RFC3339), URL: "https://example.com/jobs/1", AppliedBy: []AppliedBy{{UserID: "u_2", Username: "alex"}}}
	s.Jobs["job_2"] = JobPosting{ID: "job_2", Company: "365 Retail Markets", Title: "Frontend Engineer", Location: "Remote", Compensation: "$120k-$140k", PostedAt: now.Add(-48 * time.Hour).Format(time.RFC3339), URL: "https://example.com/jobs/2", AppliedBy: nil}
	return s
}

func (s *Store) CreateUser(email, username, passwordHash string, keywords []string) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := "u_" + strings.ReplaceAll(strings.ToLower(username), " ", "-")
	if id == "u_" {
		id = "u_new"
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
	return out
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
	p := Profile{ID: u.ID, Username: u.Username, Email: u.Email, IsPrivate: u.IsPrivate, Preferences: s.Preferences[userID]}
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
		_ = v
	}
	if raw, ok := patch["preferences"].(map[string]any); ok {
		pref := s.Preferences[userID]
		if v, ok := raw["emailOptIn"].(bool); ok {
			pref.EmailOptIn = v
		}
		s.Preferences[userID] = pref
	}
	p := Profile{ID: u.ID, Username: u.Username, Email: u.Email, IsPrivate: u.IsPrivate, Preferences: s.Preferences[userID]}
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

func (s *Store) AddGeneratedJob() JobPosting {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.GeneratedJobCount++
	id := "job_generated_" + time.Now().Format("20060102150405") + "_" + strings.ToLower(time.Now().Format("150405.000000000"))
	if _, exists := s.Jobs[id]; exists {
		id = id + "_" + strings.ToLower(strings.ReplaceAll(time.Now().Format("150405.000000000"), ".", ""))
	}

	companies := []string{"Acme Systems", "Northstar Labs", "Blue River Tech"}
	titles := []string{"Backend Engineer", "Platform Engineer", "Full Stack Engineer"}
	locations := []string{"Remote", "Austin, TX", "Minneapolis, MN"}
	idx := s.GeneratedJobCount % len(companies)

	posting := JobPosting{
		ID:           id,
		Company:      companies[idx],
		Title:        titles[idx],
		Location:     locations[idx],
		Compensation: "$125k-$165k",
		PostedAt:     time.Now().Format(time.RFC3339),
		URL:          "https://example.com/jobs/" + id,
		AppliedBy:    nil,
	}

	s.Jobs[id] = posting
	return posting
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
	}
	return out
}
