package models

import (
	"time"

	"gorm.io/gorm"
)

func MakeMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},    // 用户
		&Post{},    // 文章
		&Comment{}, // 评论
	)
}

type Model struct {
	ID        int       `gorm:"primaryKey;auto_increment" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Gorm Scopes

// Paginate 分页
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

// 通用 CRUD

// Create 创建数据
func Create[T any](db *gorm.DB, data *T) (*T, error) {
	result := db.Create(data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// Get 单条数据查询
func Get[T any](db *gorm.DB, data *T, query string, args ...any) (*T, error) {
	result := db.Where(query, args...).First(data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// List 列表查询
func List[T any](db *gorm.DB, data T, slt, order, query string, args ...any) (T, error) {
	db = db.Model(data).Select(slt).Order(order)
	if query != "" {
		db = db.Where(query, args...)
	}
	result := db.Find(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

// Count 统计数据
func Count[T any](db *gorm.DB, data *T, where ...any) (int, error) {
	var total int64
	db = db.Model(data)
	if len(where) > 0 {
		db = db.Where(where[0], where[1:]...)
	}
	result := db.Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(total), nil
}
