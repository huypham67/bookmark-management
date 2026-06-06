package link

import (
	"context"

	"github.com/rs/zerolog/log"
)

// GetOriginalURL retrieves the original URL for a given shortened code.
func (s *service) GetOriginalURL(ctx context.Context, code string) (string, error) {
	url, err := s.linkRepo.GetLink(ctx, code)

	if err != nil {
		log.Error().
			Err(err).
			Str("code", code).
			Msg("failed to retrieve original URL from Redis")
		return "", err
	}

	return url, nil
}
