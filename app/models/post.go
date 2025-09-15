package models

import "gorm.io/gorm"

type Post struct {
	Model
	Title    string `gorm:"type:varchar(100);not null" json:"title"`
	Content  string `gorm:"type:varchar(1024)" json:"content"`
	UserId   int    `json:"user_id"`
	IsDelete bool   `json:"is_delete"`

	User *User `gorm:"foreignKey:UserId" json:"user"`
}

func GetBlogPost(db *gorm.DB, id int) (data *Post, err error) {
	result := db.Preload("User").Where(Post{Model: Model{ID: id}}).Where("is_delete = false").First(&data)
	return data, result.Error
}

// GetBlogPostList 前台文章列表
func GetBlogPostList(db *gorm.DB, page, size int) (data []Post, total int64, err error) {
	db = db.Model(Post{})
	db = db.Where("is_delete = false")

	db = db.Count(&total)
	result := db.Preload("User").Order("id DESC").Scopes(Paginate(page, size)).Find(&data)

	return data, total, result.Error
}

// DeletePost 物理删除文章
func DeletePost(db *gorm.DB, ids []int) (int64, error) {
	result := db.Where("id IN ?", ids).Delete(&Post{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func SaveOrUpdatePost(db *gorm.DB, post *Post) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var result *gorm.DB
		// 先添加/更新 文章，获取其 ID
		if post.ID == 0 {
			result = db.Create(&post)
		} else {
			result = db.Model(&post).Where("id", post.ID).Updates(post)
		}
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}
