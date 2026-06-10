package ratelimitprovider

import (
	"fmt"

	"github.com/huypham67/bookmark-service-monolithic/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
)

// NewLimiter builds a rate limit Limiter over the shared Redis client, loading its policy from environment variables with the given prefix.
func NewLimiter(client *redis.Client, envPrefix string) (ratelimit.Limiter, error) {
	cfg, err := LoadConfig(envPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to load rate limit config: %w", err)
	}

	store := ratelimit.NewRedis(client)
	return ratelimit.NewFixedWindow(store, cfg.Limit, cfg.Window), nil
}
