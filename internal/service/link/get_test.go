package link

import (
	"context"
	"errors"
	"testing"

	"github.com/huypham67/bookmark-service/internal/repository/link/mocks"
	utilsMocks "github.com/huypham67/bookmark-service/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_GetOriginalURL(t *testing.T) {
	t.Parallel()

	type args struct {
		code string
	}

	testCases := []struct {
		name           string
		args           args
		setupMocks     func(context.Context, *mocks.Repository)
		verifyResponse func(*testing.T, string, error)
	}{
		{
			name: "should get original URL successfully",
			args: args{
				code: "abc1234",
			},
			setupMocks: func(ctx context.Context, mockRepo *mocks.Repository) {
				mockRepo.
					On("GetLink", ctx, "abc1234").
					Return("https://google.com", nil).
					Once()
			},
			verifyResponse: func(t *testing.T, url string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://google.com", url)
			},
		},
		{
			name: "should return error when code does not exist",
			args: args{
				code: "missing",
			},
			setupMocks: func(ctx context.Context, mockRepo *mocks.Repository) {
				mockRepo.
					On("GetLink", ctx, "missing").
					Return("", errors.New("code not found")).
					Once()
			},
			verifyResponse: func(t *testing.T, url string, err error) {
				assert.Error(t, err)
				assert.Empty(t, url)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockRepo := new(mocks.Repository)
			mockCodeGen := new(utilsMocks.CodeGenerator)
			tc.setupMocks(ctx, mockRepo)

			service := NewService(mockRepo, mockCodeGen)

			url, err := service.GetOriginalURL(ctx, tc.args.code)

			tc.verifyResponse(t, url, err)

			mockRepo.AssertExpectations(t)
			mockCodeGen.AssertExpectations(t)
		})
	}
}
