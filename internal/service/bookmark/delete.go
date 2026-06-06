package bookmark

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Delete deletes an existing bookmark for the user.
func (s *service) Delete(ctx context.Context, userID, bookmarkID string) error {
	_, err := s.bookmarkRepo.GetByIDAndUserID(ctx, bookmarkID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().
				Str("user_id", userID).
				Str("bookmark_id", bookmarkID).
				Msg("bookmark not found")
			return ErrBookmarkNotFound
		}

		log.Error().
			Err(err).
			Str("user_id", userID).
			Str("bookmark_id", bookmarkID).
			Msg("failed to fetch bookmark")
		return ErrInternalServerError
	}

	if err := s.bookmarkRepo.Delete(ctx, bookmarkID, userID); err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Str("bookmark_id", bookmarkID).
			Msg("failed to delete bookmark")
		return ErrInternalServerError
	}

	log.Info().
		Str("bookmark_id", bookmarkID).
		Str("user_id", userID).
		Msg("bookmark deleted successfully")

	return nil
}
