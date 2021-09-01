package router

import (
	"leekbox/api"
	"leekbox/config"
	"leekbox/dao"
	_ "leekbox/docs"
	"leekbox/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Create(db *dao.GormDB, config config.Configuration) *gin.Engine {
	userHandler := api.UserAPI{
		DB:        db,
		UserEvent: &api.UserEvent{},
	}
	indexHander := api.IndexAPI{
		DB: db,
	}
	app := gin.Default()
	app.Static("/static", "static")
	// 注册html模板 渲染过滤器 需要用到的html模板
	app.SetFuncMap(utils.FuncMapUnion)
	app.LoadHTMLGlob("templates/*.html")
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("version", func(c *gin.Context) {
		c.JSON(200, config.VERSION)
	})
	app.GET("/index", indexHander.Index)

	user := app.Group("/api/user")
	{
		user.GET("/info", api.AuthMiddleWare(), userHandler.GetUserInfo)
		user.POST("/check", userHandler.CheckUserId)
		user.POST("/signup", userHandler.UserSignup)
		user.POST("/login", userHandler.UserLogin)
	}
	return app
}
