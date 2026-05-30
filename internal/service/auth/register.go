package auth

import (
	"context"
	"errors"

	authDTO "github.com/huypham67/bookmark-service/internal/dto/auth"
	"github.com/huypham67/bookmark-service/internal/model"
	"github.com/huypham67/bookmark-service/pkg/dbutils"
	"github.com/rs/zerolog/log"
)

// RegisterUser registers a new user by validating input, hashing password, and saving to database.
func (s *service) RegisterUser(ctx context.Context, req authDTO.RegisterUserRequest) (*model.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, dbutils.ErrRecordNotFoundType) {
		log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("failed to check if email exists")
		return nil, ErrInternalServerError
	}

	if existingUser != nil {
		log.Warn().
			Str("email", req.Email).
			Msg("email already registered")
		return nil, ErrEmailAlreadyRegistered
	}

	// Check if username already exists
	existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, dbutils.ErrRecordNotFoundType) {
		log.Error().
			Err(err).
			Str("username", req.Username).
			Msg("failed to check if username exists")
		return nil, ErrInternalServerError
	}

	if existingUser != nil {
		log.Warn().
			Str("username", req.Username).
			Msg("username already exists")
		return nil, ErrUsernameAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to hash password")
		return nil, ErrInternalServerError
	}

	// Create new user
	user := &model.User{
		DisplayName: req.DisplayName,
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPassword,
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error().
			Err(err).
			Str("email", req.Email).
			Str("username", req.Username).
			Msg("failed to register user")
		return nil, ErrInternalServerError
	}

	log.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Str("username", user.Username).
		Msg("user registered successfully")

	return user, nil
}
