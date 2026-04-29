package notifications

import (
	"encoding/json"
	httpx "job-crawler/apps/api/internal/platform/http"
	"net/http"
)

type Handler struct{ Service Service }

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		uid := httpx.UserIDFromContext(r.Context())
		result, err := h.Service.Get(uid)
		if err != nil {
			httpx.WriteError(w, 500, "failed to load settings")
			return
		}
		httpx.WriteJSON(w, 200, result)
		return
	}
	if r.Method != http.MethodPatch {
		httpx.WriteError(w, 405, "method not allowed")
		return
	}
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpx.WriteError(w, 400, "invalid payload")
		return
	}
	uid := httpx.UserIDFromContext(r.Context())
	result, err := h.Service.Update(uid, payload)
	if err != nil {
		httpx.WriteError(w, 500, "failed to update settings")
		return
	}
	httpx.WriteJSON(w, 200, result)
}
