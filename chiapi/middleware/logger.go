package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"github.com/rickbassham/example-go/pkg/logging"
	"github.com/rickbassham/example-go/pkg/tracing"
)

// Logger middleware logs each request and adds the logger to the request
// context for use in request handlers.
func Logger(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			traceID := tracing.FromContext(r.Context())

			l := log.With(
				zap.String("direction", "incoming"),
				zap.String("trace_id", traceID),
				zap.String("method", r.Method),
				zap.String("path", path),
				zap.String("query", r.URL.RawQuery),
				zap.String("referer", r.Referer()),
				zap.String("user_agent", r.UserAgent()))

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(logging.WithLogger(r.Context(), l))
			next.ServeHTTP(ww, r)

			rctx := chi.RouteContext(r.Context())

			routePattern := strings.Join(rctx.RoutePatterns, "")
			routePattern = strings.Replace(routePattern, "/*", "", -1)

			l = l.With(zap.String("route_pattern", routePattern))

			status := ww.Status()
			if status >= 300 && status < 400 {
				// if we are doing a redirect, log where we are going
				l = l.With(zap.String("location", ww.Header().Get("Location")))
			}

			l = l.With(
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
				zap.Duration("duration", time.Since(start)),
			)

			if status < 400 {
				l.Info("request complete")
			} else if status < 500 {
				l.Warn("request complete")
			} else {
				l.Error("request complete")
			}
		}

		return http.HandlerFunc(fn)
	}
}
