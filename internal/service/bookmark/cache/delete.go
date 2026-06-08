package cache

import (
	"context"

	"github.com/rs/zerolog/log"
)

// Delete deletes a bookmark and invalidates the user's list cache.
func (s *bookmarkCacheService) Delete(ctx context.Context, userID, bookmarkID string) error {
	if err := s.bookmarkService.Delete(ctx, userID, bookmarkID); err != nil {
		return err
	}

	// Invalidate all list caches for this user
	hashKey := buildUserCacheKey(userID)
	if err := s.cacheRepo.DeleteCacheByHashKey(ctx, hashKey); err != nil {
		log.Warn().
			Err(err).
			Str("user_id", userID).
			Msg("failed to invalidate user cache after bookmark deletion")
	}

	return nil
}
