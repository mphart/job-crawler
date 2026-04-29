package httpx

import "net/http"

type ctxKey string
const UserIDKey ctxKey = "user_id"

func AuthPassthrough(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        uid := r.Header.Get("X-User-Id")
        if uid == "" { uid = "u_1" }
        ctx := r.Context()
        ctx = withUserID(ctx, uid)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
