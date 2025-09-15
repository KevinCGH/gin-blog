package middleware

import (
	"gin-blog/app/handler"
	"gin-blog/app/models"
	"gin-blog/config"
	g "gin-blog/internal/global"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// JWTAuth 基于 JWT 的授权
// 如果存在 session，则直接从 session 中获取用户信息
// 如果不存在 session，则从 Authorization 中获取 token，并解析 token 获取用户信息，并设置到 session 中
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("[middleware-JWTAuth] auth")

		db := c.MustGet(g.CTX_DB).(*gorm.DB)

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			handler.ReturnError(c, g.ErrTokenNotExist, nil)
			return
		}

		// token 的正确格式： `Bearer [tokenString]`
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handler.ReturnError(c, g.ErrTokenType, nil)
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &handler.CustomClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(config.Conf.JWT.Secret), nil
		})
		if err != nil {
			handler.ReturnError(c, g.ErrTokenWrong, err)
			return
		}
		claims, _ := token.Claims.(*handler.CustomClaims)

		// 判断 token 已过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			handler.ReturnError(c, g.ErrTokenRuntime, nil)
			return
		}

		user, err := models.GetUserById(db, claims.UserID)
		if err != nil {
			handler.ReturnError(c, g.ErrUserNotExist, err)
			return
		}

		// session
		session := sessions.Default(c)
		session.Set(g.CTX_USER_AUTH, claims.UserID)
		session.Save()

		c.Set(g.CTX_USER_AUTH, user)
	}
}
