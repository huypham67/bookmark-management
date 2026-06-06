package auth

import (
	"context"
	"testing"

	authDTO "github.com/huypham67/bookmark-service/internal/dto/auth"
	"github.com/huypham67/bookmark-service/internal/model"
	userMocks "github.com/huypham67/bookmark-service/internal/repository/user/mocks"
	"github.com/huypham67/bookmark-service/pkg/dbutils"
	jwtutilsMocks "github.com/huypham67/bookmark-service/pkg/jwtutils/mocks"
	"github.com/huypham67/bookmark-service/pkg/security"
	securityMocks "github.com/huypham67/bookmark-service/pkg/security/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_LoginUser(t *testing.T) {
	t.Parallel()

	type args struct {
		request authDTO.LoginRequest
	}

	testCases := []struct {
		name           string
		args           args
		setupMocks     func(context.Context, *userMocks.Repository, *securityMocks.PasswordHasher, *jwtutilsMocks.TokenGenerator)
		verifyResponse func(*testing.T, string, error)
	}{
		{
			name: "should login user successfully",
			args: args{
				request: authDTO.LoginRequest{
					Username: "testuser",
					Password: "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtutilsMocks.TokenGenerator,
			) {
				user := &model.User{
					BaseModel: model.BaseModel{
						ID: "user-id-123",
					},
					DisplayName: "Test User",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "$2a$10$hashedpassword123456789",
				}

				userRepo.
					On("GetByUsername", ctx, "testuser").
					Return(user, nil).
					Once()

				passwordHasher.
					On("Compare", "$2a$10$hashedpassword123456789", "password123").
					Return(nil).
					Once()

				tokenGenerator.
					On("GenerateToken", "user-id-123", "Test User", "testuser@gmail.com").
					Return("jwt-token-123", nil).
					Once()
			},
			verifyResponse: func(t *testing.T, token string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "jwt-token-123", token)
			},
		},
		{
			name: "should return error when user not found",
			args: args{
				request: authDTO.LoginRequest{
					Username: "nonexistent",
					Password: "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtutilsMocks.TokenGenerator,
			) {
				userRepo.
					On("GetByUsername", ctx, "nonexistent").
					Return(nil, dbutils.ErrRecordNotFoundType).
					Once()
			},
			verifyResponse: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.ErrorIs(t, err, ErrInvalidCredentials)
			},
		},
		{
			name: "should return error when password is invalid",
			args: args{
				request: authDTO.LoginRequest{
					Username: "testuser",
					Password: "wrongpassword",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtutilsMocks.TokenGenerator,
			) {
				user := &model.User{
					BaseModel: model.BaseModel{
						ID: "user-id-123",
					},
					DisplayName: "Test User",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "$2a$10$hashedpassword123456789",
				}

				userRepo.
					On("GetByUsername", ctx, "testuser").
					Return(user, nil).
					Once()

				passwordHasher.
					On("Compare", "$2a$10$hashedpassword123456789", "wrongpassword").
					Return(security.ErrPasswordMismatch).
					Once()
			},
			verifyResponse: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.ErrorIs(t, err, ErrInvalidCredentials)
			},
		},
		{
			name: "should return error when getting user fails",
			args: args{
				request: authDTO.LoginRequest{
					Username: "testuser",
					Password: "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtutilsMocks.TokenGenerator,
			) {
				userRepo.
					On("GetByUsername", ctx, "testuser").
					Return(nil, assert.AnError).
					Once()
			},
			verifyResponse: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
		{
			name: "should return error when token generation fails",
			args: args{
				request: authDTO.LoginRequest{
					Username: "testuser",
					Password: "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtutilsMocks.TokenGenerator,
			) {
				user := &model.User{
					BaseModel: model.BaseModel{
						ID: "user-id-123",
					},
					DisplayName: "Test User",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "$2a$10$hashedpassword123456789",
				}

				userRepo.
					On("GetByUsername", ctx, "testuser").
					Return(user, nil).
					Once()

				passwordHasher.
					On("Compare", "$2a$10$hashedpassword123456789", "password123").
					Return(nil).
					Once()

				tokenGenerator.
					On("GenerateToken", "user-id-123", "Test User", "testuser@gmail.com").
					Return("", assert.AnError).
					Once()
			},
			verifyResponse: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
		{
			name: "should return error when context is cancelled",
			args: args{
				request: authDTO.LoginRequest{
					Username: "testuser",
					Password: "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtutilsMocks.TokenGenerator,
			) {
				userRepo.
					On("GetByUsername", ctx, "testuser").
					Return(nil, context.Canceled).
					Once()
			},
			verifyResponse: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var ctx context.Context
			userRepo := userMocks.NewRepository(t)
			passwordHasher := securityMocks.NewPasswordHasher(t)
			tokenGenerator := jwtutilsMocks.NewTokenGenerator(t)

			// For context cancellation test, create a cancelled context
			if tc.name == "should return error when context is cancelled" {
				cancelledCtx, cancel := context.WithCancel(context.Background())
				cancel()
				ctx = cancelledCtx
			} else {
				ctx = context.Background()
			}

			tc.setupMocks(ctx, userRepo, passwordHasher, tokenGenerator)

			authService := NewService(userRepo, passwordHasher, tokenGenerator)

			token, err := authService.LoginUser(ctx, tc.args.request)

			tc.verifyResponse(t, token, err)

			userRepo.AssertExpectations(t)
		})
	}
}
