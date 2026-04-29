package httpx

import (
    "net/http"
    "strings"

    authx "job-crawler/apps/api/internal/platform/auth"
)

type ctxKey string
const UserIDKey ctxKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/healthz" || r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/signup" {
            next.ServeHTTP(w, r)
            return
        }

        authHeader := r.Header.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            WriteError(w, http.StatusUnauthorized, "missing bearer token")
            return
        }

        token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
        userID, err := authx.ParseToken(token)
        if err != nil || userID == "" {
            WriteError(w, http.StatusUnauthorized, "invalid bearer token")
            return
        }

        ctx := withUserID(r.Context(), userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
