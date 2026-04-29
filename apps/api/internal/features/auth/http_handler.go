package auth

import (
    "encoding/json"
    "net/http"

    httpx "job-crawler/apps/api/internal/platform/http"
)

type Handler struct{ Service Service }

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { httpx.WriteError(w, 405, "method not allowed"); return }
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { httpx.WriteError(w, 400, "invalid payload"); return }
    out, err := h.Service.Login(req.Email, req.Password)
    if err != nil { httpx.WriteError(w, 401, err.Error()); return }
    httpx.WriteJSON(w, 200, out)
}

func (h Handler) Signup(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { httpx.WriteError(w, 405, "method not allowed"); return }
    var req SignupRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { httpx.WriteError(w, 400, "invalid payload"); return }
    out, err := h.Service.Signup(req.Email, req.Username, req.Password, req.Keywords)
    if err != nil { httpx.WriteError(w, 400, err.Error()); return }
    httpx.WriteJSON(w, 200, out)
}
