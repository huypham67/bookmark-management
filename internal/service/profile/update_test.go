package profile

import (
	"context"
	"errors"
	"testing"

	profileDTO "github.com/huypham67/bookmark-service/internal/dto/profile"
	"github.com/huypham67/bookmark-service/internal/model"
	"github.com/huypham67/bookmark-service/internal/repository/user/mocks"
	"github.com/huypham67/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_UpdateUserInfo(t *testing.T) {
	t.Parallel()

	type args struct {
		userID string
		req    profileDTO.UpdateUserRequest
	}

	testCases := []struct {
		name           string
		args           args
		setupMocks     func(context.Context, *mocks.Repository)
		verifyResponse func(*testing.T, error)
	}{
		{
			name: "should update user successfully when email belongs to current user",
			args: args{
				userID: "user-id",
				req: profileDTO.UpdateUserRequest{
					DisplayName: "Updated User",
					Email:       "john@example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				mockRepo *mocks.Repository,
			) {
				mockRepo.
					On(
						"GetByID",
						ctx,
						"user-id",
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "user-id",
							},
						},
						nil,
					).
					Once()

				mockRepo.
					On(
						"GetByEmail",
						ctx,
						"john@example.com",
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "user-id",
							},
						},
						nil,
					).
					Once()

				mockRepo.
					On(
						"Update",
						ctx,
						mock.MatchedBy(func(user *model.User) bool {
							return user.ID == "user-id" &&
								user.DisplayName == "Updated User" &&
								user.Email == "john@example.com"
						}),
					).
					Return(nil).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				err error,
			) {
				assert.NoError(t, err)
			},
		},
		{
			name: "should update user successfully when email does not exist",
			args: args{
				userID: "user-id",
				req: profileDTO.UpdateUserRequest{
					DisplayName: "Updated User",
					Email:       "new@example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				mockRepo *mocks.Repository,
			) {
				mockRepo.
					On(
						"GetByID",
						ctx,
						"user-id",
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "user-id",
							},
						},
						nil,
					).
					Once()

				mockRepo.
					On(
						"GetByEmail",
						ctx,
						"new@example.com",
					).
					Return(
						nil,
						dbutils.ErrRecordNotFoundType,
					).
					Once()

				mockRepo.
					On(
						"Update",
						ctx,
						mock.AnythingOfType("*model.User"),
					).
					Return(nil).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				err error,
			) {
				assert.NoError(t, err)
			},
		},
		{
			name: "should return error when user not found",
			args: args{
				userID: "nonexistent-id",
				req: profileDTO.UpdateUserRequest{
					DisplayName: "Updated User",
					Email:       "john@example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				mockRepo *mocks.Repository,
			) {
				mockRepo.
					On(
						"GetByID",
						ctx,
						"nonexistent-id",
					).
					Return(
						nil,
						dbutils.ErrRecordNotFoundType,
					).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				err error,
			) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error when checking email fails",
			args: args{
				userID: "user-id",
				req: profileDTO.UpdateUserRequest{
					DisplayName: "Updated User",
					Email:       "john@example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				mockRepo *mocks.Repository,
			) {
				mockRepo.
					On(
						"GetByID",
						ctx,
						"user-id",
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "user-id",
							},
						},
						nil,
					).
					Once()

				mockRepo.
					On(
						"GetByEmail",
						ctx,
						"john@example.com",
					).
					Return(
						nil,
						errors.New("database error"),
					).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				err error,
			) {
				assert.Error(t, err)
			},
		},
		{
			name: "should return error when email belongs to another user",
			args: args{
				userID: "user-id",
				req: profileDTO.UpdateUserRequest{
					DisplayName: "Updated User",
					Email:       "existing@example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				mockRepo *mocks.Repository,
			) {
				mockRepo.
					On(
						"GetByID",
						ctx,
						"user-id",
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "user-id",
							},
						},
						nil,
					).
					Once()

				mockRepo.
					On(
						"GetByEmail",
						ctx,
						"existing@example.com",
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "another-user",
							},
						},
						nil,
					).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				err error,
			) {
				assert.Equal(
					t,
					err,
					ErrEmailAlreadyRegistered,
				)
			},
		},
		{
			name: "should return error when updating user fails",
			args: args{
				userID: "user-id",
				req: profileDTO.UpdateUserRequest{
					DisplayName: "Updated User",
					Email:       "john@example.com",
				},
			},
			setupMocks: func(
				ctx context.Context,
				mockRepo *mocks.Repository,
			) {
				mockRepo.
					On(
						"GetByID",
						ctx,
						"user-id",
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "user-id",
							},
						},
						nil,
					).
					Once()

				mockRepo.
					On(
						"GetByEmail",
						ctx,
						"john@example.com",
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "user-id",
							},
						},
						nil,
					).
					Once()

				mockRepo.
					On(
						"Update",
						ctx,
						mock.AnythingOfType("*model.User"),
					).
					Return(
						errors.New("update failed"),
					).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				err error,
			) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockRepo := new(mocks.Repository)

			tc.setupMocks(
				ctx,
				mockRepo,
			)

			service := NewService(
				mockRepo,
			)

			err := service.UpdateUserInfo(
				ctx,
				tc.args.userID,
				tc.args.req,
			)

			tc.verifyResponse(
				t,
				err,
			)

			mockRepo.AssertExpectations(t)
		})
	}
}
