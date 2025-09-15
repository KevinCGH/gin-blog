package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserById(t *testing.T) {
	db := initDB()
	db.AutoMigrate(User{})

	var user = User{Username: "test", Password: "123456", Nickname: "Test", Email: "test@test.com"}
	db.Create(&user)

	{
		res, err := GetUserById(db, user.ID)
		assert.Nil(t, err)
		assert.Equal(t, "test", res.Username)
	}
}
