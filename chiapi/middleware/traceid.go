package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
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

		r = r.WithContext(WithTraceID(r.Context(), traceID))
		next.ServeHTTP(w, r)
	})
}

var (
	traceIDKey = contextKey("traceID")
)

// WithTraceID adds the traceID to the request context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID retrieves the traceID from the request context.
func GetTraceID(ctx context.Context) string {
	if val, ok := ctx.Value(traceIDKey).(string); ok {
		return val
	}

	return ""
}
