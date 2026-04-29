package httpx

import (
    "encoding/json"
    "net/http"
)

type ErrorResponse struct { Error string `json:"error"` }

func WriteJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
    WriteJSON(w, status, ErrorResponse{Error: msg})
}
