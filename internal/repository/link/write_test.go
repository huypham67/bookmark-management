package link

import (
	"context"
	"testing"
	"time"

	"github.com/huypham67/bookmark-service/pkg/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_SaveLink(t *testing.T) {
	t.Parallel()

	type args struct {
		code string
		url  string
		exp  int64
	}

	testCases := []struct {
		name   string
		args   args
		verify func(*testing.T, context.Context, Repository, *redis.MockRedis, args)
	}{
		{
			name: "should save link successfully",
			args: args{
				code: "abc1234",
				url:  "https://www.google.com",
				exp:  1234,
			},
			verify: func(t *testing.T, ctx context.Context, repo Repository, mockRedis *redis.MockRedis, a args) {
				url, err := repo.GetLink(ctx, a.code)
				require.NoError(t, err)
				require.Equal(t, a.url, url)
			},
		},
		{
			name: "should overwrite existing link successfully",
			args: args{
				code: "abc1234",
				url:  "https://www.google.com/v2",
				exp:  1234,
			},
			verify: func(t *testing.T, ctx context.Context, repo Repository, mockRedis *redis.MockRedis, a args) {
				err := repo.SaveLink(ctx, a.code, "https://www.google.com/v1", 1234)
				require.NoError(t, err)
				err = repo.SaveLink(ctx, a.code, a.url, a.exp)
				require.NoError(t, err)
				url, err := repo.GetLink(ctx, a.code)
				require.NoError(t, err)
				require.Equal(t, a.url, url)
			},
		},
		{
			name: "should expire link after TTL exceeded",
			args: args{
				code: "abc1234",
				url:  "https://www.google.com",
				exp:  1,
			},
			verify: func(t *testing.T, ctx context.Context, repo Repository, mockRedis *redis.MockRedis, a args) {
				exists, err := repo.CheckExists(ctx, a.code)
				require.NoError(t, err)
				assert.True(t, exists)

				mockRedis.Server.FastForward(2 * time.Second)

				exists, err = repo.CheckExists(ctx, a.code)

				require.NoError(t, err)
				assert.False(t, exists)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			repo, mockRedis := newTestRepository(t)

			err := repo.SaveLink(ctx, tc.args.code, tc.args.url, tc.args.exp)

			require.NoError(t, err)

			tc.verify(t, ctx, repo, mockRedis, tc.args)
		})
	}
}

