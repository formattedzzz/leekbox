package handler

import (
	"fmt"
	"leekbox/dao"
	"leekbox/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 这里的router本质就是一个struct类型为gin.Engine的指针
func initRouter(router *gin.Engine) {
	router.LoadHTMLFiles("templates/index.html")
}

func IndexHandler(ctx *gin.Context) {
	ctx.DefaultQuery("userid", "87")
	name := ctx.Query("name")
	userid := ctx.Query("userid")
	if name == "" {
		name = "undefined"
	}
	if userid == "" {
		userid = "730811"
	}
	fmt.Println(dao.DB)
	users := []model.User{}
	dao.DB.Table("users").Find(&users)
	ctx.HTML(http.StatusOK, "index.html", gin.H{"items": users})
}
