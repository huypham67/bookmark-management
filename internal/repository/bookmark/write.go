package bookmark

import (
	"context"

	"github.com/huypham67/bookmark-service/internal/model"
)

// Create saves a new bookmark to the database.
func (r *repository) Create(ctx context.Context, bookmark *model.Bookmark) error {
	return r.db.WithContext(ctx).Create(bookmark).Error
}
