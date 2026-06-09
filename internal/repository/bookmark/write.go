package bookmark

import (
	"context"

	"github.com/huypham67/bookmark-service/internal/model"
	"github.com/huypham67/bookmark-service/pkg/dbutils"
)

// Create saves a new bookmark to the database.
func (r *repository) Create(ctx context.Context, bookmark *model.Bookmark) error {
	if err := r.db.WithContext(ctx).Create(bookmark).Error; err != nil {
		return dbutils.ClassifyError(err)
	}
	return nil
}

// NextCodeInt generates the next integer code for a new bookmark.
func (r *repository) NextCodeInt(ctx context.Context) (int64, error) {
	query := "SELECT COALESCE(MAX(code_int), 0) + 1 FROM bookmarks"
	if r.db.Dialector.Name() == "postgres" {
		query = "SELECT nextval(pg_get_serial_sequence('bookmarks', 'code_int'))"
	}

	var n int64
	if err := r.db.WithContext(ctx).Raw(query).Scan(&n).Error; err != nil {
		return 0, dbutils.ClassifyError(err)
	}
	return n, nil
}

// Update updates an existing bookmark for a specific user, only updating non-nil fields.
// Returns the number of rows affected and any error encountered.
func (r *repository) Update(ctx context.Context, id, userID string, updates *model.Bookmark) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates)

	if result.Error != nil {
		return 0, dbutils.ClassifyError(result.Error)
	}

	return result.RowsAffected, nil
}

// Delete deletes a bookmark for a specific user.
// Returns the number of rows affected and any error encountered.
func (r *repository) Delete(ctx context.Context, id, userID string) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Bookmark{})

	if result.Error != nil {
		return 0, dbutils.ClassifyError(result.Error)
	}

	return result.RowsAffected, nil
}
