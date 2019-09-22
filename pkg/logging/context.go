package logging

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

func (k contextKey) String() string {
	return "context key: " + string(k)
}

var (
	loggerKey = contextKey("logger")
)

// WithLogger adds the logger to the request context.
func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext retrieves the logger from the context. If there is no logger on the context, it will
// return a valid no-op logger.
func FromContext(ctx context.Context) *zap.Logger {
	if val, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return val
	}

	return zap.NewNop()
}
