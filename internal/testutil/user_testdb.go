package testutil

import (
	"github.com/huypham67/bookmark-service/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const TestPassword = "password123"

type UserTestDB struct {
	baseTestDB
}

func (u *UserTestDB) MigrateDB() error {
	return u.db.AutoMigrate(&model.User{})
}

func (u *UserTestDB) SeedData() error {
	db := u.db.Session(&gorm.Session{SkipHooks: true})

	users := []model.User{
		{
			BaseModel: model.BaseModel{
				ID: "user-uuid-1",
			},
			DisplayName: "Test User 1",
			Username:    "testuser1",
			Email:       "testuser1@gmail.com",
			Password:    hashPassword(TestPassword),
		},
		{
			BaseModel: model.BaseModel{
				ID: "user-uuid-2",
			},
			DisplayName: "Test User 2",
			Username:    "testuser2",
			Email:       "testuser2@gmail.com",
			Password:    hashPassword(TestPassword),
		},
		{
			BaseModel: model.BaseModel{
				ID: "user-uuid-3",
			},
			DisplayName: "Test User 3",
			Username:    "testuser3",
			Email:       "testuser3@gmail.com",
			Password:    hashPassword(TestPassword),
		},
	}

	err := db.CreateInBatches(users, 10).Error
	if err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		panic(err)
	}

	return string(hash)
}
