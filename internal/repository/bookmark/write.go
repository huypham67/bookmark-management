package bookmark

import (
	"context"

	"github.com/huypham67/bookmark-service/internal/model"
)

// Create saves a new bookmark to the database.
func (r *repository) Create(ctx context.Context, bookmark *model.Bookmark) error {
	return r.db.WithContext(ctx).Create(bookmark).Error
}

// Update updates an existing bookmark for a specific user, only updating non-nil fields.
func (r *repository) Update(ctx context.Context, id, userID string, updates *model.Bookmark) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates).Error
}
