package handler

import (
	"gin-blog/app/models"
	g "gin-blog/internal/global"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Comment struct {
}

type AddCommentReq struct {
	PostId  int    `json:"post_id" form:"post_id"`
	Content string `json:"content" form:"content"`
}

type CommentQuery struct {
	PageQuery
}

// GetCommentList 获取评论列表
func (*Comment) GetCommentList(c *gin.Context) {
	var query CommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	data, total, err := models.GetCommentList(db, query.Page, query.Size)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[models.Comment]{
		List:  data,
		Total: total,
		Size:  query.Size,
		Page:  query.Page,
	})
}

func (*Comment) GetCommentListByPost(c *gin.Context) {
	var query CommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	postId, err := strconv.Atoi(c.Param(("id")))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	data, total, err := models.GetCommentListByPostId(db, postId, query.Page, query.Size)

	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[models.Comment]{
		List:  data,
		Total: total,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// SaveComment 保存评论（只能新增，不能编辑）
func (*Comment) SaveComment(c *gin.Context) {
	var req AddCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	// 过滤评论内容，防止 XSS 攻击
	req.Content = template.HTMLEscapeString(req.Content)
	db := GetDB(c)

	auth, _ := CurrentUserAuth(c)
	comment, err := models.AddComment(db, auth.ID, req.PostId, req.Content)

	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, comment)
}
