package worker

import (
    "net/http"

    httpx "job-crawler/apps/api/internal/platform/http"
)

type Handler struct {
    Service Service
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

    job := h.Service.Tick()
    httpx.WriteJSON(w, http.StatusOK, map[string]any{
        "status": "ok",
        "generatedJobId": job.ID,
    })
}
