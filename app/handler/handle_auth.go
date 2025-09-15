package handler

import (
	"errors"
	"gin-blog/app/models"
	"gin-blog/config"
	g "gin-blog/internal/global"
	"gin-blog/internal/utils"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserAuth struct{}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=20"`
}

type LoginVO struct {
	models.User
	Token string `json:"token"`
}

type CustomClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (u *UserAuth) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)

	userAuth, err := models.GetUserByName(db, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 检查密码是否正确
	if !utils.BcryptCheck(req.Password, userAuth.Password) {
		ReturnError(c, g.ErrPassword, nil)
		return
	}

	// 登录成功，生成 Token
	claims := CustomClaims{
		userAuth.ID,
		userAuth.Username,
		jwt.RegisteredClaims{
			//Issuer:    "12",
			//Subject:   "",
			//Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			//NotBefore: jwt.NewNumericDate(time.Now()),
			//IssuedAt:  jwt.NewNumericDate(time.Now()),
			//ID: strconv.Itoa(userAuth.ID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	conf := config.Conf.JWT
	tokenString, err := token.SignedString([]byte(conf.Secret))
	if err != nil {
		ReturnError(c, g.ErrTokenCreate, err)
		return
	}

	slog.Info("用户登录成功：" + userAuth.Username)

	session := sessions.Default(c)
	session.Set(g.CTX_USER_AUTH, userAuth.ID)
	session.Save()

	ReturnSuccess(c, LoginVO{
		User:  *userAuth,
		Token: tokenString,
	})
}

func (*UserAuth) Register(c *gin.Context) {
	var regreq RegisterReq
	if err := c.ShouldBindJSON(&regreq); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}
	regreq.Email = utils.Format(regreq.Email)

	auth, err := models.GetUserByName(GetDB(c), regreq.Email)
	if err != nil {
		var flag bool = false
		if errors.Is(err, gorm.ErrRecordNotFound) {
			flag = true
		}
		if !flag {
			ReturnError(c, g.ErrDbOp, err)
			return
		}
	}

	if auth != nil {
		ReturnError(c, g.ErrUserExist, err)
		return
	}

	info := utils.GenEmailVerificationInfo(regreq.Email, regreq.Password)
	// TODO code 记录到 redis 中
	emailData := utils.GetEmailData(regreq.Email, info)
	err = utils.SendEmail(regreq.Email, emailData)
	if err != nil {
		ReturnError(c, g.ErrSendEmail, err)
		return
	}

	ReturnSuccess(c, nil)
}

// VerifyCode 验证邮箱
func (*UserAuth) VerifyCode(c *gin.Context) {
	var code string
	if code = c.Query("info"); code == "" {
		returnErrorPage(c)
		return
	}

	username, password, err := utils.ParseEmailVerificationInfo(code)
	if err != nil {
		returnErrorPage(c)
		return
	}

	//注册用户
	_, err = models.CreateNewUser(GetDB(c), username, password)
	if err != nil {
		returnErrorPage(c)
		return
	}

	// 注册成功，返回成功页面
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>注册成功</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #5cb85c;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>注册成功</h1>
                <p>恭喜您，注册成功！</p>
            </div>
        </body>
        </html>
    `))
}

func (*UserAuth) Logout(c *gin.Context) {
	c.Set(g.CTX_USER_AUTH, nil)

	// 已经退出登录
	auth, _ := CurrentUserAuth(c)
	if auth == nil {
		ReturnSuccess(c, nil)
		return
	}

	session := sessions.Default(c)
	session.Delete(g.CTX_USER_AUTH)
	session.Save()

	// TODO: Redis 中在线状态

	ReturnSuccess(c, nil)
}

func returnErrorPage(c *gin.Context) {
	c.Data(http.StatusInternalServerError, "text/html;charset=utf-8", []byte(`
		<!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>注册失败</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #d9534f;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>注册失败</h1>
                <p>请重试。</p>
            </div>
        </body>
        </html>
	`))
}
