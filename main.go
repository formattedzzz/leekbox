package main

import (
	"fmt"
	_ "leekbox/dao"
	"leekbox/router"
	"time"

	"github.com/gin-gonic/gin"
)

var loc *time.Location

func init() {
	initLoc()
}

func initLoc() {
	if lac, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		loc = lac
	} else {
		fmt.Println(err)
	}
}

func main() {
	recover := func() {
		// 在必要的模块 panic之后需要重启一下
		if fatal := recover(); fatal != nil {
			fmt.Println("panic captured", fatal)
			return
		}
	}
	defer recover()

	app := gin.Default()
	router.InitRouter(app)

	err := app.Run(":7003")
	if err != nil {
		fmt.Println("something error", err)
	}
}
