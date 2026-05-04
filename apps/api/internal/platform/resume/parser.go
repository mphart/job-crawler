package resume

import (
	"encoding/base64"
	"regexp"
	"strings"
)

type ParsedSignals struct {
	Keywords     []string
	RoleFamilies []string
	Locations    []string
}

func ParseBase64(contentBase64 string) (ParsedSignals, error) {
	data, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return ParsedSignals{}, err
	}

	text := extractPDFText(data)
	normalized := strings.ToLower(text)

	return ParsedSignals{
		Keywords:     topMatches(normalized, keywordDictionary()),
		RoleFamilies: topMatches(normalized, roleFamilyDictionary()),
		Locations:    topMatches(normalized, locationDictionary()),
	}, nil
}

func extractPDFText(data []byte) string {
	// Heuristic extraction: keep printable text runs from decoded bytes.
	text := string(data)
	runRegex := regexp.MustCompile(`[A-Za-z0-9\+\#\.\-_,/\s]{3,}`)
	runs := runRegex.FindAllString(text, -1)
	return strings.Join(runs, " ")
}

func topMatches(text string, dictionary []string) []string {
	out := make([]string, 0, 12)
	seen := map[string]struct{}{}
	for _, token := range dictionary {
		pattern := `\b` + regexp.QuoteMeta(strings.ToLower(token)) + `\b`
		if regexp.MustCompile(pattern).FindStringIndex(text) == nil {
			continue
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		out = append(out, token)
	}
	return out
}

func keywordDictionary() []string {
	return []string{
		"software engineer", "backend", "frontend", "full stack", "golang", "react",
		"mechanical engineer", "electrical engineer", "manufacturing", "quality engineer",
		"data engineer", "product manager", "devops", "qa", "automation", "cad",
	}
}

func roleFamilyDictionary() []string {
	return []string{
		"software", "mechanical", "electrical", "product", "operations", "manufacturing",
	}
}

func locationDictionary() []string {
	return []string{
		"remote", "united states", "austin", "minneapolis", "new york", "san francisco",
	}
}
