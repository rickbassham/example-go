package identity

import "context"

type contextKey string

func (k contextKey) String() string {
	return "context key: " + string(k)
}

var (
	userKey = contextKey("user")
)

// WithUser adds the user to the request context.
func WithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// FromContext retrieves the user from the context.
func FromContext(ctx context.Context) string {
	if val, ok := ctx.Value(userKey).(string); ok {
		return val
	}

	return ""
}
