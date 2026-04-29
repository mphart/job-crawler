package feed

import (
	httpx "job-crawler/apps/api/internal/platform/http"
	"net/http"
	"strings"
)

type Handler struct{ Service Service }

func jobID(path, suffix string) string {
	return strings.TrimSuffix(strings.TrimPrefix(path, "/api/feed/"), suffix)
}
func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, 405, "method not allowed")
		return
	}
	uid := httpx.UserIDFromContext(r.Context())
	jobs, err := h.Service.List(uid, r.URL.Query().Get("search"), r.URL.Query().Get("sortBy"))
	if err != nil {
		httpx.WriteError(w, 500, "failed to list feed")
		return
	}
	httpx.WriteJSON(w, 200, jobs)
}
func (h Handler) Action(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/apply") {
		h.Apply(w, r)
		return
	}
	if strings.HasSuffix(r.URL.Path, "/reject") {
		h.Reject(w, r)
		return
	}
	httpx.WriteError(w, 404, "not found")
}
func (h Handler) Apply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, 405, "method not allowed")
		return
	}
	uid := httpx.UserIDFromContext(r.Context())
	if err := h.Service.Apply(uid, jobID(r.URL.Path, "/apply")); err != nil {
		httpx.WriteError(w, 404, err.Error())
		return
	}
	w.WriteHeader(204)
}
func (h Handler) Reject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, 405, "method not allowed")
		return
	}
	uid := httpx.UserIDFromContext(r.Context())
	if err := h.Service.Reject(uid, jobID(r.URL.Path, "/reject")); err != nil {
		httpx.WriteError(w, 404, err.Error())
		return
	}
	w.WriteHeader(204)
}
