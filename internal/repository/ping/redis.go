package ping

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisPinger struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) Pinger {
	return &RedisPinger{
		client: client,
	}
}

func (p *RedisPinger) Ping(ctx context.Context) error {
	return p.client.Ping(ctx).Err()
}
