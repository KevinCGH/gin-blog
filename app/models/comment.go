package models

import "gorm.io/gorm"

type Comment struct {
	Model
	Content string `gorm:"unique;type:varchar(1024)" json:"content"`
	UserId  int    `json:"user_id"`
	PostId  int    `json:"post_id"`

	User User `gorm:"foreignKey:UserId" json:"user"`
	Post Post `gorm:"foreignKey:PostId" json:"post"`
}

// GetCommentList 获取博客评论列表
func GetCommentList(db *gorm.DB, page, size int) (data []Comment, total int64, err error) {
	return GetCommentListByPostId(db, -1, page, size)
}

// GetCommentListByPostId 获取博客 文章对应的评论列表
func GetCommentListByPostId(db *gorm.DB, postId, page, size int) (data []Comment, total int64, err error) {
	var list []Comment
	tx := db.Model(&Comment{})
	tx.Count(&total).Preload("User").Preload("Post")
	if postId > -1 {
		tx.Where("post_id = ?", postId)
	}
	tx.Order("id DESC").Scopes(Paginate(page, size))
	if err := tx.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func AddComment(db *gorm.DB, userId, postId int, content string) (*Comment, error) {
	comment := Comment{
		Content: content,
		PostId:  postId,
		UserId:  userId,
	}
	result := db.Create(&comment)
	return &comment, result.Error

}
