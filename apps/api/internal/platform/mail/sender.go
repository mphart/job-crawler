package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Config struct {
	Host, Username, Password, From string
	Port                           int
}

type Sender struct{ cfg Config }

func NewSender(cfg Config) *Sender {
	if cfg.Port <= 0 {
		cfg.Port = 587
	}
	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.Username = strings.TrimSpace(cfg.Username)
	cfg.Password = strings.TrimSpace(cfg.Password)
	cfg.From = strings.TrimSpace(cfg.From)
	return &Sender{cfg: cfg}
}

// HasHost reports whether an SMTP relay hostname is configured (other fields may still be missing).
func (s *Sender) HasHost() bool {
	return s.cfg.Host != ""
}

// Enabled is true when HasHost is true (legacy callers); prefer ValidateSMTP before sending.
func (s *Sender) Enabled() bool {
	return s.HasHost()
}

// ValidateSMTP returns an error if the relay is partially configured or the From address will be rejected by providers like Brevo.
func (s *Sender) ValidateSMTP() error {
	if !s.HasHost() {
		return fmt.Errorf("SMTP host is empty")
	}
	if s.cfg.Username == "" {
		return fmt.Errorf("SMTP username is empty (use Brevo SMTP login, e.g. xxxx@smtp-brevo.com)")
	}
	if s.cfg.Password == "" {
		return fmt.Errorf("SMTP password is empty (use the SMTP key from Brevo, not your Brevo login password)")
	}
	if s.cfg.From == "" {
		return fmt.Errorf("SMTP From is empty: set WORKER_SMTP_FROM / API_SMTP_FROM to an email address you have verified under Brevo > Senders & IP > Domains")
	}
	if !strings.Contains(s.cfg.From, "@") {
		return fmt.Errorf("SMTP From %q must be a valid email address", s.cfg.From)
	}
	lower := strings.ToLower(s.cfg.From)
	if strings.HasSuffix(lower, ".local") || strings.HasSuffix(lower, ".test") || strings.HasSuffix(lower, ".invalid") {
		return fmt.Errorf("SMTP From %q is a placeholder; Brevo requires a verified real sender address (set WORKER_SMTP_FROM)", s.cfg.From)
	}
	return nil
}

func (s *Sender) SendHTML(to, subject, htmlBody string) error {
	if !s.HasHost() {
		return nil
	}
	if err := s.ValidateSMTP(); err != nil {
		return err
	}
	to = strings.TrimSpace(to)
	if to == "" || strings.ContainsAny(to, "\r\n") {
		return fmt.Errorf("invalid recipient address")
	}
	if strings.ContainsAny(s.cfg.From, "\r\n") || strings.ContainsAny(subject, "\r\n") {
		return fmt.Errorf("invalid mail headers")
	}
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	headers := "To: " + to + "\r\nFrom: " + s.cfg.From + "\r\nSubject: " + subject + "\r\n"
	mime := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n"
	msg := []byte(headers + mime + htmlBody)
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, msg)
}

// SendWelcomeRegistration sends the formal registration confirmation (no-op if SMTP host is not configured).
func (s *Sender) SendWelcomeRegistration(toEmail, legalName string) error {
	return s.SendHTML(toEmail, welcomeRegistrationSubject(), welcomeRegistrationHTML(legalName))
}
