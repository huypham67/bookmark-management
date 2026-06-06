package user

import (
	"testing"

	"github.com/huypham67/bookmark-service/internal/testutil"
	"gorm.io/gorm"
)

func newTestRepository(t *testing.T) (Repository, *gorm.DB) {
	t.Helper()

	testDB := testutil.NewTestDB(t, &testutil.UserTestDB{})
	repo := NewRepository(testDB)

	return repo, testDB
}
