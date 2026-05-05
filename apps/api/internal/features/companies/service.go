package companies

import (
	"job-crawler/apps/api/internal/platform/db"
	"sort"
	"strings"
)

type CompanyStore interface {
	SearchCompanies(q string) ([]string, error)
}

type InMemoryStore struct{ Inner *db.Store }

func (s InMemoryStore) SearchCompanies(q string) ([]string, error) {
	return s.Inner.SearchCompanies(q), nil
}

type MySQLStore struct{ Inner *db.MySQLAuthStore }

func (s MySQLStore) SearchCompanies(q string) ([]string, error) {
	return s.Inner.SearchCompanies(q)
}

type Service struct{ Store CompanyStore }

func (s Service) Search(q string) ([]string, error) {
	found, err := s.Store.SearchCompanies(q)
	if err != nil {
		return nil, err
	}
	return mergeWithCatalog(q, found), nil
}

func mergeWithCatalog(q string, fromStore []string) []string {
	needle := strings.ToLower(strings.TrimSpace(q))
	seen := map[string]struct{}{}
	out := make([]string, 0, len(fromStore)+12)

	appendCompany := func(company string) {
		company = strings.TrimSpace(company)
		if company == "" {
			return
		}
		key := strings.ToLower(company)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		out = append(out, company)
	}

	for _, company := range fromStore {
		appendCompany(company)
	}
	for _, company := range verifiedCompanyCatalog {
		if needle == "" || strings.Contains(strings.ToLower(company), needle) {
			appendCompany(company)
		}
	}

	sort.Strings(out)
	if len(out) > 25 {
		return out[:25]
	}
	return out
}

var verifiedCompanyCatalog = []string{
	"Amazon",
	"Apple",
	"Cloudflare",
	"Cisco",
	"CrowdStrike",
	"Capital One",
	"Datadog",
	"DoorDash",
	"Data Bricks",
	"Google",
	"HubSpot",
	"LinkedIn",
	"Meta",
	"Microsoft",
	"Netflix",
	"NVIDIA",
	"OpenAI",
	"Palantir",
	"Salesforce",
	"Shopify",
	"Snowflake",
	"Spotify",
	"Square",
	"Stripe",
	"ServiceNow",
	"American Express",
	"Uber",
}
