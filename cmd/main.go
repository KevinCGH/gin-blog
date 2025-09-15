package main

import (
	"gin-blog/app/routers"
	"gin-blog/config"
	gb "gin-blog/internal/bootstrap"
	"gin-blog/internal/middleware"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	conf := config.Conf
	_ = gb.InitLogger(conf)
	db := gb.InitDatabase(conf)

	// 初始化 gin 服务
	gin.SetMode(conf.Server.Mode)
	r := gin.New()
	r.SetTrustedProxies([]string{"*"})
	if conf.Server.Mode == "debug" {
		r.Use(gin.Logger(), gin.Recovery()) // gin 自带的日志和恢复中间件
	} else {
		// 后续可以自定义中间件
		r.Use(gin.Logger(), gin.Recovery())
	}
	r.Use(middleware.WithGormDB(db))
	r.Use(middleware.WithCookieStore(conf.Session.Name, conf.Session.Salt))
	routers.SetupRouter(r)

	serverAddr := config.Conf.Server.Port
	if serverAddr[0] == ':' || strings.HasPrefix(serverAddr, "0.0.0.0:") {
		log.Printf("Serving HTTP on (http://localhost%s/) ... \n", strings.Split(serverAddr, ":")[1])
	} else {
		log.Printf("Serving HTTP on (http://%s/) ... \n", serverAddr)
	}
	err := r.Run(serverAddr)
	if err != nil {
		return
	}
}
