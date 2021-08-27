package router

import (
	_ "leekbox/docs"
	"leekbox/handler"
	"leekbox/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(app *gin.Engine) {
	app.Static("/static", "static")
	// 注册html模板 渲染过滤器 需要用到的html模板
	app.SetFuncMap(utils.FuncMapUnion)
	app.LoadHTMLGlob("templates/*.html")
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.GET("/index", handler.IndexHandler)
	user := app.Group("/api/user")
	{
		user.GET("/info", handler.AuthMiddleWare(), handler.GetUserInfo)
		user.POST("/check", handler.CheckUserId)
		user.POST("/signup", handler.UserSignup)
		user.POST("/login", handler.UserLogin)
	}
}
