package bookmark

import (
	"context"

	"github.com/huypham67/bookmark-service/internal/model"
)

// GetPaginatedByUserID fetches a paginated list of bookmarks for a specific user.
func (r *repository) GetPaginatedByUserID(ctx context.Context, userID string, offset, limit int64, sort string) ([]*model.Bookmark, error) {
	var bookmarks []*model.Bookmark

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order(sort + " DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&bookmarks).Error

	return bookmarks, err
}

// CountByUserID counts the total number of bookmarks for a specific user.
func (r *repository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	return count, err
}

// GetByIDAndUserID fetches a bookmark by ID, verifying it belongs to the specified user.
func (r *repository) GetByIDAndUserID(ctx context.Context, id, userID string) (*model.Bookmark, error) {
	var bookmark *model.Bookmark

	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&bookmark).Error

	return bookmark, err
}
