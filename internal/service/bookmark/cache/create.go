package cache

import (
	"context"

	bookmarkDTO "github.com/huypham67/bookmark-service/internal/dto/bookmark"
	"github.com/huypham67/bookmark-service/internal/model"
	"github.com/rs/zerolog/log"
)

// Create creates a new bookmark and invalidates the user's list cache.
func (s *bookmarkCacheService) Create(ctx context.Context, userID string, req bookmarkDTO.CreateBookmarkRequest) (*model.Bookmark, error) {
	bookmark, err := s.bookmarkService.Create(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	// Invalidate all list caches for this user
	hashKey := buildUserCacheKey(userID)
	if err := s.cacheRepo.DeleteCacheByHashKey(ctx, hashKey); err != nil {
		log.Warn().
			Err(err).
			Str("user_id", userID).
			Msg("failed to invalidate user cache after bookmark creation")
	}

	return bookmark, nil
}
