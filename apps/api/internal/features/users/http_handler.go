package users

import (
	httpx "job-crawler/apps/api/internal/platform/http"
	"net/http"
)

type Handler struct{ Service Service }

func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, 405, "method not allowed")
		return
	}
	result, err := h.Service.Search(r.URL.Query().Get("q"))
	if err != nil {
		httpx.WriteError(w, 500, "search failed")
		return
	}
	httpx.WriteJSON(w, 200, result)
}
