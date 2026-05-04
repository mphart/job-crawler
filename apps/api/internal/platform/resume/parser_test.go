package resume

import (
	"encoding/base64"
	"testing"
)

func TestParseBase64_ExtractsHeuristicSignals(t *testing.T) {
	raw := "Experienced Mechanical Engineer with CAD and manufacturing background. Open to Remote and Austin roles."
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))

	signals, err := ParseBase64(encoded)
	if err != nil {
		t.Fatalf("expected parse to succeed, got %v", err)
	}

	if len(signals.Keywords) == 0 {
		t.Fatalf("expected keyword extraction to return at least one signal")
	}

	foundMechanical := false
	for _, k := range signals.Keywords {
		if k == "mechanical engineer" {
			foundMechanical = true
			break
		}
	}
	if !foundMechanical {
		t.Fatalf("expected mechanical engineer keyword, got %v", signals.Keywords)
	}
}
