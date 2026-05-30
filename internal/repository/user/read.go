package user

import (
	"context"

	"github.com/huypham67/bookmark-service/internal/model"
	"github.com/huypham67/bookmark-service/pkg/dbutils"
)

// GetByEmail retrieves a user by their email address.
func (r *repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, dbutils.ClassifyError(err)
	}
	return user, nil
}

// GetByUsername retrieves a user by their username.
func (r *repository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user *model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, dbutils.ClassifyError(err)
	}
	return user, nil
}

// GetByID retrieves a user by their ID.
func (r *repository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	var user *model.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, dbutils.ClassifyError(err)
	}
	return user, nil
}
