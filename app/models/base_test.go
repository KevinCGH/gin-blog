package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type user struct {
	UUID      uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	Age       int
	Enabled   bool
}

func initDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	db.AutoMigrate(user{})
	return db
}

func TestCreate(t *testing.T) {
	db := initDB()

	val, err := Create(db, &user{
		Name:    "test",
		Age:     11,
		Enabled: true,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, val.UUID)
}

func TestCount(t *testing.T) {
	db := initDB()

	db.Create(&user{Name: "user1", Email: "user1@example.com", Age: 10})
	count, err := Count(db, &user{})
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	db.Create(&user{Name: "user2", Email: "user2@example.com", Age: 20})
	count, err = Count(db, &user{})
	assert.Nil(t, err)
	assert.Equal(t, 2, count)

	db.Create(&user{Name: "user3", Email: "user3@example.com", Age: 30})
	count, err = Count(db, &user{})
	assert.Nil(t, err)
	assert.Equal(t, 3, count)

	count, err = Count(db, &user{}, "age >= ?", 20)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
}
