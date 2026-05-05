package companies

import (
	httpx "job-crawler/apps/api/internal/platform/http"
	"net/http"
	"strings"
)

type Handler struct{ Service Service }

func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(q) < 2 {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"companies": []map[string]any{}})
		return
	}
	companies, err := h.Service.Search(q)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to search companies")
		return
	}
	results := make([]map[string]any, 0, len(companies))
	for _, company := range companies {
		results = append(results, map[string]any{
			"name":       company,
			"isVerified": true,
		})
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"companies": results})
}
