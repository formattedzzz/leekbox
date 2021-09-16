package router

import (
	"leekbox/api"
	"leekbox/api/auth"
	"leekbox/api/stream"
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
	indexHander := api.IndexAPI{DB: db, Config: &config}
	streamHander := stream.New([]string{"localhost"})
	roomHander := api.RoomAPI{DB: db, Stream: streamHander}

	app := gin.Default()
	app.Static("/static", "static")
	app.SetFuncMap(utils.FuncMapUnion)
	app.LoadHTMLGlob("templates/*.html")
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("/index", indexHander.Index)
	app.GET("/version", indexHander.Version)

	// /api开头的接口 有token的话默认取一下token
	app.Use(auth.AttachToken())
	user := app.Group("/api/user")
	{
		user.POST("/check", userHandler.CheckUserId)
		user.POST("/signup", userHandler.UserSignup)
		user.POST("/login", userHandler.UserLogin)
		user.GET("/info", auth.AuthMiddleWare(), userHandler.GetUserInfo)
		user.PUT("/update", auth.AuthMiddleWare(), userHandler.UpdateUserInfo)
		user.GET("/rooms", auth.AuthMiddleWare(), userHandler.GetUserSubRooms)
	}
	room := app.Group("/api/room")
	{
		room.GET("/:id", roomHander.GetRoomInfo)
		room.GET("/comments", roomHander.GetRoomComments)
		room.POST("/create", auth.AuthMiddleWare(), roomHander.CreateNewRoom)
		room.PUT("/update", auth.AuthMiddleWare(), roomHander.UpdateRoomInfo)
		room.POST("/subscribe", auth.AuthMiddleWare(), roomHander.CreateSubscribe)
	}
	comment := app.Group("/api/comment")
	{
		comment.POST("/create", auth.AuthMiddleWare(), roomHander.CreateNewComment)
	}
	app.GET("/api/stream", streamHander.Handler)
	return app
}
