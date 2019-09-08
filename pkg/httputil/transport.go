package httputil

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/rickbassham/example-go/chiapi/middleware"
)

// HeaderTransport is an http.RoundTripper that will add a header to all requests.
func HeaderTransport(key, value string, old http.RoundTripper) http.RoundTripper {
	if old == nil {
		old = http.DefaultTransport
	}

	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.Header.Set(key, value)

		return old.RoundTrip(req)
	})
}

// TraceIDTransport ensures a trace id header is set on all outgoing requests.
func TraceIDTransport(old http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		ctx := req.Context()

		// use the existing trace id if it exists
		traceID := middleware.GetTraceID(ctx)
		if traceID == "" {
			traceID = uuid.New().String()
			req = req.WithContext(middleware.WithTraceID(ctx, traceID))
		}

		req.Header.Set("X-Trace-Id", traceID)

		return old.RoundTrip(req)
	})
}

// DefaultLogTransport will add the given logger to all request contexts.
func DefaultLogTransport(log *zap.Logger, old http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req = req.WithContext(middleware.WithLogger(req.Context(), log))

		return old.RoundTrip(req)
	})
}

// LogTransport will log every outgoing request.
func LogTransport(old http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		l := middleware.GetLogger(req.Context()).With(
			zap.String("direction", "outgoing"),
			zap.String("url", req.URL.String()),
			zap.String("method", req.Method),
		)

		start := time.Now()

		resp, err := old.RoundTrip(req)

		l = l.With(
			zap.Duration("duration", time.Since(start)),
		)

		if err != nil {
			l.Error("request error", zap.Error(err))
		}

		if resp != nil {
			status := resp.StatusCode

			l = l.With(
				zap.Int("status", status),
			)

			if status >= 300 && status < 400 {
				l = l.With(zap.String("location", resp.Header.Get("Location")))
			}

			if status < 400 {
				l.Info("request complete")
			} else if status < 500 {
				l.Warn("request complete")
			} else {
				l.Error("request complete")
			}
		}

		return resp, err
	})
}

// APIKeyTransport is an http.RoundTripper that will add basic authentication
// headers to all requests.
func APIKeyTransport(apiKey string, old http.RoundTripper) http.RoundTripper {
	if old == nil {
		old = http.DefaultTransport
	}

	return HeaderTransport("X-Api-Key", apiKey, old)
}

// BasicAuthTransport is an http.RoundTripper that will add basic authentication
// headers to all requests.
func BasicAuthTransport(username, password string, old http.RoundTripper) http.RoundTripper {
	if old == nil {
		old = http.DefaultTransport
	}

	hashed := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	headerVal := fmt.Sprintf("Basic %s", hashed)

	return HeaderTransport("Authorization", headerVal, old)
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
