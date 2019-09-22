package cache

type contextKey string

func (k contextKey) String() string {
	return "cache context key: " + string(k)
}
