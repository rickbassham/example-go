package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/rickbassham/example-go/pkg/tracing"
)

// TraceID tries to read the X-Trace-Id request header, and if empty, creates
// a new traceID to use. This traceID will be added to the response X-Trace-Id
// header and to the request context.
func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var traceID string

		if r.Header != nil {
			traceID = r.Header.Get("X-Trace-Id")

			if traceID == "" {
				traceID = uuid.New().String()
			}
		}

		w.Header().Set("X-Trace-Id", traceID)

		r = r.WithContext(tracing.WithTraceID(r.Context(), traceID))
		next.ServeHTTP(w, r)
	})
}
