package ping

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisPinger implements the Pinger interface for Redis.
type RedisPinger struct {
	client *redis.Client
}

// NewRedis creates a new RedisPinger with the provided Redis client.
func NewRedis(client *redis.Client) Pinger {
	return &RedisPinger{
		client: client,
	}
}

// Ping checks the health of the Redis connection by sending a PING command and returning any error encountered.
func (p *RedisPinger) Ping(ctx context.Context) error {
	return p.client.Ping(ctx).Err()
}
