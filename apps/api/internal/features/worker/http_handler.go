package worker

import (
	"encoding/json"
	"net/http"

	"job-crawler/apps/api/internal/platform/db"
	httpx "job-crawler/apps/api/internal/platform/http"
)

type Handler struct {
	Service     Service
	WorkerToken string
}

func (h Handler) Tick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if h.WorkerToken != "" && r.Header.Get("X-Worker-Token") != h.WorkerToken {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid worker token")
		return
	}

	var jobs []db.ScrapedJob
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&jobs)
	}

	if len(jobs) == 0 {
		httpx.WriteError(w, http.StatusBadRequest, "no scraped jobs provided")
		return
	}
	inserted := h.Service.Ingest(jobs)
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":   "ok",
		"inserted": inserted,
	})
}

func (h Handler) Preferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if h.WorkerToken != "" && r.Header.Get("X-Worker-Token") != h.WorkerToken {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid worker token")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"keywords": h.Service.Keywords(),
	})
}

func (h Handler) DigestCandidates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if h.WorkerToken != "" && r.Header.Get("X-Worker-Token") != h.WorkerToken {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid worker token")
		return
	}
	candidates, err := h.Service.DigestCandidates()
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get candidates")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"candidates": candidates})
}

func (h Handler) NotificationSent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if h.WorkerToken != "" && r.Header.Get("X-Worker-Token") != h.WorkerToken {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid worker token")
		return
	}
	var payload struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.UserID == "" {
		httpx.WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := h.Service.MarkNotificationSent(payload.UserID); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to mark sent")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}
