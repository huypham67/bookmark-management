// Package auth provides authentication services for user registration and login operations.
// It handles password hashing, token generation, and user validation.
package auth

import (
	"context"
	"errors"

	authDTO "github.com/huypham67/bookmark-service/internal/dto/auth"
	"github.com/huypham67/bookmark-service/internal/model"
	"github.com/huypham67/bookmark-service/internal/repository/user"
	"github.com/huypham67/bookmark-service/pkg/jwtutils"
	"github.com/huypham67/bookmark-service/pkg/security"
)

var (
	ErrUserAlreadyExists   = errors.New("username or email already exists")
	ErrInvalidCredentials  = errors.New("Invalid username or password")
	ErrInternalServerError = errors.New("internal server error")
)

// Service defines the interface for authentication operations, including user registration and login.
//
//go:generate mockery --name=Service --output=./mocks --outpkg=mocks --filename=mock_service.go
type Service interface {
	RegisterUser(ctx context.Context, req authDTO.RegisterUserRequest) (*model.User, error)
	LoginUser(ctx context.Context, req authDTO.LoginRequest) (string, error)
}

type service struct {
	userRepo       user.Repository
	passwordHasher security.PasswordHasher
	tokenGenerator jwtutils.TokenGenerator
}

// NewService creates a new instance of the authentication service with the provided dependencies.
func NewService(
	userRepo user.Repository,
	passwordHasher security.PasswordHasher,
	tokenGenerator jwtutils.TokenGenerator,
) Service {
	return &service{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}
