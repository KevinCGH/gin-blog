package routers

import (
	"gin-blog/app/handler"
	"gin-blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

var (
	userAuthAPI handler.UserAuth
	userAPI     handler.User
	postAPI     handler.Post
	commentAPI  handler.Comment
)

func SetupRouter(r *gin.Engine) {
	setupBaseRouter(r)
	setupBlogRouter(r)
}

// setupBaseRouter 通用接口，全部不需要 登录 + 鉴权
func setupBaseRouter(r *gin.Engine) {
	base := r.Group("/api")

	base.POST("/login", userAuthAPI.Login)            // ✅ 登录
	base.POST("/register", userAuthAPI.Register)      // ✅ 注册
	base.GET("/logout", userAuthAPI.Logout)           // ✅ 登出
	base.GET("/email/verify", userAuthAPI.VerifyCode) // ✅ 验证注册
}

// setupBlogRouter 博客前端接口，大部分不需要登录，部分需要登录
func setupBlogRouter(r *gin.Engine) {
	base := r.Group("/api/front")

	post := base.Group("/post")
	{
		post.GET("/list", postAPI.GetPostList)                         // ✅ 前台文章列表
		post.GET("/:id", postAPI.GetPostInfo)                          // ✅ 文章详情
		post.GET("/:id/comment/list", commentAPI.GetCommentListByPost) // ✅ 获取文章对应的评论列表
	}

	comment := base.Group("/comment")
	{
		comment.GET("/list", commentAPI.GetCommentList) // ✅ 前台评论列表
	}

	// 需要登录
	base.Use(middleware.JWTAuth())
	{
		base.GET("/user/info", userAPI.GetInfo)

		base.POST("/post", postAPI.SaveOrUpdate)      // ✅ 创建文章
		base.DELETE("/post/:id", postAPI.DeletePost)  // ✅ 删除自己的文章
		base.POST("/comment", commentAPI.SaveComment) // ✅ 创建评论
	}

}
