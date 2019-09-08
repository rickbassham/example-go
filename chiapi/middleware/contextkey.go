package middleware

type contextKey string

func (k contextKey) String() string {
	return "middleware context key: " + string(k)
}
