package middleware

import (
	g "gin-blog/internal/global"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WithGormDB 将 gorm.DB 注入到 gin.Context
// handler 中通过 c.MustGet(g.CTX_DB).(*gorm.DB) 获取使用
func WithGormDB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(g.CTX_DB, db)
		c.Next()
	}
}

func WithCookieStore(name, secret string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{Path: "/", MaxAge: 600})
	return sessions.Sessions(name, store)
}
