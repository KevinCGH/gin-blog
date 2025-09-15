package handler

import (
	"errors"
	"gin-blog/app/models"
	g "gin-blog/internal/global"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Response 响应结构体
type Response[T any] struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 响应消息
	Data    T      `json:"data"`    // 响应数据
}

// ReturnHttpResponse HTTP 码 + 业务码 + 消息 + 数据
func ReturnHttpResponse(c *gin.Context, httpCode, code int, msg string, data any) {
	c.JSON(httpCode, Response[any]{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// ReturnResponse 业务码 + 数据
func ReturnResponse(c *gin.Context, r g.Result, data any) {
	ReturnHttpResponse(c, http.StatusOK, r.Code(), r.Msg(), data)
}

func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c, g.OkResult, data)
}

// ReturnError
// 所有可预料的错误 = 业务错误 + 系统错误，在业务层面处理，返回 HTTP 200 状态码
// 对于不可预料的错误，会触发 panic，由 gin 中间件捕获，并返回 HTTP 500 状态码
// err 是业务错误，data 是错误数据 （可以是 error 或 string）
func ReturnError(c *gin.Context, r g.Result, data any) {
	slog.Info("[Func-ReturnError] " + r.Msg())

	var val = r.Msg()

	if data != nil {
		switch v := data.(type) {
		case error:
			val = v.Error()
		case string:
			val = v
		}
		slog.Error(val)
	}

	c.AbortWithStatusJSON(
		http.StatusOK,
		Response[any]{
			Code:    r.Code(),
			Message: r.Msg(),
			Data:    val,
		})
}

// PageQuery 分页获取数据
type PageQuery struct {
	Page    int    `form:"page_num"`  // 当前页数（从 1 开始）
	Size    int    `form:"page_size"` // 每页条数
	Keyword string `form:"keyword"`   // 搜索关键字
}

type PageResult[T any] struct {
	Page  int   `json:"page_num"`  // 当前页数（从 1 开始）
	Size  int   `json:"page_size"` // 每页条数
	Total int64 `json:"total"`     // 总条数
	List  []T   `json:"page_data"` // 分页数据
}

// GetDB 获取 *gorm.DB
func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(g.CTX_DB).(*gorm.DB)
}

/*
	CurrentUser 获取当前登录用户信息

获取当前登录用户信息
1. 能从 gin Context 上获取到 user 对象，说明本次请求链路中获取过了
2. 从 session 中获取 uid
3. 根据 uid 获取用户信息，并设置到 gin Context 上
*/
func CurrentUserAuth(c *gin.Context) (*models.User, error) {
	key := g.CTX_USER_AUTH

	// 1
	if cache, exist := c.Get(key); exist && cache != nil {
		slog.Debug("[Func-CurrentUserAuth] get from cache: " + cache.(*models.User).Username)
		return cache.(*models.User), nil
	}

	// 2
	session := sessions.Default(c)
	id := session.Get(key)
	if id == nil {
		return nil, errors.New("session 中没有 user_auth_id")
	}

	// 3
	db := GetDB(c)
	user, err := models.GetUserById(db, id.(int))
	if err != nil {
		return nil, err
	}

	c.Set(key, user)
	return user, nil
}
