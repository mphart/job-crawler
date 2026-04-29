package notify

import (
	"fmt"
	"net/smtp"
	"strings"

	"job-crawler/apps/worker/internal/platform/config"
	"job-crawler/apps/worker/internal/scraper"
)

type Candidate struct {
	UserID    string               `json:"userId"`
	Email     string               `json:"email"`
	Username  string               `json:"username"`
	Frequency string               `json:"frequency"`
	Jobs      []scraper.ScrapedJob `json:"jobs"`
}

func SendDigest(cfg config.Config, candidate Candidate) error {
	if strings.TrimSpace(cfg.SMTPHost) == "" {
		return nil
	}
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)
	subject := "Subject: New jobs matching your preferences\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := buildDigestHTML(candidate.Username, candidate.Frequency, candidate.Jobs)
	msg := []byte(subject + mime + body)
	return smtp.SendMail(addr, auth, cfg.SMTPFrom, []string{candidate.Email}, msg)
}

func buildDigestHTML(username, frequency string, jobs []scraper.ScrapedJob) string {
	var b strings.Builder
	b.WriteString("<h2>Job Crawler Digest</h2>")
	b.WriteString("<p>Hi " + htmlEscape(username) + ", here are fresh jobs from your " + htmlEscape(frequency) + " notification:</p>")
	b.WriteString("<ul>")
	for i, job := range jobs {
		if i >= 12 {
			break
		}
		b.WriteString("<li><strong>" + htmlEscape(job.Title) + "</strong> at " + htmlEscape(job.Company) + " (" + htmlEscape(job.Location) + ")")
		if strings.TrimSpace(job.Compensation) != "" {
			b.WriteString(" - " + htmlEscape(job.Compensation))
		}
		b.WriteString(` - <a href="` + htmlEscape(job.URL) + `">View posting</a></li>`)
	}
	b.WriteString("</ul>")
	b.WriteString("<p>Good luck with your search.</p>")
	return b.String()
}

func htmlEscape(value string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return replacer.Replace(value)
}
