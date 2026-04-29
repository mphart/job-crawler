package httpx

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string

func withUserID(ctx context.Context, userID string) context.Context { return context.WithValue(ctx, UserIDKey, userID) }
func UserIDFromContext(ctx context.Context) string {
    if v, ok := ctx.Value(UserIDKey).(string); ok { return v }
    return ""
}

type Handlers struct {
    AuthLogin            http.HandlerFunc
    AuthSignup           http.HandlerFunc
    FeedList             http.HandlerFunc
    FeedAction           http.HandlerFunc
    ProfileUpdateMe      http.HandlerFunc
    ProfileGetByID       http.HandlerFunc
    NotificationSettings http.HandlerFunc
    UsersSearch          http.HandlerFunc
}

func NewRouter(h Handlers) http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { WriteJSON(w, 200, map[string]string{"status":"ok"}) })
    mux.HandleFunc("/api/auth/login", h.AuthLogin)
    mux.HandleFunc("/api/auth/signup", h.AuthSignup)
    mux.HandleFunc("/api/feed", h.FeedList)
    mux.HandleFunc("/api/feed/", func(w http.ResponseWriter, r *http.Request) {
        if strings.HasSuffix(r.URL.Path, "/apply") || strings.HasSuffix(r.URL.Path, "/reject") {
            h.FeedAction(w, r)
            return
        }
        WriteError(w, 404, "not found")
    })
    mux.HandleFunc("/api/profiles/me", h.ProfileUpdateMe)
    mux.HandleFunc("/api/profiles/", h.ProfileGetByID)
    mux.HandleFunc("/api/notifications/settings", h.NotificationSettings)
    mux.HandleFunc("/api/users/search", h.UsersSearch)
    return CorsMiddleware(AuthMiddleware(mux))
}
