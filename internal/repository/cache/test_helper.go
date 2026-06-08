package cache

import (
	"testing"

	"github.com/huypham67/bookmark-service/pkg/redis"
)

func newTestRepository(t *testing.T) (Repository, *redis.MockRedis) {
	t.Helper()

	mockRedis := redis.NewMockRedis(t)
	repo := NewRedis(mockRedis.Client)

	return repo, mockRedis
}
