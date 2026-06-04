package bookmark

import (
	"context"

	"github.com/huypham67/bookmark-service/internal/model"
	"github.com/rs/zerolog/log"
)

// List retrieves a paginated list of bookmarks for the user.
func (s *service) List(ctx context.Context, userID string, page, limit int64, sort string) ([]*model.Bookmark, *PaginationResult, error) {

	offset := (page - 1) * limit

	bookmarks, err := s.bookmarkRepo.GetPaginatedByUserID(ctx, userID, offset, limit, sort)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Int64("page", page).
			Int64("limit", limit).
			Msg("failed to fetch paginated bookmarks")
		return nil, nil, ErrInternalServerError
	}

	total, err := s.bookmarkRepo.CountByUserID(ctx, userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to count bookmarks")
		return nil, nil, ErrInternalServerError
	}

	log.Info().
		Str("user_id", userID).
		Int64("page", page).
		Int64("limit", limit).
		Int64("total", total).
		Msg("bookmarks listed successfully")

	return bookmarks, &PaginationResult{Total: total}, nil
}
