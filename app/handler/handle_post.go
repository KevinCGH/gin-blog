package handler

import (
	"gin-blog/app/models"
	g "gin-blog/internal/global"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Post struct {
}

type AddOrEditPostReq struct {
	ID      int    `json:"id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type PostQuery struct {
	PageQuery
}

func (*Post) GetPostList(c *gin.Context) {
	var query PostQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, _, err := models.GetBlogPostList(GetDB(c), query.Page, query.Size)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

func (*Post) GetPostInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(("id")))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	val, err := models.GetBlogPost(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, val)
}

func (*Post) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditPostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	auth, _ := CurrentUserAuth(c)

	post := models.Post{
		Model:   models.Model{ID: req.ID},
		Title:   req.Title,
		Content: req.Content,
		UserId:  auth.ID,
	}

	err := models.SaveOrUpdatePost(db, &post)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, post)
}

func (*Post) DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(("id")))
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	val, err := models.GetBlogPost(db, id)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	auth, _ := CurrentUserAuth(c)
	if val.UserId != auth.ID {
		ReturnError(c, g.ErrUserHasNoPermission, nil)
		return
	}

	rows, err := models.DeletePost(GetDB(c), []int{id})
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}
