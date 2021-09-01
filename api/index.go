package api

import (
	"leekbox/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeDB interface {
	GetUserList(page, limit int) ([]model.User, error)
}

type IndexAPI struct {
	DB HomeDB
}

func (this *IndexAPI) Index(ctx *gin.Context) {
	ctx.DefaultQuery("userid", "87")
	name := ctx.Query("name")
	userid := ctx.Query("userid")
	if name == "" {
		name = "undefined"
	}
	if userid == "" {
		userid = "730811"
	}
	if users, err := this.DB.GetUserList(1, 100); err != nil {
		ctx.HTML(http.StatusOK, "index.html", gin.H{"items": nil, "userid": userid})
	} else {
		ctx.HTML(http.StatusOK, "index.html", gin.H{"items": users, "userid": userid})
	}
}
