package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// Logger middleware logs each request and adds the logger to the request
// context for use in request handlers.
func Logger(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			traceID := GetTraceID(r.Context())

			l := log.With(
				zap.String("trace_id", traceID),
				zap.String("method", r.Method),
				zap.String("path", path),
				zap.String("query", r.URL.RawQuery),
				zap.String("referer", r.Referer()),
				zap.String("user_agent", r.UserAgent()))

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(WithLogger(r.Context(), l))
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

var (
	loggerKey = contextKey("logger")
)

// WithLogger adds the logger to the request context.
func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// GetLogger retrieves the logger from the context. If there is no logger on the context, it will
// return a valid NoOp logger.
func GetLogger(ctx context.Context) *zap.Logger {
	if val, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return val
	}

	return zap.NewNop()
}
