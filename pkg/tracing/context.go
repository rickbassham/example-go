package tracing

import "context"

type contextKey string

func (k contextKey) String() string {
	return "context key: " + string(k)
}

var (
	traceIDKey = contextKey("traceID")
)

// WithTraceID adds the traceID to the request context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// FromContext retrieves the traceID from the request context.
func FromContext(ctx context.Context) string {
	if val, ok := ctx.Value(traceIDKey).(string); ok {
		return val
	}

	return ""
}
