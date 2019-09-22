package cache

import (
	"context"

	"github.com/go-redis/redis"
)

// Client represents the functions needed for this wrapper.
type Client interface {
	AddHook(redis.Hook)
	WithContext(context.Context) *redis.Client
}

// Cache is a logged and instrumented wrapper around a redis client.
type Cache struct {
	client Client
}

// New creates a new Cache.
func New(client Client) *Cache {
	client.AddHook(LoggerHook{})
	client.AddHook(NewRelicHook{})

	return &Cache{
		client: client,
	}
}

// GetValue retrieves the value from redis.
func (c *Cache) GetValue(ctx context.Context) (string, error) {
	return c.client.WithContext(ctx).Get("value").Result()
}
