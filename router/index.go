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
	indexHander := api.IndexAPI{DB: db}
	roomHander := api.RoomAPI{DB: db}
	app := gin.Default()
	app.Static("/static", "static")
	app.SetFuncMap(utils.FuncMapUnion)
	app.LoadHTMLGlob("templates/*.html")
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("/version", func(c *gin.Context) {
		c.JSON(200, config.VERSION)
	})
	app.GET("/index", indexHander.Index)

	// /api开头的接口 有token的话默认取一下token
	app.Use(api.AttachToken())
	user := app.Group("/api/user")
	{
		user.POST("/check", userHandler.CheckUserId)
		user.POST("/signup", userHandler.UserSignup)
		user.POST("/login", userHandler.UserLogin)
		user.GET("/info", api.AuthMiddleWare(), userHandler.GetUserInfo)
		user.POST("/update", api.AuthMiddleWare(), userHandler.UpdateUserInfo)
	}
	room := app.Group("/api/room")
	{
		room.GET("/:id", roomHander.GetRoomInfo)
		room.POST("/create", api.AuthMiddleWare(), roomHander.CreateNewRoom)
		room.POST("/update", api.AuthMiddleWare(), roomHander.UpdateRoomInfo)
	}
	return app
}
