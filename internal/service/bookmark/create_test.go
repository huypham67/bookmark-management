package bookmark

import (
	"context"
	"testing"

	bookmarkDTO "github.com/huypham67/bookmark-service/internal/dto/bookmark"
	"github.com/huypham67/bookmark-service/internal/model"
	bookmarkMocks "github.com/huypham67/bookmark-service/internal/repository/bookmark/mocks"
	"github.com/huypham67/bookmark-service/pkg/dbutils"
	utilsMocks "github.com/huypham67/bookmark-service/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func matchBookmark(expected *model.Bookmark) interface{} {
	return mock.MatchedBy(func(actual *model.Bookmark) bool {
		return actual.Description == expected.Description &&
			actual.URL == expected.URL &&
			actual.Code == expected.Code &&
			actual.UserID == expected.UserID
	})
}

func TestService_Create(t *testing.T) {
	t.Parallel()

	type args struct {
		userID  string
		request bookmarkDTO.CreateBookmarkRequest
	}

	testCases := []struct {
		name           string
		args           args
		setupMocks     func(context.Context, *bookmarkMocks.Repository, *utilsMocks.CodeGenerator)
		verifyResponse func(*testing.T, *model.Bookmark, error)
	}{
		{
			name: "should create bookmark successfully",
			args: args{
				userID: "user-id-123",
				request: bookmarkDTO.CreateBookmarkRequest{
					Description: "Test Bookmark",
					URL:         "https://example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				bookmarkRepo *bookmarkMocks.Repository,
				codeGenerator *utilsMocks.CodeGenerator,
			) {
				codeGenerator.
					On("Generate", bookmarkCodeLength).
					Return("abc123", nil).
					Once()

				expectedBookmark := &model.Bookmark{
					Description: "Test Bookmark",
					URL:         "https://example.com",
					Code:        "abc123",
					UserID:      "user-id-123",
				}

				bookmarkRepo.
					On(
						"Create",
						ctx,
						matchBookmark(expectedBookmark),
					).
					Return(nil).
					Once()
			},
			verifyResponse: func(t *testing.T, bm *model.Bookmark, err error) {
				assert.NoError(t, err)
				require.NotNil(t, bm)

				assert.Equal(t, "Test Bookmark", bm.Description)
				assert.Equal(t, "https://example.com", bm.URL)
				assert.Equal(t, "abc123", bm.Code)
				assert.Equal(t, "user-id-123", bm.UserID)
			},
		},
		{
			name: "should return error when code generation fails",
			args: args{
				userID: "user-id-123",
				request: bookmarkDTO.CreateBookmarkRequest{
					Description: "Test Bookmark",
					URL:         "https://example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				bookmarkRepo *bookmarkMocks.Repository,
				codeGenerator *utilsMocks.CodeGenerator,
			) {
				codeGenerator.
					On("Generate", bookmarkCodeLength).
					Return("", assert.AnError).
					Once()
			},
			verifyResponse: func(t *testing.T, bm *model.Bookmark, err error) {
				assert.Error(t, err)
				assert.Nil(t, bm)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
		{
			name: "should return error when bookmark code already exists",
			args: args{
				userID: "user-id-123",
				request: bookmarkDTO.CreateBookmarkRequest{
					Description: "Test Bookmark",
					URL:         "https://example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				bookmarkRepo *bookmarkMocks.Repository,
				codeGenerator *utilsMocks.CodeGenerator,
			) {
				codeGenerator.
					On("Generate", bookmarkCodeLength).
					Return("abc123", nil).
					Once()

				expectedBookmark := &model.Bookmark{
					Description: "Test Bookmark",
					URL:         "https://example.com",
					Code:        "abc123",
					UserID:      "user-id-123",
				}

				bookmarkRepo.
					On("Create", ctx, matchBookmark(expectedBookmark)).
					Return(dbutils.ErrDuplicationType).
					Once()
			},
			verifyResponse: func(t *testing.T, bm *model.Bookmark, err error) {
				assert.Error(t, err)
				assert.Nil(t, bm)
				assert.ErrorIs(t, err, ErrBookmarkAlreadyExists)
			},
		},
		{
			name: "should return error when user not found",
			args: args{
				userID: "nonexistent-user",
				request: bookmarkDTO.CreateBookmarkRequest{
					Description: "Test Bookmark",
					URL:         "https://example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				bookmarkRepo *bookmarkMocks.Repository,
				codeGenerator *utilsMocks.CodeGenerator,
			) {
				codeGenerator.
					On("Generate", bookmarkCodeLength).
					Return("abc123", nil).
					Once()

				expectedBookmark := &model.Bookmark{
					Description: "Test Bookmark",
					URL:         "https://example.com",
					Code:        "abc123",
					UserID:      "nonexistent-user",
				}

				bookmarkRepo.
					On("Create", ctx, matchBookmark(expectedBookmark)).
					Return(dbutils.ErrForeignKeyViolationType).
					Once()
			},
			verifyResponse: func(t *testing.T, bm *model.Bookmark, err error) {
				assert.Error(t, err)
				assert.Nil(t, bm)
				assert.ErrorIs(t, err, ErrBadRequest)
			},
		},
		{
			name: "should return error when database operation fails",
			args: args{
				userID: "user-id-123",
				request: bookmarkDTO.CreateBookmarkRequest{
					Description: "Test Bookmark",
					URL:         "https://example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				bookmarkRepo *bookmarkMocks.Repository,
				codeGenerator *utilsMocks.CodeGenerator,
			) {
				codeGenerator.
					On("Generate", bookmarkCodeLength).
					Return("abc123", nil).
					Once()

				expectedBookmark := &model.Bookmark{
					Description: "Test Bookmark",
					URL:         "https://example.com",
					Code:        "abc123",
					UserID:      "user-id-123",
				}

				bookmarkRepo.
					On("Create", ctx, matchBookmark(expectedBookmark)).
					Return(assert.AnError).
					Once()
			},
			verifyResponse: func(t *testing.T, bm *model.Bookmark, err error) {
				assert.Error(t, err)
				assert.Nil(t, bm)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
		{
			name: "should return error when context is cancelled",
			args: args{
				userID: "user-id-123",
				request: bookmarkDTO.CreateBookmarkRequest{
					Description: "Test Bookmark",
					URL:         "https://example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				bookmarkRepo *bookmarkMocks.Repository,
				codeGenerator *utilsMocks.CodeGenerator,
			) {
				codeGenerator.
					On("Generate", bookmarkCodeLength).
					Return("abc123", nil).
					Once()

				expectedBookmark := &model.Bookmark{
					Description: "Test Bookmark",
					URL:         "https://example.com",
					Code:        "abc123",
					UserID:      "user-id-123",
				}

				bookmarkRepo.
					On("Create", ctx, matchBookmark(expectedBookmark)).
					Return(context.Canceled).
					Once()
			},
			verifyResponse: func(t *testing.T, bm *model.Bookmark, err error) {
				assert.Error(t, err)
				assert.Nil(t, bm)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var ctx context.Context
			bookmarkRepo := bookmarkMocks.NewRepository(t)
			codeGenerator := utilsMocks.NewCodeGenerator(t)

			// For context cancellation test, create a cancelled context
			if tc.name == "should return error when context is cancelled" {
				cancelledCtx, cancel := context.WithCancel(context.Background())
				cancel()
				ctx = cancelledCtx
			} else {
				ctx = context.Background()
			}

			tc.setupMocks(ctx, bookmarkRepo, codeGenerator)

			bookmarkService := NewService(bookmarkRepo, codeGenerator)

			bm, err := bookmarkService.Create(ctx, tc.args.userID, tc.args.request)

			tc.verifyResponse(t, bm, err)

			bookmarkRepo.AssertExpectations(t)
			codeGenerator.AssertExpectations(t)
		})
	}
}
