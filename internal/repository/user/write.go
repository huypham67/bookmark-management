package user

import (
	"context"

	"github.com/huypham67/bookmark-service/internal/model"
)

// Create saves a new user to the database.
func (r *repository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update updates user's display_name and email.
func (r *repository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"display_name": user.DisplayName,
		"email":        user.Email,
	}).Error
}

