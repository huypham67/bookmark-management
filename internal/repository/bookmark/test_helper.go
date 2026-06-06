package bookmark

import (
	"testing"

	"github.com/huypham67/bookmark-service/internal/testutil"
	"gorm.io/gorm"
)

func newTestRepository(t *testing.T) (Repository, *gorm.DB) {
	t.Helper()

	testDB := testutil.NewTestDB(t, &testutil.BookmarkTestDB{})
	repo := NewRepository(testDB)

	return repo, testDB
}
