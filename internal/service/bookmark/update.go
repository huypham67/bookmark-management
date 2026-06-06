package bookmark

import (
	"context"

	bookmarkDTO "github.com/huypham67/bookmark-service/internal/dto/bookmark"
	"github.com/huypham67/bookmark-service/internal/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Update updates an existing bookmark for the user.
func (s *service) Update(ctx context.Context, userID, bookmarkID string, req bookmarkDTO.UpdateBookmarkRequest) error {
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

	updates := &model.Bookmark{}
	if req.Description != nil {
		updates.Description = *req.Description
	}
	if req.URL != nil {
		updates.URL = *req.URL
	}

	if err := s.bookmarkRepo.Update(ctx, bookmarkID, userID, updates); err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Str("bookmark_id", bookmarkID).
			Msg("failed to update bookmark")
		return ErrInternalServerError
	}

	log.Info().
		Str("bookmark_id", bookmarkID).
		Str("user_id", userID).
		Msg("bookmark updated successfully")

	return nil
}
