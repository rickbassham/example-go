package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
)

// Authenticator is used to verify a valid token is on the request. The param unauthorized should be
// a request handler to render your 401 response. If it is nil, a simple default response will be
// written.
func Authenticator(unauthorized http.Handler) func(next http.Handler) http.Handler {
	if unauthorized == nil {
		unauthorized = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "unauthorized")
		})
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())
			if err != nil {
				unauthorized.ServeHTTP(w, r)
				return
			}

			if token == nil || !token.Valid {
				unauthorized.ServeHTTP(w, r)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
