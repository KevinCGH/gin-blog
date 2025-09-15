package models

import (
	"errors"
	"gin-blog/internal/utils"
	"log/slog"
	"strconv"

	"gorm.io/gorm"
)

type User struct {
	Model
	Username string `gorm:"unique;type:varchar(50);not null" json:"username"`
	Password string `gorm:"type:varchar(100);not null" json:"-"`
	Email    string `gorm:"type:varchar(30);not null" json:"email"`
	Nickname string `gorm:"type:varchar(30);not null" json:"nickname"`
	Avatar   string `gorm:"type:varchar(1024)" json:"avatar"`
}

type UserInfoVO struct {
	User
}

func GetUserById(db *gorm.DB, id int) (*User, error) {
	var user User
	result := db.Model(&user).Where("id = ?", id).First(&user)
	return &user, result.Error
}

func GetUserByName(db *gorm.DB, name string) (*User, error) {
	var user User

	result := db.Model(&user).Where("username LIKE ?", name).First(&user)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &user, result.Error
}

func CreateNewUser(db *gorm.DB, username, password string) (*User, error) {
	num, err := Count(db, &User{})
	if err != nil {
		slog.Info(err.Error())
	}
	number := strconv.Itoa(num)
	pass, _ := utils.BcryptHash(password)
	user := &User{
		Username: username,
		Email:    username,
		Password: pass,
		Nickname: "游客" + number,
	}
	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, result.Error
}
