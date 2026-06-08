package cache

import (
	"context"

	bookmarkDTO "github.com/huypham67/bookmark-service/internal/dto/bookmark"
	"github.com/rs/zerolog/log"
)

// Update updates a bookmark and invalidates the user's list cache.
func (s *bookmarkCacheService) Update(ctx context.Context, userID, bookmarkID string, req bookmarkDTO.UpdateBookmarkRequest) error {
	if err := s.bookmarkService.Update(ctx, userID, bookmarkID, req); err != nil {
		return err
	}

	// Invalidate all list caches for this user
	hashKey := buildUserCacheKey(userID)
	if err := s.cacheRepo.DeleteCacheByHashKey(ctx, hashKey); err != nil {
		log.Warn().
			Err(err).
			Str("user_id", userID).
			Msg("failed to invalidate user cache after bookmark update")
	}

	return nil
}
