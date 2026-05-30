package link

import (
	"context"

	"github.com/huypham67/bookmark-service/pkg/dbutils"
)

// CheckExists checks whether the short code already exists.
func (r *repository) CheckExists(ctx context.Context, code string) (bool, error) {
	result, err := r.client.Exists(ctx, code).Result()
	if err != nil {
		return false, dbutils.ClassifyError(err)
	}
	return result > 0, nil
}

// GetLink retrieves original URL from Redis.
func (r *repository) GetLink(ctx context.Context, code string) (string, error) {
	return r.client.Get(ctx, code).Result()
}
