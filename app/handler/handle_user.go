package handler

import (
	"gin-blog/app/models"
	g "gin-blog/internal/global"

	"github.com/gin-gonic/gin"
)

type User struct{}

func (*User) GetInfo(c *gin.Context) {
	user, err := CurrentUserAuth(c)
	if err != nil {
		ReturnError(c, g.ErrTokenRuntime, nil)
		return
	}

	userInfoVO := models.UserInfoVO{User: *user}

	ReturnSuccess(c, userInfoVO)
}
