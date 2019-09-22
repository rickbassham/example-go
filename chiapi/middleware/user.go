package middleware

import (
	"net/http"

	"github.com/go-chi/jwtauth"

	"github.com/rickbassham/example-go/pkg/identity"
)

// User get's user email from JWT token and adds it to the request context.
func User(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		user, ok := claims["email"].(string)

		if ok && user != "" {
			r = r.WithContext(identity.WithUser(r.Context(), user))
		}

		next.ServeHTTP(w, r)
	})
}
