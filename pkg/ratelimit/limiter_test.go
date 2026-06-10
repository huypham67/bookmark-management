package ratelimit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/huypham67/bookmark-service-monolithic/pkg/ratelimit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFixedWindowLimiter_Allow(t *testing.T) {
	t.Parallel()

	const (
		key    = "ratelimit:user:123"
		limit  = 20
		window = 10 * time.Second
	)

	testCases := []struct {
		name        string
		setupMock   func(*mocks.Store)
		wantAllowed bool
		wantErr     bool
	}{
		{
			name: "should allow and increment when under limit",
			setupMock: func(s *mocks.Store) {
				s.On("GetCounter", mock.Anything, key).Return(limit-1, nil).Once()
				s.On("IncrementCounter", mock.Anything, key, window).Once()
			},
			wantAllowed: true,
		},
		{
			name: "should deny without incrementing when at limit",
			setupMock: func(s *mocks.Store) {
				s.On("GetCounter", mock.Anything, key).Return(limit, nil).Once()
			},
			wantAllowed: false,
		},
		{
			name: "should return error when store fails",
			setupMock: func(s *mocks.Store) {
				s.On("GetCounter", mock.Anything, key).Return(0, errors.New("redis down")).Once()
			},
			wantAllowed: false,
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := mocks.NewStore(t)
			tc.setupMock(store)

			limiter := &fixedWindowLimiter{store: store, limit: limit, window: window}

			allowed, err := limiter.Allow(context.Background(), key)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantAllowed, allowed)
		})
	}
}
