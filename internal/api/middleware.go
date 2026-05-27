package api

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const tokenKey ctxKey = "token"

func tokenMiddleware(allowedTokens []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedTokens))
	for _, t := range allowedTokens {
		allowed[t] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
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
			if len(allowed) > 0 {
				if _, ok := allowed[token]; !ok {
					writeError(w, http.StatusUnauthorized, "unauthorized", "token not authorized")
					return
				}
			}
			ctx := context.WithValue(r.Context(), tokenKey, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func tokenFromCtx(ctx context.Context) string {
	return ctx.Value(tokenKey).(string)
}
