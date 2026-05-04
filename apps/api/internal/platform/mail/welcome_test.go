package mail

import (
	"strings"
	"testing"
)

func TestWelcomeRegistrationHTML_escapesName(t *testing.T) {
	out := welcomeRegistrationHTML(`Pat <script>alert(1)</script>`)
	if strings.Contains(out, "<script>") {
		t.Fatalf("expected script stripped or escaped")
	}
	if !strings.Contains(out, "&lt;") {
		t.Fatalf("expected escaped angle bracket in name")
	}
}
