package api

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const tokenKey ctxKey = "token"

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Authorization header required")
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			writeError(w, http.StatusUnauthorized, "unauthorized", "Authorization token must not be empty")
			return
		}
		ctx := context.WithValue(r.Context(), tokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func tokenFromCtx(ctx context.Context) string {
	return ctx.Value(tokenKey).(string)
}
