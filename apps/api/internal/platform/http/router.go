package httpx

import (
    "context"
    "net/http"
    "strings"

    "job-crawler/apps/api/internal/features/auth"
    "job-crawler/apps/api/internal/features/feed"
    "job-crawler/apps/api/internal/features/notifications"
    "job-crawler/apps/api/internal/features/profiles"
    "job-crawler/apps/api/internal/features/users"
)

type contextKey string

func withUserID(ctx context.Context, userID string) context.Context { return context.WithValue(ctx, UserIDKey, userID) }
func UserIDFromContext(ctx context.Context) string {
    if v, ok := ctx.Value(UserIDKey).(string); ok { return v }
    return "u_1"
}

type Handlers struct {
    Auth *auth.Handler
    Feed *feed.Handler
    Profiles *profiles.Handler
    Notifications *notifications.Handler
    Users *users.Handler
}

func NewRouter(h Handlers) http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { WriteJSON(w, 200, map[string]string{"status":"ok"}) })
    mux.HandleFunc("/api/auth/login", h.Auth.Login)
    mux.HandleFunc("/api/auth/signup", h.Auth.Signup)
    mux.HandleFunc("/api/feed", h.Feed.List)
    mux.HandleFunc("/api/feed/", func(w http.ResponseWriter, r *http.Request) {
        if strings.HasSuffix(r.URL.Path, "/apply") { h.Feed.Apply(w, r); return }
        if strings.HasSuffix(r.URL.Path, "/reject") { h.Feed.Reject(w, r); return }
        WriteError(w, 404, "not found")
    })
    mux.HandleFunc("/api/profiles/me", h.Profiles.UpdateMe)
    mux.HandleFunc("/api/profiles/", h.Profiles.GetByID)
    mux.HandleFunc("/api/notifications/settings", h.Notifications.Update)
    mux.HandleFunc("/api/users/search", h.Users.Search)
    return AuthPassthrough(mux)
}
