package ping

import (
	"testing"

	"github.com/huypham67/bookmark-service/pkg/redis"
)

func newTestPinger(t *testing.T) (Pinger, *redis.MockRedis) {
	t.Helper()

	mockRedis := redis.NewMockRedis(t)
	pinger := NewRedis(mockRedis.Client)

	return pinger, mockRedis
}
