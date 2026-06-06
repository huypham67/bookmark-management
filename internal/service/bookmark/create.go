package bookmark

import (
	"context"
	"errors"

	bookmarkDTO "github.com/huypham67/bookmark-service/internal/dto/bookmark"
	"github.com/huypham67/bookmark-service/internal/model"
	"github.com/huypham67/bookmark-service/pkg/dbutils"
	"github.com/rs/zerolog/log"
)

// Create creates a new bookmark for the user.
func (s *service) Create(ctx context.Context, userID string, req bookmarkDTO.CreateBookmarkRequest) (*model.Bookmark, error) {
	code, err := s.codeGenerator.Generate(bookmarkCodeLength)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to generate bookmark code")
		return nil, ErrInternalServerError
	}

	bm := &model.Bookmark{
		Description: req.Description,
		URL:         req.URL,
		Code:        code,
		UserID:      userID,
	}

	if err := s.bookmarkRepo.Create(ctx, bm); err != nil {
		switch {
		case errors.Is(err, dbutils.ErrDuplicationType):
			log.Warn().
				Str("user_id", userID).
				Str("code", code).
				Msg("bookmark code already exists")
			return nil, ErrBookmarkAlreadyExists
		case errors.Is(err, dbutils.ErrForeignKeyViolationType):
			log.Warn().
				Str("user_id", userID).
				Msg("user not found")
			return nil, ErrBadRequest
		default:
			log.Error().
				Err(err).
				Str("user_id", userID).
				Str("url", req.URL).
				Msg("failed to create bookmark")
			return nil, ErrInternalServerError
		}
	}

	log.Info().
		Str("bookmark_id", bm.ID).
		Str("user_id", userID).
		Str("code", bm.Code).
		Msg("bookmark created successfully")

	return bm, nil
}
