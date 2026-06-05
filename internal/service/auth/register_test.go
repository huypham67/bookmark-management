package auth

import (
	"context"
	"errors"
	"testing"

	authDTO "github.com/huypham67/bookmark-service/internal/dto/auth"
	"github.com/huypham67/bookmark-service/internal/model"
	userMocks "github.com/huypham67/bookmark-service/internal/repository/user/mocks"
	jwtutilsMocks "github.com/huypham67/bookmark-service/pkg/jwtutils/mocks"
	securityMocks "github.com/huypham67/bookmark-service/pkg/security/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func expectedAuthRegisteredUser() *model.User {
	return &model.User{
		DisplayName: "Test Display Name",
		Username:    "testuser",
		Email:       "testuser@gmail.com",
		Password:    "$2a$10$hashedpassword123456789",
	}
}

func matchAuthUser(expected *model.User) interface{} {
	return mock.MatchedBy(func(actual *model.User) bool {
		return actual.DisplayName == expected.DisplayName &&
			actual.Username == expected.Username &&
			actual.Email == expected.Email &&
			actual.Password == expected.Password
	})
}

func TestService_RegisterUser(t *testing.T) {
	t.Parallel()

	type args struct {
		request authDTO.RegisterUserRequest
	}

	testCases := []struct {
		name           string
		args           args
		setupMocks     func(context.Context, *userMocks.Repository, *securityMocks.PasswordHasher, *jwtutilsMocks.TokenGenerator)
		verifyResponse func(*testing.T, *model.User, error)
	}{
		{
			name: "should register user successfully",
			args: args{
				request: authDTO.RegisterUserRequest{
					DisplayName: "Test Display Name",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtutilsMocks.TokenGenerator,
			) {
				passwordHasher.
					On("Hash", "password123").
					Return("$2a$10$hashedpassword123456789", nil).
					Once()

				expectedUser := expectedAuthRegisteredUser()

				userRepo.
					On(
						"Create",
						ctx,
						matchAuthUser(expectedUser),
					).
					Return(nil).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.NoError(t, err)
				require.NotNil(t, user)

				assert.Equal(t, "Test Display Name", user.DisplayName)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "testuser@gmail.com", user.Email)
				assert.Equal(t, "$2a$10$hashedpassword123456789", user.Password)
			},
		},
		{
			name: "should return error when user already exists",
			args: args{
				request: authDTO.RegisterUserRequest{
					DisplayName: "Test Display Name",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtutilsMocks.TokenGenerator,
			) {
				passwordHasher.
					On("Hash", "password123").
					Return("$2a$10$hashedpassword123456789", nil).
					Once()

				expectedUser := expectedAuthRegisteredUser()

				duplicateError := errors.New(`ERROR: duplicate key value violates unique constraint "idx_users_email" (SQLSTATE 23505)`)

				userRepo.
					On(
						"Create",
						ctx,
						matchAuthUser(expectedUser),
					).
					Return(duplicateError).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.ErrorIs(t, err, ErrUserAlreadyExists)
			},
		},
		{
			name: "should return error when password hashing fails",
			args: args{
				request: authDTO.RegisterUserRequest{
					DisplayName: "Test Display Name",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtutilsMocks.TokenGenerator,
			) {
				passwordHasher.
					On("Hash", "password123").
					Return("", assert.AnError).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
		{
			name: "should return error when creating user in database fails",
			args: args{
				request: authDTO.RegisterUserRequest{
					DisplayName: "Test Display Name",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtutilsMocks.TokenGenerator,
			) {
				passwordHasher.
					On("Hash", "password123").
					Return("$2a$10$hashedpassword123456789", nil).
					Once()

				expectedUser := expectedAuthRegisteredUser()

				userRepo.
					On(
						"Create",
						ctx,
						matchAuthUser(expectedUser),
					).
					Return(assert.AnError).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			userRepo := new(userMocks.Repository)
			passwordHasher := securityMocks.NewPasswordHasher(t)
			tokenGenerator := jwtutilsMocks.NewTokenGenerator(t)

			tc.setupMocks(ctx, userRepo, passwordHasher, tokenGenerator)

			authService := NewService(userRepo, passwordHasher, tokenGenerator)

			user, err := authService.RegisterUser(
				ctx,
				tc.args.request,
			)

			tc.verifyResponse(t, user, err)

			userRepo.AssertExpectations(t)
		})
	}
}
