package middleware

import (
	"net/http"
)

// Version middleware adds the X-Version header to the response.
func Version(v string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Version", v)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
