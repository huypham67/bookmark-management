package testutil

import (
	"time"

	"github.com/huypham67/bookmark-service/internal/model"
	"gorm.io/gorm"
)

const (
	TestUserID1 = "user-uuid-1"
	TestUserID2 = "user-uuid-2"
)

var baseTime = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

type BookmarkTestDB struct {
	UserTestDB
}

func (b *BookmarkTestDB) MigrateDB() error {
	// Enable foreign key constraints for SQLite
	if err := b.db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return err
	}

	return b.db.AutoMigrate(&model.User{}, &model.Bookmark{})
}

func (b *BookmarkTestDB) SeedData() error {
	if err := b.UserTestDB.SeedData(); err != nil {
		return err
	}

	db := b.db.Session(&gorm.Session{SkipHooks: true})

	bookmarks := []*model.Bookmark{
		{
			BaseModel:   model.BaseModel{ID: "bookmark-1-1", CreatedAt: baseTime.Add(0 * time.Second)},
			Description: "Test Bookmark 1-1",
			URL:         "https://example.com/1",
			Code:        "code1001",
			UserID:      TestUserID1,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-1-2", CreatedAt: baseTime.Add(1 * time.Second)},
			Description: "Test Bookmark 1-2",
			URL:         "https://example.com/2",
			Code:        "code1002",
			UserID:      TestUserID1,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-1-3", CreatedAt: baseTime.Add(2 * time.Second)},
			Description: "Test Bookmark 1-3",
			URL:         "https://example.com/3",
			Code:        "code1003",
			UserID:      TestUserID1,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-1-4", CreatedAt: baseTime.Add(3 * time.Second)},
			Description: "Test Bookmark 1-4",
			URL:         "https://example.com/4",
			Code:        "code1004",
			UserID:      TestUserID1,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-1-5", CreatedAt: baseTime.Add(4 * time.Second)},
			Description: "Test Bookmark 1-5",
			URL:         "https://example.com/5",
			Code:        "code1005",
			UserID:      TestUserID1,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-1-6", CreatedAt: baseTime.Add(5 * time.Second)},
			Description: "Test Bookmark 1-6",
			URL:         "https://example.com/6",
			Code:        "code1006",
			UserID:      TestUserID1,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-1-7", CreatedAt: baseTime.Add(6 * time.Second)},
			Description: "Test Bookmark 1-7",
			URL:         "https://example.com/7",
			Code:        "code1007",
			UserID:      TestUserID1,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-1-8", CreatedAt: baseTime.Add(7 * time.Second)},
			Description: "Test Bookmark 1-8",
			URL:         "https://example.com/8",
			Code:        "code1008",
			UserID:      TestUserID1,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-2-1", CreatedAt: baseTime.Add(0 * time.Second)},
			Description: "Test Bookmark 2-1",
			URL:         "https://example.com/11",
			Code:        "code2001",
			UserID:      TestUserID2,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-2-2", CreatedAt: baseTime.Add(1 * time.Second)},
			Description: "Test Bookmark 2-2",
			URL:         "https://example.com/12",
			Code:        "code2002",
			UserID:      TestUserID2,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-2-3", CreatedAt: baseTime.Add(2 * time.Second)},
			Description: "Test Bookmark 2-3",
			URL:         "https://example.com/13",
			Code:        "code2003",
			UserID:      TestUserID2,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-2-4", CreatedAt: baseTime.Add(3 * time.Second)},
			Description: "Test Bookmark 2-4",
			URL:         "https://example.com/14",
			Code:        "code2004",
			UserID:      TestUserID2,
		},
		{
			BaseModel:   model.BaseModel{ID: "bookmark-2-5", CreatedAt: baseTime.Add(4 * time.Second)},
			Description: "Test Bookmark 2-5",
			URL:         "https://example.com/15",
			Code:        "code2005",
			UserID:      TestUserID2,
		},
	}

	err := db.CreateInBatches(bookmarks, 10).Error
	if err != nil {
		return err
	}

	return nil
}
