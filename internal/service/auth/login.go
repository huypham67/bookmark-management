package auth

import (
	"context"
	"errors"

	authDTO "github.com/huypham67/bookmark-service/internal/dto/auth"
	"github.com/huypham67/bookmark-service/pkg/dbutils"
	"github.com/rs/zerolog/log"
)

// LoginUser authenticates a user by validating credentials and returns a JWT token.
func (s *service) LoginUser(ctx context.Context, req authDTO.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, dbutils.ErrRecordNotFoundType) {
		log.Error().
			Err(err).
			Str("username", req.Username).
			Msg("failed to get user by username")
		return "", ErrInternalServerError
	}

	if user == nil {
		log.Warn().
			Str("username", req.Username).
			Msg("user not found")
		return "", ErrInvalidCredentials
	}

	// Validate password
	if err := s.passwordHasher.Compare(user.Password, req.Password); err != nil {
		log.Warn().
			Str("username", req.Username).
			Msg("invalid password")
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.tokenGenerator.GenerateToken(user.ID, user.DisplayName, user.Email)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to generate token")
		return "", ErrInternalServerError
	}

	log.Info().
		Str("user_id", user.ID).
		Str("username", user.Username).
		Msg("user logged in successfully")

	return token, nil
}
