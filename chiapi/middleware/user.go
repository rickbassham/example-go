package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth"
)

// User get's user email from JWT token and adds it to the request context.
func User(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		user, ok := claims["email"].(string)

		if ok && user != "" {
			r = r.WithContext(WithUser(r.Context(), user))
		}

		next.ServeHTTP(w, r)
	})
}

var (
	userKey = contextKey("user")
)

// WithUser adds the user to the request context.
func WithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetUser retrieves the user from the context.
func GetUser(ctx context.Context) string {
	if val, ok := ctx.Value(userKey).(string); ok {
		return val
	}

	return ""
}
