package link

import (
	"context"
	"errors"
	"testing"

	linkDTO "github.com/huypham67/bookmark-service/internal/dto/link"
	"github.com/huypham67/bookmark-service/internal/repository/link/mocks"
	utilsMocks "github.com/huypham67/bookmark-service/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_ShortenURL(t *testing.T) {
	t.Parallel()

	type args struct {
		request linkDTO.ShortenURLRequest
	}

	testCases := []struct {
		name           string
		args           args
		setupMocks     func(context.Context, *mocks.Repository, *utilsMocks.CodeGenerator)
		verifyResponse func(*testing.T, string, error)
	}{
		{
			name: "should shorten URL successfully when code does not exist",
			args: args{
				request: linkDTO.ShortenURLRequest{
					Url: "https://google.com",
					Exp: 3600,
				},
			},
			setupMocks: func(ctx context.Context, mockRepo *mocks.Repository, mockCodeGen *utilsMocks.CodeGenerator) {
				mockCodeGen.
					On("Generate", shortCodeLength).
					Return("abc1234", nil).
					Once()

				mockRepo.
					On("CheckExists", ctx, "abc1234").
					Return(false, nil).
					Once()

				mockRepo.
					On(
						"SaveLink",
						ctx,
						"abc1234",
						"https://google.com",
						int64(3600),
					).
					Return(nil).
					Once()
			},
			verifyResponse: func(t *testing.T, code string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "abc1234", code)
			},
		},
		{
			name: "should return error when code generation fails",
			args: args{
				request: linkDTO.ShortenURLRequest{
					Url: "https://google.com",
					Exp: 3600,
				},
			},
			setupMocks: func(ctx context.Context, mockRepo *mocks.Repository, mockCodeGen *utilsMocks.CodeGenerator) {
				mockCodeGen.
					On("Generate", shortCodeLength).
					Return("", errors.New("code generation failed")).
					Once()
			},
			verifyResponse: func(t *testing.T, code string, err error) {
				assert.Error(t, err)
				assert.Empty(t, code)
			},
		}, {
			name: "should return error when checking code existence fails",
			args: args{
				request: linkDTO.ShortenURLRequest{
					Url: "https://google.com",
					Exp: 3600,
				},
			},
			setupMocks: func(ctx context.Context, mockRepo *mocks.Repository, mockCodeGen *utilsMocks.CodeGenerator) {
				mockCodeGen.
					On("Generate", shortCodeLength).
					Return("abc1234", nil).
					Once()

				mockRepo.
					On("CheckExists", ctx, "abc1234").
					Return(false, errors.New("redis error")).
					Once()
			},
			verifyResponse: func(t *testing.T, code string, err error) {
				assert.Error(t, err)
				assert.Empty(t, code)
			},
		}, {
			name: "should retry code generation when code already exists",
			args: args{
				request: linkDTO.ShortenURLRequest{
					Url: "https://google.com",
					Exp: 3600,
				},
			},
			setupMocks: func(ctx context.Context, mockRepo *mocks.Repository, mockCodeGen *utilsMocks.CodeGenerator) {
				// First attempt
				mockCodeGen.
					On("Generate", shortCodeLength).
					Return("abc1234", nil).
					Once()

				mockRepo.
					On("CheckExists", ctx, "abc1234").
					Return(true, nil).
					Once()

				// Second attempt
				mockCodeGen.
					On("Generate", shortCodeLength).
					Return("def5678", nil).
					Once()

				mockRepo.
					On("CheckExists", ctx, "def5678").
					Return(false, nil).
					Once()

				mockRepo.
					On(
						"SaveLink",
						ctx,
						"def5678",
						"https://google.com",
						int64(3600),
					).
					Return(nil).
					Once()
			},
			verifyResponse: func(t *testing.T, code string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "def5678", code)
			},
		}, {
			name: "should return error when saving link fails",
			args: args{
				request: linkDTO.ShortenURLRequest{
					Url: "https://google.com",
					Exp: 3600,
				},
			},
			setupMocks: func(ctx context.Context, mockRepo *mocks.Repository, mockCodeGen *utilsMocks.CodeGenerator) {
				mockCodeGen.
					On("Generate", shortCodeLength).
					Return("abc1234", nil).
					Once()

				mockRepo.
					On("CheckExists", ctx, "abc1234").
					Return(false, nil).
					Once()

				mockRepo.
					On(
						"SaveLink",
						ctx,
						"abc1234",
						"https://google.com",
						int64(3600),
					).
					Return(errors.New("save error")).
					Once()
			},
			verifyResponse: func(t *testing.T, code string, err error) {
				assert.Error(t, err)
				assert.Empty(t, code)
			},
		},
		{
			name: "should return error when context is cancelled",
			args: args{
				request: linkDTO.ShortenURLRequest{
					Url: "https://google.com",
					Exp: 3600,
				},
			},
			setupMocks: func(ctx context.Context, mockRepo *mocks.Repository, mockCodeGen *utilsMocks.CodeGenerator) {
				mockCodeGen.
					On("Generate", shortCodeLength).
					Return("abc1234", nil).
					Once()

				mockRepo.
					On("CheckExists", ctx, "abc1234").
					Return(false, context.Canceled).
					Once()
			},
			verifyResponse: func(t *testing.T, code string, err error) {
				assert.Error(t, err)
				assert.Empty(t, code)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var ctx context.Context
			mockRepo := new(mocks.Repository)
			mockCodeGen := new(utilsMocks.CodeGenerator)

			// For context cancellation test, create a cancelled context
			if tc.name == "should return error when context is cancelled" {
				cancelledCtx, cancel := context.WithCancel(context.Background())
				cancel()
				ctx = cancelledCtx
			} else {
				ctx = context.Background()
			}

			tc.setupMocks(ctx, mockRepo, mockCodeGen)

			service := NewService(mockRepo, mockCodeGen)

			code, err := service.ShortenURL(ctx, tc.args.request)

			tc.verifyResponse(t, code, err)

			mockRepo.AssertExpectations(t)
			mockCodeGen.AssertExpectations(t)
		})
	}
}
